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
			"estimated_time": None,
			"current_fps": None,
			"average_fps": None,
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

	def _run_ffmpeg(self, ffmpeg_command: List[str], interval_callback: callable=None, *args):
		if type(ffmpeg_command) != list:
			raise TypeError("ffmpeg_command argument must be a list of strings")

		with Popen(ffmpeg_command, stdout=PIPE, stderr=STDOUT, universal_newlines=True) as p:
			for line in p.stdout:
				print(line)
				# TODO: Parse the line (current_frame (may equal None), fps, time, speed)
				if interval_callback != None:
					interval_callback("current_frame(maybe)", "fps", "time", "speed", *args)

				self.socket_io.sleep(0.1)

	def _run_job(self, job_info: Dict):
		input_file = Path(job_info["file"])
		framerate = job_info["media_info"]["frame_rate"]
		is_hevc = job_info["is_hevc"]
		has_stereo = job_info["has_stereo"]
		is_interlaced = job_info["is_interlaced"]

		output_file = Path.cwd() / "output.mkv"

		if output_file.exists():
			output_file.unlink()

		job_start_time = time.time()

		encode_only = ["ffmpeg", "-i", str(input_file), "-map", "0:a?", "-map", "0:s?", "-c", "copy", "-map", "0:v", "-vcodec", "hevc", str(output_file)]
		downmixed_audio: Path = None

		if not has_stereo:
			# Extract audio to cwd/job_uuid-extracted-audio.mkv
			extracted_audio = Path.cwd() / f"{job_info['uuid']}-extracted-audio.mkv"
			if extracted_audio.exists():
				extracted_audio.unlink()
			extract_audio_command = ["ffmpeg", "-i", str(input_file), "-map", "-0:v?", "-map", "0:a:0", "-map", "-0:s?", "-c", "copy", str(extracted_audio)]
			self.__current_job_status["stage"] = "Extracting Audio"
			self._run_ffmpeg(extract_audio_command) # TODO: Add interval callback function

			# Downmix audio
			downmixed_audio = Path.cwd() / f"{job_info['uuid']}-downmixed-audio.mkv"
			downmix_audio_command = ["ffmpeg", "-i", str(extracted_audio), "-map", "0:a", "-c", "aac", "-af", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", str(downmixed_audio)]
			self.__current_job_status["stage"] = "Downmixing Audio"
			self._run_ffmpeg(downmix_audio_command) # TODO: Add interval callback function

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
		self._run_ffmpeg(final_ffmpeg_command) # TODO: Add interval callback function

		# Remove original file
		delete_successful = False
		if input_file.exists():
			try:
				input_file.unlink()
				delete_successful = True
			except PermissionError:
				print(f"Could not delete {input_file}")

		delete_successful = False #! TODO: Remove!

		# Move output.mkv to take the original file's place
		if output_file.exists():
			if delete_successful:
				shutil_move(str(output_file), input_file.with_suffix(output_file.suffix))	# Retains the file suffix with the new name
			else:
				shutil_move(str(output_file), input_file.with_name(f"{input_file.stem}-RedCedarSmart").with_suffix(output_file.suffix))	# Retains the file suffix with the new name

		if len(self.__waiting_jobs) > 0:
			next_job = self.__waiting_jobs.pop(0)
			self.__start_job(next_job)
		else:
			self.active = False

	def update_job_status_from_json(self, json_string: str, job_start_time: float) -> Tuple[bool, object]:
		"""Outputs data from the given json_string

		Args:
			json_string (str): The string to get data from

		Returns:
			Tuple[bool, object]: If the function outputted and why the function did not output if the first option is False.
		"""
		sanitized_json_string = json_string.replace('Progress: ', '')
		json_obj = {}
		try:
			json_obj = loads(sanitized_json_string)
		except JSONDecodeError as e:
			# TODO: Log error message to a file
			error_message = f"json decode error: {e} on {sanitized_json_string}"
			return (False, error_message)

		if json_obj["State"] != "WORKING":
			return (False, "HandBrakeCLI is not in State: WORKING")

		working = json_obj["Working"]
		current_time = time.time()

		self.__current_job_status = {
			"percentage": f"{round(working['Progress'] * 100, 2)}%",
			"elapsed_time": chop_ms(timedelta(seconds=(current_time - job_start_time))),
			"estimated_time": chop_ms(timedelta(seconds=working['ETASeconds'])),
			"current_fps": round(working['Rate'], 3),
			"average_fps": round(working['RateAvg'], 3),
			"stage": "None"
		}

		return (True, "")

def chop_ms(delta):
	return delta - timedelta(microseconds=delta.microseconds)
