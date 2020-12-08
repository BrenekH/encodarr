import time
from copy import copy
from datetime import datetime, timedelta
from flask_socketio import SocketIO
from logging import INFO, getLogger, ERROR, WARNING, StreamHandler, FileHandler, Formatter
from pathlib import Path
from shutil import move as shutil_move
from typing import Dict, List
from subprocess import PIPE, Popen, STDOUT

# Setup logging for runner.py
# Create a custom logger
logger = getLogger(__name__)

# Create handlers
console_handler = StreamHandler()
# file_handler = FileHandler("/config/log.log")
console_handler.setLevel(WARNING)
# file_handler.setLevel(INFO)

# Create formatters and add it to handlers
console_format = Formatter("%(name)s|%(levelname)s|%(lineno)d|%(message)s")
# file_format = Formatter("%(asctime)s|%(name)s|%(levelname)s|%(lineno)d|%(message)s")
console_handler.setFormatter(console_format)
# file_handler.setFormatter(file_format)

# Add handlers to the logger
logger.addHandler(console_handler)
# logger.addHandler(file_handler)

class JobRunner:
	def __init__(self, socket_io: SocketIO):
		self.socket_io = socket_io
		self.active = False
		self.__completed_jobs = [] # Contains dictionaries with file, datetime_completed(in UTC), warnings, and errors keys
		self.__current_job_status = {
			"percentage": None,
			"job_elapsed_time": None,
			"stage_estimated_time_remaining": None,
			"fps": None,
			"stage_elapsed_time": None,
			"stage": None
		}
		self.__waiting_jobs = [] # In case new_job gets called even when runner is active
		self.__current_job = {
			"file": None,
			"uuid": None
		}
		self.__running = False

	def stop(self):
		logger.info("Stopping JobRunner")
		self.__running = False

	def new_job(self, job_info: Dict):
		self.active = True
		if self.__current_job["uuid"] != None:
			self.__waiting_jobs.append(job_info)
			return

		self.__start_job(job_info)

	def completed_jobs(self):
		to_return = copy(self.__completed_jobs)
		self.__completed_jobs = []
		return to_return

	def __start_job(self, job_info: Dict):
		self.__running = True
		self.__current_job = {
			"file": job_info["file"],
			"uuid": job_info["uuid"]
		}

		self.emit_current_job()
		self.socket_io.start_background_task(self._run_job, job_info)

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

				self.socket_io.sleep(0.05)
			exit_code = p.poll()
			while exit_code == None:
				self.socket_io.sleep(0.1)
				logger.debug("Rechecking for non-None exit code")
				exit_code = p.poll()
			logger.debug(f"ffmpeg command returned exit code: {exit_code}")
			if exit_code != 0 and subtitle_failure:
				raise SubtitleError("Subtitle failure detected while running ffmpeg")
			return exit_code

	def _run_job(self, job_info: Dict):
		input_file = Path(job_info["file"])
		is_hevc = job_info["is_hevc"]
		has_stereo = job_info["has_stereo"]
		is_interlaced = job_info["is_interlaced"]

		logger.info(f"Running job {input_file} which has characteristics: [is_hevc: {is_hevc}, has_stereo: {has_stereo}, is_interlaced: {is_interlaced}]")

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
			self.emit_current_job_status()
			logger.info("Starting extraction of audio")
			stage_start_time = time.time()
			self._run_ffmpeg(extract_audio_command, self.update_job_status, stage_start_time, job_start_time, job_info)

			# Downmix audio
			downmixed_audio = Path.cwd() / f"{job_info['uuid']}-downmixed-audio.mka"
			downmix_audio_command = ["ffmpeg", "-i", str(extracted_audio), "-map", "0:a", "-c", "aac", "-af", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", str(downmixed_audio)]
			self.__current_job_status["stage"] = "Downmix Audio"
			self.emit_current_job_status()
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
		self.emit_current_job_status()
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
			logger.info(f"Completed final encode for {input_file}")

		if downmixed_audio != None and downmixed_audio.exists():
			downmixed_audio.unlink()

		if not self.__running:
			return

		if not critical_failure:
			# Remove original file
			delete_successful = False
			if input_file.exists():
				try:
					input_file.unlink()
					delete_successful = True
				except PermissionError:
					current_job_warnings.append(f"Could not delete {input_file}, adding '-RedCedar' as a suffix")
					logger.warning(f"Could not delete {input_file}", exc_info=True)

			# Move output.mkv to take the original file's place
			if output_file.exists():
				if delete_successful:
					shutil_move(str(output_file), input_file.with_suffix(output_file.suffix))	# Retains the file suffix with the new name
				else:
					shutil_move(str(output_file), input_file.with_name(f"{input_file.stem}-RedCedar").with_suffix(output_file.suffix))	# Retains the file suffix with the new name
			logger.info("Output file copied")
		else:
			try:
				output_file.unlink()
			except PermissionError:
				logger.debug("Could not remove output_file during critical failure cleanup")

		if self.__running:
			self.__completed_jobs.append({
				"file": job_info["file"],
				"datetime_completed": datetime.utcnow().timestamp(),
				"warnings": current_job_warnings,
				"errors": current_job_errors
			})

		self.__current_job_status = {
			"percentage": None,
			"job_elapsed_time": None,
			"stage_estimated_time_remaining": None,
			"fps": None,
			"stage_elapsed_time": None,
			"stage": None
		}
		self.emit_current_job_status()

		self.__current_job = {"file": None, "uuid": None}

		if len(self.__waiting_jobs) > 0:
			next_job = self.__waiting_jobs.pop(0)
			self.__start_job(next_job)
		else:
			self.active = False

	def update_job_status(self, fps: float, current_file_time_str: str, current_speed: float, stage_start_time: float, job_start_time: float, job_info: Dict):
		if fps == None:
			fps = "N/A"

		if current_file_time_str == None or current_speed == None:
			self.emit_current_job_status()
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
		self.emit_current_job_status()

		return (True, "")

	def emit_current_job_status(self):
		self.emit_event("current_job_status_update", {"percentage": str(self.__current_job_status["percentage"]),
												"job_elapsed_time": str(self.__current_job_status["job_elapsed_time"]),
												"stage_estimated_time_remaining": str(self.__current_job_status["stage_estimated_time_remaining"]),
												"fps": str(self.__current_job_status["fps"]),
												"stage_elapsed_time": str(self.__current_job_status["stage_elapsed_time"]),
												"stage": str(self.__current_job_status["stage"])})

	def emit_current_job(self):
		self.emit_event("current_job_update", {"file": str(self.__current_job["file"])})

	def emit_event(self, event_name: str, data):
		logger.debug(f"Emitting event {event_name} with data: {data}")
		self.socket_io.emit(event_name, data, namespace="/updates")

def chop_ms(delta):
	return delta - timedelta(microseconds=delta.microseconds)

class SubtitleError(RuntimeError):
	pass
