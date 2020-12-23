import requests, signal, time
from datetime import timedelta
from json import dumps, loads
from logging import getLogger, WARNING, StreamHandler, Formatter
from pathlib import Path
from shutil import move as shutil_move
from requests_toolbelt import MultipartEncoder
from typing import Dict, List
from subprocess import PIPE, Popen, STDOUT

# Setup logging for runner.py
# Create a custom logger
logger = getLogger(__name__)

# Create handlers
console_handler = StreamHandler()
console_handler.setLevel(WARNING)

# Create formatters and add it to handlers
console_format = Formatter("%(name)s|%(levelname)s|%(lineno)d|%(message)s")
console_handler.setFormatter(console_format)

# Add handlers to the logger
logger.addHandler(console_handler)

class JobRunner:
	def __init__(self, controller_ip: str="localhost:5000", runner_name=""):
		self.controller_ip = controller_ip
		self.runner_name = runner_name

		self.__current_job_status = {
			"percentage": None,
			"job_elapsed_time": None,
			"stage_estimated_time_remaining": None,
			"fps": None,
			"stage_elapsed_time": None,
			"stage": None
		}
		self.__running = True

		self.__current_uuid = None

		# Setup self.stop as a handler to handle terminate signals
		signal.signal(signal.SIGINT, self.stop)
		signal.signal(signal.SIGTERM, self.stop)

	def stop(self, *args):
		# This method accepts *args because the signal module calls with extra info that we don't care about when shutting down
		logger.info("Stopping RedCedarRunner")
		self.__running = False

	def run(self):
		"""Runs the JobRunner
		"""
		while self.__running:
			new_job_info = self.new_job_from_controller()
			self.__start_job(new_job_info)

	def new_job_from_controller(self):
		"""Sends a get request to the controller for a new job
		"""
		for i in range(100):
			if self.__running:
				r = requests.get(f"http://{self.controller_ip}/api/v1/job/request", headers={"redcedar-runner-name": self.runner_name}, stream=True)

				if r.status_code != 200:
					logger.warning(f"Received status code {r.status_code} from controller because of error: {r.content}. Retrying in {i} seconds")
					if not self.__running:
						return None
					time.sleep(i)
					continue

				job_info = loads(r.headers.get("x-rc-job-info"))
				input_file = Path.cwd() / f"input{Path(job_info['file']).suffix}" # Creates an input file with the same suffix as the input

				if input_file.exists():
					input_file.unlink()

				with input_file.open("wb") as f:
					for chunk in r.iter_content(1024):
						f.write(chunk)

				job_info["in_file"] = str(input_file)

				return job_info
			else:
				return None

		raise RuntimeError(f"Controller did not respond with new job after 100 tries")

	def __start_job(self, job_info: Dict):
		self.__current_uuid = job_info["uuid"]
		self._run_job(job_info)

	def _run_ffmpeg(self, ffmpeg_command: List[str], interval_callback=None, *args):
		if type(ffmpeg_command) != list:
			raise TypeError("ffmpeg_command argument must be a list of strings")

		subtitle_failure = False
		logger.info(f"Running ffmpeg using command list {ffmpeg_command}")

		if not self.__running:
			logger.debug(f"Stopping ffmpeg because __running is {self.__running}")
			return

		with Popen(ffmpeg_command, stdout=PIPE, stderr=STDOUT, universal_newlines=True) as p:
			for line in p.stdout:
				logger.debug(f"Got FFMPEG line: {line}")
				if "Subtitle codec" in line and "is not supported" in line:
					subtitle_failure = True
					logger.debug(f"Triggered subtitle_failure on line: {line}")
				fps: float = None
				current_time: str = None
				speed: float = None
				split_line = line.split(" ")
				for _line in split_line:
					if "fps" in _line:
						try:
							fps = float(_line.replace("fps=", ""))
						except:
							fps = None
					elif "time" in _line:
						try:
							current_time = _line.replace("time=", "")
						except:
							current_time = None
					elif "speed" in _line:
						try:
							speed = float(_line.replace("speed=", "").replace("x", ""))
						except:
							speed = None

				if interval_callback != None:
					interval_callback(fps, current_time, speed, *args)

				if not self.__running:
					return

				time.sleep(0.05)
			exit_code = p.poll()
			while exit_code == None:
				time.sleep(0.1)
				logger.debug("Rechecking for non-None exit code")
				exit_code = p.poll()
			logger.debug(f"ffmpeg command returned exit code: {exit_code}")
			if exit_code != 0 and subtitle_failure:
				raise SubtitleError("Subtitle failure detected while running ffmpeg")
			return exit_code

	def _run_job(self, job_info: Dict):
		input_file = Path(job_info["in_file"])
		is_hevc = job_info["is_hevc"]
		has_stereo = job_info["has_stereo"]
		is_interlaced = job_info["is_interlaced"]

		logger.info(f"Running job {job_info['file']} which has characteristics: [is_hevc: {is_hevc}, has_stereo: {has_stereo}, is_interlaced: {is_interlaced}]")

		current_job_warnings, current_job_errors = ([], [])
		critical_failure = False

		output_file = Path.cwd() / "output.mkv"

		if output_file.exists():
			output_file.unlink()

		job_start_time = time.time()

		downmixed_audio: Path = None

		if not has_stereo:
			# Extract audio to cwd/job_uuid-extracted-audio.mkv
			extracted_audio = Path.cwd() / f"{job_info['uuid']}-extracted-audio.mka"
			if extracted_audio.exists():
				extracted_audio.unlink()
			extract_audio_command = ["ffmpeg", "-i", str(input_file), "-map", "-0:v?", "-map", "0:a:0", "-map", "-0:s?", "-c", "copy", str(extracted_audio)]
			self.__current_job_status["stage"] = "Extract Audio"
			self.send_current_job_status()
			logger.info("Starting extraction of audio")
			stage_start_time = time.time()
			self._run_ffmpeg(extract_audio_command, self.update_job_status, stage_start_time, job_start_time, job_info)

			# Downmix audio
			downmixed_audio = Path.cwd() / f"{job_info['uuid']}-downmixed-audio.mka"
			downmix_audio_command = ["ffmpeg", "-i", str(extracted_audio), "-map", "0:a", "-c", "aac", "-af", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", str(downmixed_audio)]
			self.__current_job_status["stage"] = "Downmix Audio"
			self.send_current_job_status()
			logger.info("Starting downmixing of audio")
			stage_start_time = time.time()
			self._run_ffmpeg(downmix_audio_command, self.update_job_status, stage_start_time, job_start_time, job_info)

			if extracted_audio.exists():
				extracted_audio.unlink()

		encode_inputs = ["-i", str(input_file)]
		tracks_to_copy = ["-map", "0:s?", "-map", "0:a"]
		encoding_commands = []

		if downmixed_audio != None:
			encode_inputs.append("-i")
			encode_inputs.append(str(downmixed_audio))
			tracks_to_copy.append("-map")
			tracks_to_copy.append("1:a")

		if is_hevc:
			tracks_to_copy.append("-map")
			tracks_to_copy.append("0:v")
		else:
			encoding_commands = ["-map", "0:v", "-vcodec", "hevc"]

		final_ffmpeg_command = ["ffmpeg"] + encode_inputs + tracks_to_copy + ["-c", "copy"] + encoding_commands + [str(output_file)]
		self.__current_job_status["stage"] = "Final encode"
		self.send_current_job_status()
		logger.info("Starting final encode")
		stage_start_time = time.time()

		try:
			ffmpeg_exit_code = self._run_ffmpeg(final_ffmpeg_command, self.update_job_status, stage_start_time, job_start_time, job_info)
		except SubtitleError:
			if output_file.exists():
				try:
					output_file.unlink()
				except PermissionError:
					current_job_errors.append("Could not delete interim output file while trying without subtitles")
					logger.critical("Could not delete interim output file while trying without subtitles", exc_info=True)
					critical_failure = True
			if not critical_failure:
				current_job_warnings.append("Final encode failed because of suspected unsupported subtitle codec, retrying without subtitles")
				logger.warning("Final encode failed because of suspected unsupported subtitle codec, retrying without subtitles")
				tracks_to_copy = tracks_to_copy[2:]
				final_ffmpeg_command = ["ffmpeg"] + encode_inputs + tracks_to_copy + ["-c", "copy"] + encoding_commands + [str(output_file)]
				ffmpeg_exit_code = self._run_ffmpeg(final_ffmpeg_command, self.update_job_status, stage_start_time, job_start_time, job_info)

		if ffmpeg_exit_code != 0 and self.__running and not critical_failure:
			critical_failure = True
			logger.critical(f"Final encode ffmpeg command returned non-zero exit code: {ffmpeg_exit_code}")
			current_job_errors.append(f"Final encode ffmpeg command returned non-zero exit code: {ffmpeg_exit_code}")
		else:
			if self.__running:
				logger.info(f"Completed final encode for {input_file}")

		if downmixed_audio != None and downmixed_audio.exists():
			downmixed_audio.unlink()

		if not self.__running:
			return

		history_entry = {
			"file": job_info["file"],
			"datetime_completed": time.time(),
			"warnings": current_job_warnings,
			"errors": current_job_errors
		}

		if not critical_failure:
			self.send_job_complete(history_entry, output_file)
			if input_file.exists():
				input_file.unlink()
			if output_file.exists():
				output_file.unlink()
		else:
			try:
				output_file.unlink()
			except PermissionError:
				logger.warning("Could not remove output file during critical failure cleanup")

			try:
				input_file.unlink()
			except PermissionError:
				logger.warning("Could not remove input file during critical failure cleanup")

			self.send_job_complete(history_entry, None)

	def update_job_status(self, fps: float, current_file_time_str: str, current_speed: float, stage_start_time: float, job_start_time: float, job_info: Dict):
		if fps == None:
			fps = "N/A"

		if current_file_time_str == None or current_speed == None:
			self.send_current_job_status()
			return

		total_length = str([track for track in job_info["media_info"]["tracks"] if track["kind_of_stream"] == "General"][0]["duration"])
		total_length = float(f"{total_length[:-3]}.{total_length[-3:]}")

		reversed_times = [float(x) for x in current_file_time_str.split(":")[::-1]] # Now in format [seconds, minutes, hours, etc]
		time_conversion_mask = [1.0, 60.0, 3600] # Multipliers to convert specific value to seconds
		current_file_time_seconds = sum([reversed_times[i] * time_conversion_mask[i] for i in range(len(reversed_times))])

		current_time = time.time()
		remaining_time = total_length - current_file_time_seconds

		if current_speed == 0:
			current_speed = 0.000000000000000000001

		self.__current_job_status = {
			"percentage": f"{round((current_file_time_seconds / total_length) * 100, 2)}",
			"job_elapsed_time": chop_ms(timedelta(seconds=(current_time - job_start_time))),
			"stage_estimated_time_remaining": chop_ms(timedelta(seconds=remaining_time / current_speed)),
			"fps": fps,
			"stage_elapsed_time": chop_ms(timedelta(seconds=(current_time - stage_start_time))),
			"stage": self.__current_job_status["stage"]
		}
		self.send_current_job_status()

		return (True, "")

	def send_current_job_status(self):
		def x():
			if self.__current_uuid == None:
				logger.warning(f"Current job status failed to send because self.__current_uuid is None")
				return
			r = requests.post(f"http://{self.controller_ip}/api/v1/job/status", json={
												"uuid": self.__current_uuid,
												"status": {"percentage": str(self.__current_job_status["percentage"]),
													"job_elapsed_time": str(self.__current_job_status["job_elapsed_time"]),
													"stage_estimated_time_remaining": str(self.__current_job_status["stage_estimated_time_remaining"]),
													"fps": str(self.__current_job_status["fps"]),
													"stage_elapsed_time": str(self.__current_job_status["stage_elapsed_time"]),
													"stage": str(self.__current_job_status["stage"])}
												})
			if r.status_code != 200:
				if r.status_code == 409:
					logger.critical("Detected self as unresponsive while send job status. Shutting down")
					self.__running = False
					return
				logger.warning(f"Current job status failed to send because of error: {r.content}")

		# Just trying to see if this makes the ui less laggy and weird
		# start_new_thread(x, ())
		x()

	def send_job_complete(self, history_entry, output_file_path: Path):
		# TODO: Allow for output_file_path to be None, signalling to not copy any file
		for i in range(100):
			if self.__current_uuid == None:
				logger.error(f"Failed to send job complete signal because self.__current_uuid is None")
				return

			with output_file_path.open("rb") as f:
				m = MultipartEncoder(fields={"file": (output_file_path.name, f)})
				r = requests.post(f"http://{self.controller_ip}/api/v1/job/complete",
					data=m,
					headers={
						"Content-Type": m.content_type,
						"x-rc-history-entry": dumps({"uuid": self.__current_uuid, "history": history_entry})
					})

			if r.status_code != 200:
				if r.status_code == 409:
					logger.warning("Detected self as unresponsive while sending job complete signal")
					return
				logger.warning(f"Job complete failed to send because of error: {r.content}. Retrying in {i} seconds...")
				if not self.__running:
					logger.error("Exiting without sending job complete signal to controller")
					return
				time.sleep(i)
			else:
				return

def chop_ms(delta):
	return delta - timedelta(microseconds=delta.microseconds)

class SubtitleError(RuntimeError):
	pass
