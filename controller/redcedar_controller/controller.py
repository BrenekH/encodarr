import time
from collections import deque
from copy import deepcopy
from flask_socketio import SocketIO
from json import dump, load
from logging import getLogger, WARNING, StreamHandler, Formatter
from pathlib import Path
from pymediainfo import MediaInfo
from typing import Dict, List
from uuid import uuid4

# Setup logging for controller.py
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

class JobController:
	def __init__(self, socket_io: SocketIO, path_to_search: Path=Path.cwd(), config_dir: Path=Path("/config")) -> None:
		self.socket_io = socket_io
		self.__path_to_search = path_to_search

		self.__file_system_check_offset = 15 * 60 # 15 minutes in seconds

		self.__last_file_system_check = 0

		self.__job_queue = deque()

		self.__job_history = deque()

		# TODO: Save this on exit
		self.__dispatched_jobs = {}

		self.runner = None

		self.__running = False

		self.__history_file = config_dir / "history.json"

	def start(self) -> None:
		if not self.__history_file.exists():
			with self.__history_file.open("w") as f:
				f.write("{\"history\": []}")

		with self.__history_file.open() as f:
			for history_obj in load(f)["history"]:
				self.__job_history.appendleft(history_obj)

		self.__running = True
		self.__run()

	def get_new_job(self) -> Dict:
		while len(self.__job_queue) == 0:
			self.socket_io.sleep(0.5)

		to_send = self.__job_queue.popleft()
		self.__dispatched_jobs[to_send["uuid"]] = to_send
		self.emit_current_jobs()
		return to_send

	def update_job_status(self, status_info: Dict):
		# TODO: Make sure this works
		self.__dispatched_jobs[status_info["uuid"]]["status"] = status_info["status"]
		self.emit_current_jobs()

	def job_complete(self, history_entry: Dict):
		# TODO: Make sure this works
		del self.__dispatched_jobs[history_entry["uuid"]]
		self.__job_history.appendleft(history_entry["history"])
		self.__save_job_history()
		self.emit_current_jobs()

	def stop(self) -> None:
		logger.info("Stopping JobController")
		self.__running = False

	def get_job_history(self):
		return deepcopy(list(self.__job_history))

	def get_job_queue(self):
		return list(self.__job_queue)

	def __run(self) -> None:
		while self.__running:
			# If last time file system was checked is greater than wait time: parse through all files
			if time.time() - self.__file_system_check_offset > self.__last_file_system_check:
				logger.debug("Starting file system check")
				self.__last_file_system_check = time.time()
				video_file_paths = self.get_video_file_paths()
				for video_file in video_file_paths:
					# Have we already added this file?
					if len([job for job in self.get_job_queue() if job["file"] == str(video_file)]) > 0:
						continue

					# Have we already dispatched this file?
					if len([job for job in self.__dispatched_jobs if job["file"] == str(video_file)]) > 0:
						continue

					if "Plex Versions" in str(video_file): # Is the file 'optimized' by Plex?
						continue

					media_info = MediaInfo.parse(str(video_file))
					if not len(media_info.video_tracks) > 0:
						continue

					if media_info.video_tracks[0].color_primaries == "BT.2020": # Is the file HDR?
						continue

					is_hevc = media_info.video_tracks[0].format == "HEVC"
					has_stereo = len([audio_track for audio_track in media_info.audio_tracks if audio_track.channel_s == 2]) > 0
					# Is not HEVC (we assume all HEVC is progressive) and video track scan type is not 'Progressive'
					is_interlaced = not is_hevc and media_info.video_tracks[0].scan_type != "Progressive"

					if is_hevc and has_stereo and not is_interlaced:
						continue

					to_append = {
						"uuid": str(uuid4()),
						"file": str(video_file),
						"is_hevc": is_hevc,
						"has_stereo": has_stereo,
						"is_interlaced": is_interlaced,
						"media_info": media_info.to_data()
					}
					if "Movies" in str(video_file):
						self.__job_queue.append(to_append)
					else:
						self.__job_queue.appendleft(to_append)
					logger.info(f"Added to job queue: {to_append['file']}")
					self.socket_io.sleep(0.05)
				logger.debug("File system check complete")

			self.socket_io.sleep(0.075)

	def get_video_file_paths(self) -> List[Path]:
		video_file_types = [".m4v", ".mp4", ".mkv", ".avi", ".mov", ".webm", ".ogg", ".m4p", ".wmv", ".qt"]
		return [x for x in self.__path_to_search.glob("**/*") if x.is_file() and x.suffix in video_file_types]

	def __save_job_history(self):
		to_save = {"history": self.get_job_history()}
		with self.__history_file.open("w") as f:
			dump(to_save, f, indent=4)

	def emit_current_jobs(self):
		self.emit_event("current_jobs_update", self.__dispatched_jobs)

	def emit_event(self, event_name: str, data):
		logger.debug(f"Emitting event {event_name} with data: {data}")
		self.socket_io.emit(event_name, data, namespace="/updates")
