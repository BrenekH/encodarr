import sys, time
from datetime import timedelta
from flask_socketio import SocketIO, emit
from json import dump, load, loads
from json.decoder import JSONDecodeError
from pathlib import Path
from subprocess import DEVNULL, PIPE, Popen
from typing import List, Tuple

class RedCedar:
	def __init__(self, socketio: SocketIO, current_working_directory: Path=Path.cwd()):
		"""Object to manage and run RedCedar operations.

		Args:
			current_working_directory (Path, optional): The directory that RedCedar uses for the output file and the memory file. Defaults to Path.cwd().
		"""
		self.cwd = current_working_directory
		self.socket_io = socketio
		
		self.video_file_paths = []

		self.total_start_time = 0.0
		self.current_start_time = 0.0

		self._completed_videos = {"completed": []}
		self._completed_videos_path = self.cwd / "completed_videos.json"

		if self._completed_videos_path.exists():
			self._completed_videos = load(open(self._completed_videos_path))
		else:
			self.save_completed_videos_json()

		self.output_file = self.cwd / "output.m4v"

	def run(self, path_to_search: Path=None):
		if path_to_search == None:
			path_to_search = self.cwd

		# Make sure output.m4v doesn't exist
		if sys.version_info[1] >= 8:
			self.output_file.unlink(missing_ok=True)
		else:
			if self.output_file.exists():
				self.output_file.unlink()

		self.video_file_paths = self.get_video_file_paths(path_to_search)

		self.total_start_time = time.time()

		for indx, path in enumerate(self.video_file_paths, start=1):	# enumerate returns (the index, the value) on every loop
			# print(f"{Fore.GREEN}Starting {self.path_color}{path}{Fore.RESET}")
			self.current_start_time = time.time()

			# Make sure file isn't an 'optimized version' by Plex
			if "Plex Versions" in str(path):
				# print(f"{Fore.YELLOW}Skipping {self.path_color}{path}{Fore.YELLOW} because it contains '{Fore.BLUE}Plex Versions{Fore.YELLOW}'{Fore.RESET}")
				continue

			# Make sure we haven't already encoded this file previously
			if self.check_video_complete(path):
				# print(f"{Fore.YELLOW}Skipping {self.path_color}{path}{Fore.YELLOW} because it is marked in {Fore.BLUE}completed_videos.json{Fore.YELLOW} as already transcoded.{Fore.RESET}")
				continue

			if self.output_file.exists():
				self.output_file.unlink()

			# Run handbrake cli and save to output.m4v
			handbrake_command = f"HandBrakeCLI -i \"{path}\" -o output.m4v -e x265 --optimize --json"
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
							self.output_from_json(json_string, indx)
						record_json, json_string = (False, "")
			
			self.printer.finish()	# Ensures the next print statement doesn't overwrite the progress line

			# Remove original file
			if sys.version_info[1] >= 8:
				path.unlink(missing_ok=True)
			else:
				if path.exists():
					path.unlink()

			# Move output.m4v to take the original file's place
			if self.output_file.exists():
				self.output_file.rename(path.with_suffix(self.output_file.suffix))	# Retains the .m4v suffix with the new name

			self.mark_video_complete(path)
			# print(f"{Fore.MAGENTA}Conversion complete{Fore.RESET} for {self.path_color}{path}{Fore.RESET}.")
		
		print("RedCedar Complete")

	def output_from_json(self, json_string, job_number=1) -> Tuple[bool, str]:
		"""Outputs data from the given json_string

		Args:
			json_string (str): The string to get data from

		Returns:
			Tuple[bool, str]: If the function outputted and why the function did not output if the first option is False.
		"""
		sanitized_json_string = json_string.replace('Progress: ', '')
		json_obj = {}
		try:
			json_obj = loads(sanitized_json_string)
		except JSONDecodeError as e:
			# TODO: Log error message to a file
			error_message = f"json decode error: {e} on {sanitized_json_string}"
			return (False, error_message)

		# print(f"BEGIN\n\n{json_obj}\nEND")
		# Output format: Time: {}; File: 1 / 100; Current ETA: {}; Current Time: {}; {}%;
		if json_obj["State"] != "WORKING":
			return (False, "HandBrakeCLI is not in State: WORKING")
		working = json_obj["Working"]
		current_time = time.time()
		# self.printer.output(f"Total Time: {timedelta(seconds=(current_time - self.total_start_time))}; File: {job_number}/{len(self.video_file_paths)}; Current ETA: {timedelta(seconds=working['ETASeconds'])}; Current Time: {timedelta(seconds=(current_time - self.current_start_time))}; {round(working['Progress'] * 100, 2)}%; FPS: {round(working['Rate'], 3)}; Avg FPS: {round(working['RateAvg'], 3)}")
		self.broadcast_status_update(timedelta(seconds=(current_time - self.total_start_time)), f"{job_number}/{len(self.video_file_paths)}", timedelta(seconds=working['ETASeconds']), timedelta(seconds=(current_time - self.current_start_time)), f"{round(working['Progress'] * 100, 2)}%", round(working['Rate'], 3), round(working['RateAvg'], 3))

	def get_video_file_paths(self, top_path: Path) -> List[Path]:
		video_file_types = [".m4v", ".mp4", ".mkv", ".avi"]
		return [x for x in top_path.glob("**/*") if x.is_file() and x.suffix in video_file_types]

	def save_completed_videos_json(self):
		dump(self._completed_videos, open(self._completed_videos_path, "w"), indent=4)

	def mark_video_complete(self, path: Path):
		self._completed_videos["completed"].append(str(path).replace(path.suffix, ""))
		self.save_completed_videos_json()

	def check_video_complete(self, path: Path) -> bool:
		return (str(path).replace(path.suffix, "") in self._completed_videos["completed"])

	def broadcast_status_update(self, total_time, file_progress, current_eta, current_time, percentage, current_fps, avg_fps):
		self.socket_io.emit("status_update", {
												"total_time": total_time,
												"file_progress": file_progress,
												"current_eta": current_eta,
												"current_time": current_time,
												"percentage": percentage,
												"current_fps": current_fps,
												"avg_fps": avg_fps
											}, namespace="/websocket")
