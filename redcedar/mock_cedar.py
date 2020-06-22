import random, time
from copy import deepcopy
from datetime import datetime, timedelta
from flask_socketio import SocketIO, emit
from pathlib import Path
from .db import db_video_files

class MockCedar:
	def __init__(self, socketio: SocketIO, current_working_directory: Path=Path.cwd()):
		self.cwd = current_working_directory
		self.socket_io = socketio

		# External Triggers
		self.__stop = False
		self.__new_connection = False
		
		self.video_file_paths = []

		self.total_start_time = 0.0
		self.current_start_time = 0.0

		self._completed_videos = {"completed": []}

		self._skipped_videos = []
		self._completed_videos_events = []

		self._plex_version_skip_text = "Found \"Plex Versions\" in the file path"
		self._previously_transcoded_skip_text = "File was previously transcoded"

	def run(self, path_to_search: Path=None):
		if path_to_search == None:
			path_to_search = self.cwd

		self.video_file_paths = [Path(x) for x in db_video_files]

		self.total_start_time = time.time()

		for indx, path in enumerate(self.video_file_paths, start=1):
			self.current_start_time = time.time()
			self.emit_current_file_update(f"{indx}/{len(self.video_file_paths)}", path)

			if "Plex Versions" in str(path):
				self.emit_file_skip(path, self._plex_version_skip_text)
				self._skipped_videos.append({
					"file_path": str(path),
					"reason": self._plex_version_skip_text,
					"timestamp": self.get_friendly_timestamp()
				})
				continue
		
			elif random.randint(0, 20) == 5:
				self.emit_file_skip(path, self._previously_transcoded_skip_text)
				self._skipped_videos.append({
					"file_path": str(path),
					"reason": self._previously_transcoded_skip_text,
					"timestamp": self.get_friendly_timestamp()
				})
				continue

			latest_avg_fps = 0.00
			for x in range(25):
				self.socket_io.sleep(random.randint(1, 2))

				current_time = time.time()
				latest_avg_fps = round(random.uniform(0.1, 15.0), 3)
				self.emit_current_file_status_update(chop_ms(timedelta(seconds=(current_time - self.total_start_time))),
													chop_ms(timedelta(seconds=(random.uniform(current_time, current_time * 1.2) - self.current_start_time))),
													chop_ms(timedelta(seconds=(current_time - self.current_start_time))),
													f"{random.randint(0, 100)}%",
													round(random.uniform(0.1, 15.0), 3),
													latest_avg_fps)

				if self.__new_connection:
					self.emit_connect_info(f"{indx}/{len(self.video_file_paths)}", path)
					self.__new_connection = False

				if self.__stop:
					break

			if self.__stop:
				break

			time_taken = chop_ms(timedelta(seconds=(time.time() - self.current_start_time)))
			self.emit_file_complete(path, latest_avg_fps, time_taken)
			self._completed_videos_events.append({
				"file_path": str(path),
				"avg_fps": latest_avg_fps,
				"time_taken": str(time_taken),
				"timestamp": self.get_friendly_timestamp()
			})
		
		end_time = time.time()

		while not self.__stop:
			# Check for new connections and emit new_connection event when there are
			if self.__new_connection:
				self.emit_connect_info(f"{len(self.video_file_paths)}/{len(self.video_file_paths)}", "Operation Complete")
				self.emit_current_file_status_update(chop_ms(timedelta(seconds=(end_time - self.total_start_time))), "0:00:00", "0:00:00", "", "0.00", "0.00")

	def stop(self):
		self.__stop = True

	def new_connection(self):
		self.__new_connection = True

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
