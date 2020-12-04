import time
from copy import copy
from datetime import timedelta
from flask_socketio import SocketIO
from json import JSONDecodeError, loads
from pathlib import Path
from shutil import move as shutil_move
from typing import Dict, List, Tuple
from subprocess import DEVNULL, PIPE, Popen, STDOUT

class JobRunner:
	def __init__(self, socket_io: SocketIO):
		self.socket_io = socket_io
		self.active = False
		self.__current_job_status = {
			"percentage": None,
			"elapsed_time": None,
			"estimated_time_remaining": None,
			"fps": None,
			"stage_elapsed_time": None,
			"stage": None
		}
		self.__waiting_jobs = [] # In case new_job gets called even when runner is active

	def new_job(self, job_info: Dict):
		self.active = True
		if self.__current_job_status["uuid"] != None:
			self.__waiting_jobs.append(job_info)
			return

		self.__start_job(job_info)

	def get_job_status(self) -> Dict:
		return copy(self.__current_job_status)

	def __start_job(self, job_info: Dict):
		self.socket_io.start_background_task(self._run_job, job_info)

	def _run_ffmpeg(self, ffmpeg_command: List[str], interval_callback=None, *args):
		if type(ffmpeg_command) != list:
			raise TypeError("ffmpeg_command argument must be a list of strings")

		with Popen(ffmpeg_command, stdout=PIPE, stderr=STDOUT, universal_newlines=True) as p:
			for line in p.stdout:
				print(line)
				fps: float = None
				current_time: str = None
				speed: float = None
				split_line = line.split(" ")
				for _line in split_line:
					if "fps" in _line:
						fps = float(_line.replace("fps=", ""))
					elif "time" in _line:
						current_time = _line.replace("time=", "")
					elif "speed" in _line:
						speed = float(_line.replace("speed=", "").replace("x", ""))

				if interval_callback != None and fps != None and current_time != None and speed != None:
					interval_callback(fps, current_time, speed, *args)

				self.socket_io.sleep(0.1)

	def _run_job(self, job_info: Dict):
		input_file = Path(job_info["file"])
		is_hevc = job_info["is_hevc"]
		has_stereo = job_info["has_stereo"]
		is_interlaced = job_info["is_interlaced"]

		output_file = Path.cwd() / "output.mkv"

		if output_file.exists():
			output_file.unlink()

		job_start_time = time.time()

		downmixed_audio: Path = None

		if not has_stereo:
			# Extract audio to cwd/job_uuid-extracted-audio.mkv
			extracted_audio = Path.cwd() / f"{job_info['uuid']}-extracted-audio.mkv"
			if extracted_audio.exists():
				extracted_audio.unlink()
			extract_audio_command = ["ffmpeg", "-i", str(input_file), "-map", "-0:v?", "-map", "0:a:0", "-map", "-0:s?", "-c", "copy", str(extracted_audio)]
			self.__current_job_status["stage"] = "Extract Audio"
			# self.emit_current_job_status()
			stage_start_time = time.time()
			self._run_ffmpeg(extract_audio_command, self.update_job_status, stage_start_time, job_start_time, job_info)

			# Downmix audio
			downmixed_audio = Path.cwd() / f"{job_info['uuid']}-downmixed-audio.mkv"
			downmix_audio_command = ["ffmpeg", "-i", str(extracted_audio), "-map", "0:a", "-c", "aac", "-af", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", str(downmixed_audio)]
			self.__current_job_status["stage"] = "Downmix Audio"
			# self.emit_current_job_status()
			stage_start_time = time.time()
			self._run_ffmpeg(downmix_audio_command, self.update_job_status, stage_start_time, job_start_time, job_info)

		encode_inputs = ["-i", str(input_file)]
		tracks_to_copy = ["-map", "0:s?", "-map", "0:a?"]
		encoding_commands = []

		if downmixed_audio != None:
			encode_inputs.append("-i")
			encode_inputs.append(str(downmixed_audio))
			tracks_to_copy.append("-map")
			tracks_to_copy.append("1:a")

		if is_hevc:
			tracks_to_copy.append("-map")
			tracks_to_copy.append("0:v?")
		else:
			encoding_commands = ["-map", "0:v?", "-vcodec", "hevc"]

		final_ffmpeg_command = ["ffmpeg"] + encode_inputs + tracks_to_copy + encoding_commands + [str(output_file)]
		self.__current_job_status["stage"] = "Final encode"
		# self.emit_current_job_status()
		stage_start_time = time.time()
		self._run_ffmpeg(final_ffmpeg_command, self.update_job_status, stage_start_time, job_start_time, job_info)

		# Remove original file
		delete_successful = False
		if input_file.exists():
			try:
				input_file.unlink()
				delete_successful = True
			except PermissionError:
				print(f"Could not delete {input_file}")

		# Move output.mkv to take the original file's place
		if output_file.exists():
			if delete_successful:
				shutil_move(str(output_file), input_file.with_suffix(output_file.suffix))	# Retains the file suffix with the new name
			else:
				shutil_move(str(output_file), input_file.with_name(f"{input_file.stem}-RedCedarSmart").with_suffix(output_file.suffix))	# Retains the file suffix with the new name

		self.__current_job_status = {
			"percentage": None,
			"elapsed_time": None,
			"estimated_time_remaining": None,
			"fps": None,
			"stage_elapsed_time": None,
			"stage": None
		}
		# self.emit_current_job_status()

		if len(self.__waiting_jobs) > 0:
			next_job = self.__waiting_jobs.pop(0)
			self.__start_job(next_job)
		else:
			self.active = False

	def update_job_status(self, fps: float, current_file_time_str: str, current_speed: float, stage_start_time: float, job_start_time: float, job_info: Dict):
		total_length = [track for track in job_info["media_info"]["tracks"] if track["kind_of_stream"] == "General"][0]["duration"]
		total_length = float(f"{total_length[:-3]}.{total_length[-3:]}")

		reversed_times = current_file_time_str.split(":")[::-1] # Now in format [seconds, minutes, hours, etc]
		time_conversion_mask = [1.0, 60.0, 3600] # Multipliers to convert specific value to seconds
		current_file_time_seconds = sum([reversed_times[i] * time_conversion_mask[i] for i in range(reversed_times)])

		current_time = time.time()
		remaining_time = total_length - current_file_time_seconds

		if current_speed == 0:
			current_speed = 0.000000000000000000001

		self.__current_job_status = {
			"percentage": f"{round(current_file_time_seconds / total_length, 2)}%",
			"elapsed_time": chop_ms(timedelta(seconds=(current_time - job_start_time))),
			"estimated_time_remaining": chop_ms(timedelta(seconds=remaining_time / current_speed)),
			"fps": fps,
			"stage_elapsed_time": chop_ms(timedelta(seconds=(current_time - stage_start_time))),
			"stage": self.__current_job_status["stage"]
		}
		# self.emit_current_job_status()

		return (True, "")

def chop_ms(delta):
	return delta - timedelta(microseconds=delta.microseconds)
