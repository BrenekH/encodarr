import sys, time
from datetime import datetime, timedelta
from flask_socketio import SocketIO, emit
from json import dump, load, loads
from json.decoder import JSONDecodeError
from pathlib import Path
from shutil import move as shutil_move
from subprocess import DEVNULL, PIPE, Popen
from typing import List, Tuple

class RedCedar:
	def __init__(self, socketio: SocketIO, current_working_directory: Path=Path.cwd()):
		"""Class to manage and run RedCedar operations.

		Args:
			current_working_directory (Path, optional): The directory that RedCedar uses for the output file and the memory file. Defaults to Path.cwd().
		"""
		self.cwd = current_working_directory
		self.socket_io = socketio

		self.__stop = False
		self.__new_connection = False
		
		self.video_file_paths = []

		self.total_start_time = 0.0
		self.current_start_time = 0.0

		self._completed_videos = {"completed": []}
		self._completed_videos_path = self.cwd / "completed_videos.json"

		if self._completed_videos_path.exists():
			self._completed_videos = load(open(self._completed_videos_path))
		else:
			self.save_completed_videos_json()

		self._skipped_videos = []
		self._completed_videos_events = []

		self._plex_version_skip_text = "Found \"Plex Versions\" in the file path"
		self._previously_transcoded_skip_text = "File was previously transcoded"

		self.output_file = Path.cwd() / "output.m4v"

		self.__latest_avg_fps = 0.0

	def run(self, path_to_search: Path=None):
		if path_to_search == None:
			path_to_search = self.cwd

		# Make sure output.m4v doesn't exist
		if self.output_file.exists():
			self.output_file.unlink()

		self.video_file_paths = self.get_video_file_paths(path_to_search)

		self.total_start_time = time.time()

		for indx, path in enumerate(self.video_file_paths, start=1):	# enumerate returns (the index, the value) on every loop
			self.current_start_time = time.time()
			self.emit_current_file_update(f"{indx}/{len(self.video_file_paths)}", path)

			# Make sure file isn't an 'optimized version' by Plex
			if "Plex Versions" in str(path):
				self.emit_file_skip(path, self._plex_version_skip_text)
				self._skipped_videos.append({
					"file_path": str(path),
					"reason": self._plex_version_skip_text,
					"timestamp": self.get_friendly_timestamp()
				})
				continue

			# Make sure we haven't already encoded this file previously
			elif self.check_video_complete(path):
				self.emit_file_skip(path, self._previously_transcoded_skip_text)
				self._skipped_videos.append({
					"file_path": str(path),
					"reason": self._previously_transcoded_skip_text,
					"timestamp": self.get_friendly_timestamp()
				})
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
					
					# Check for new connections and emit new_connection event when there are
					if self.__new_connection:
						self.emit_connect_info(f"{indx}/{len(self.video_file_paths)}", path)
					
					# Sleep to allow for other operations
					self.socket_io.sleep(0.01)

			# Remove original file
			delete_successful = False
			if path.exists():
				try:	
					path.unlink()
					delete_successful = True
				except PermissionError:
					print(f"Could not delete {path}")

			# Move output.m4v to take the original file's place
			if self.output_file.exists():
				if delete_successful:
					shutil_move(self.output_file, path.with_suffix(self.output_file.suffix))	# Retains the .m4v suffix with the new name
				else:
					shutil_move(self.output_file, path.with_name(f"{path.stem}-New H.265 Encoded").with_suffix(self.output_file.suffix))	# Retains the .m4v suffix with the new name

			self.mark_video_complete(path)

			time_taken = chop_ms(timedelta(seconds=(time.time() - self.current_start_time)))
			self.emit_file_complete(path, self.__latest_avg_fps, time_taken)
			self._completed_videos_events.append({
				"file_path": str(path),
				"avg_fps": self.__latest_avg_fps,
				"time_taken": str(time_taken),
				"timestamp": self.get_friendly_timestamp()
			})

		end_time = time.time()

		while not self.__stop:
			# Check for new connections and emit new_connection event when there are
			if self.__new_connection:
				self.emit_connect_info(f"{len(self.video_file_paths)}/{len(self.video_file_paths)}", "Operation Complete")
				self.emit_current_file_status_update(chop_ms(timedelta(seconds=(end_time - self.total_start_time))), "0:00:00", "0:00:00", "", "0.00", "0.00")

	def output_from_json(self, json_string, job_number=1) -> Tuple[bool, object]:
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
		
		self.emit_current_file_status_update(chop_ms(timedelta(seconds=(current_time - self.total_start_time))),
											chop_ms(timedelta(seconds=working['ETASeconds'])),
											chop_ms(timedelta(seconds=(current_time - self.current_start_time))),
											f"{round(working['Progress'] * 100, 2)}%",
											round(working['Rate'], 3),
											round(working['RateAvg'], 3))

		self.__latest_avg_fps = round(working["RateAvg"], 3)

		return (True, "")

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

	def new_connection(self):
		self.__new_connection = True

	def stop(self):
		self.__stop = True

	def emit_file_skip(self, path: Path, reason: str):
		self.emit_event("file_skip", {"file_path": str(path),
									"reason": reason,
									"timestamp": self.get_friendly_timestamp()})

	def emit_file_complete(self, path: Path, avg_fps, time_taken):
		self.emit_event("file_complete", {"file_path": str(path),
										"avg_fps": str(avg_fps),
										"time_taken": str(time_taken),
										"timestamp": self.get_friendly_timestamp()})

	def emit_connect_info(self, file_count: str, current_file_path: Path):
		self.emit_event("connect_info", {"file_count": file_count,
										"current_file_path": str(current_file_path),
										"completed_files": self._completed_videos_events[::-1],
										"skipped_files": self._skipped_videos[::-1]})

	def emit_current_file_status_update(self, total_time, current_etr, current_time, percentage, current_fps, avg_fps):
		self.emit_event("current_file_status_update", {"total_time": str(total_time),
													"current_etr": str(current_etr),
													"current_time": str(current_time),
													"percentage": str(percentage),
													"current_fps": str(current_fps),
													"avg_fps": str(avg_fps)})

	def emit_current_file_update(self, file_count: str, path: Path):
		self.emit_event("current_file_update", {"file_count": file_count, "file_path": str(path)})

	def emit_event(self, event_name: str, data):
		self.socket_io.emit(event_name, data, namespace="/websocket")

	def get_friendly_timestamp(self) -> str:
		return datetime.now().strftime("%m/%d/%Y, %H:%M:%S")

def chop_ms(delta):
	return delta - timedelta(microseconds=delta.microseconds)