import time
from copy import copy
from datetime import timedelta
from flask_socketio import SocketIO
from json import JSONDecodeError, loads
from pathlib import Path
from shutil import move as shutil_move
from typing import Dict, Tuple
from subprocess import DEVNULL, PIPE, Popen

class JobRunner:
	def __init__(self, socket_io: SocketIO):
		self.socket_io = socket_io
		self.active = False
		self.__current_job_status = {
			"percentage": None,
			"elapsed_time": None,
			"estimated_time": None,
			"current_fps": None,
			"average_fps": None
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
		# Calculates the handbrake options starts the background HandBrake task
		pass

	def _run_job(self, input_file: Path, handbrake_options: str):
		output_file = Path.cwd() / "output.m4v"

		if output_file.exists():
			output_file.unlink()

		job_start_time = time.time()

		# Run handbrake cli and save to output.m4v
		handbrake_command = f"HandBrakeCLI -i \"{input_file}\" -o \"{output_file}\" {handbrake_options} --json"
		with Popen([handbrake_command], stdout=PIPE, stderr=DEVNULL, bufsize=1, shell=True, universal_newlines=True) as p:
			record_json, json_string = (False, "")
			for line in p.stdout:
				if record_json:
					json_string += line

				if line.startswith("Progress: {"):
					record_json = True
					json_string += line

				elif line.startswith("}"):
					if not json_string.strip() == "":
						self.update_job_status_from_json(json_string, job_start_time)
					record_json, json_string = (False, "")

				# Sleep to allow for other operations
				self.socket_io.sleep(0.01)

		# Remove original file
		delete_successful = False
		if input_file.exists():
			try:
				input_file.unlink()
				delete_successful = True
			except PermissionError:
				print(f"Could not delete {input_file}")

		# Move output.m4v to take the original file's place
		if output_file.exists():
			if delete_successful:
				shutil_move(str(output_file), input_file.with_suffix(output_file.suffix))	# Retains the .m4v suffix with the new name
			else:
				shutil_move(str(output_file), input_file.with_name(f"{input_file.stem}-New H.265 Encoded").with_suffix(output_file.suffix))	# Retains the .m4v suffix with the new name

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
			"average_fps": round(working['RateAvg'], 3)
		}

		return (True, "")

def chop_ms(delta):
	return delta - timedelta(microseconds=delta.microseconds)
