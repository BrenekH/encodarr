import time
from _thread import start_new_thread
from collections import deque
from copy import deepcopy
from flask_socketio import SocketIO
from json import dump, load
from logging import getLogger, WARNING, StreamHandler, Formatter
from pathlib import Path
from pymediainfo import MediaInfo
from typing import Dict, List, Union
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
		self.health_check_interval = 60  # Seconds

		self.__last_file_system_check = 0

		self.__job_queue = deque()

		self.__job_history = deque()

		self.__dispatched_jobs = {}

		self.__unresponsive_jobs_uuids = []

		self.empty_status = {
			"fps": "N/A",
			"job_elapsed_time": "N/A",
			"percentage": "N/A",
			"stage": "Intermission",
			"stage_elapsed_time": "N/A",
			"stage_estimated_time_remaining": "N/A"
		}

		self.runner = None

		self.__running = False

		self.__history_file = config_dir / "history.json"
		self.__dispatched_jobs_file = config_dir / "dispatched_jobs.json"

	def start(self) -> None:
		if not self.__history_file.exists():
			with self.__history_file.open("w") as f:
				f.write("{\"history\": []}")

		with self.__history_file.open() as f:
			for history_obj in load(f)["history"]:
				self.__job_history.appendleft(history_obj)

		self.__load_dispatched_jobs()

		self.__running = True

		start_new_thread(self.health_check, ())
		self.__run()

	def get_new_job(self, runner_name="None") -> Dict:
		if len(self.__job_queue) == 0:
			while len(self.__job_queue) == 0:
				self.socket_io.sleep(0.5)

		to_send = self.__job_queue.popleft()
		self.__dispatched_jobs[to_send["uuid"]] = to_send
		self.__dispatched_jobs[to_send["uuid"]]["runner_name"] = runner_name
		self.__dispatched_jobs[to_send["uuid"]]["last_updated"] = time.time()
		self.__dispatched_jobs[to_send["uuid"]]["status"] = deepcopy(self.empty_status)
		logger.info(f"Dispatching job for {to_send['file']}")
		self.emit_current_jobs()
		return to_send

	def update_job_status(self, status_info: Dict) -> bool:
		"""Applies status information to the specified job

		Args:
			status_info (Dict): Status information to apply

		Returns:
			bool: Whether or not the status information was applied
		"""
		if status_info["uuid"] in self.__unresponsive_jobs_uuids:
			# Runner is responsive again but the job has already been added back into the queue so we ignore this Runner
			logger.warning(f"Received status information from previously unresponsive runner")
			return False

		self.__dispatched_jobs[status_info["uuid"]]["status"] = status_info["status"]
		self.__dispatched_jobs[status_info["uuid"]]["last_updated"] = time.time()

		logger.debug(f"Received status: {status_info}")
		self.emit_current_jobs()
		return True

	def job_complete(self, history_entry: Dict) -> bool:
		"""Marks a running job as complete and saves the supplied history information

		Args:
			history_entry (Dict): The history information to save

		Returns:
			bool: Whether or not the operation was completed
		"""
		if history_entry["uuid"] in self.__unresponsive_jobs_uuids:
			# Runner is responsive again but the job has already been added back into the queue so we ignore this Runner
			logger.warning(f"Received job complete signal from previously unresponsive runner")
			return False

		del self.__dispatched_jobs[history_entry["uuid"]]
		self.__job_history.appendleft(history_entry["history"])
		self.__save_job_history()

		logger.info(f"Received job complete for {history_entry['history']['file']}")
		self.emit_current_jobs()
		return True

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
					if len([job for job in self.__dispatched_jobs if self.__dispatched_jobs[job]["file"] == str(video_file)]) > 0:
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

	def __save_dispatched_jobs(self):
		to_save = {"jobs": self.__dispatched_jobs}
		with self.__dispatched_jobs_file.open("w") as f:
			dump(to_save, f, indent=4)

	def __load_dispatched_jobs(self):
		if not self.__dispatched_jobs_file.exists():
			self.__save_dispatched_jobs()

		with self.__dispatched_jobs_file.open("r") as f:
			self.__dispatched_jobs = load(f)["jobs"]

	def emit_current_jobs(self):
		filtered_dict = {}

		for key in self.__dispatched_jobs:
			filtered_dict[key] = {
				"file": self.__dispatched_jobs[key]["file"],
				"runner_name": self.__dispatched_jobs[key]["runner_name"],
				"status": self.__dispatched_jobs[key]["status"]
			}

		self.emit_event("current_jobs_update", filtered_dict)

	def emit_event(self, event_name: str, data):
		logger.debug(f"Emitting event {event_name} with data: {data}")
		self.socket_io.emit(event_name, data, namespace="/updates")

	def health_check(self):
		while self.__running:
			self.micro_sleep(self.health_check_interval)
			logger.debug("Starting Health Check")
			keys_to_del = []
			for key in self.__dispatched_jobs:
				#? Maybe use two different settings for interval vs time dead?
				if self.__dispatched_jobs[key]["last_updated"] < time.time() - self.health_check_interval:
					# Runner is unresponsive
					logger.warning(f"Runner {self.__dispatched_jobs[key]['runner_name']} is unresponsive. Adding its job back into the queue.")
					to_append = {
						"uuid": str(uuid4()),
						"file": deepcopy(self.__dispatched_jobs[key]["file"]),
						"is_hevc": deepcopy(self.__dispatched_jobs[key]["is_hevc"]),
						"has_stereo": deepcopy(self.__dispatched_jobs[key]["has_stereo"]),
						"is_interlaced": deepcopy(self.__dispatched_jobs[key]["is_interlaced"]),
						"media_info": deepcopy(self.__dispatched_jobs[key]["media_info"])
					}
					keys_to_del.append(key)
					self.__job_queue.appendleft(to_append)
					self.__unresponsive_jobs_uuids.append(key)

			for key in keys_to_del:
				del self.__dispatched_jobs[key]

			if len(keys_to_del) > 0:
				self.emit_current_jobs()
			logger.debug("Health Check complete")

	def micro_sleep(self, seconds: Union[int, float]):
		self.socket_io.sleep(seconds - int(seconds)) # Complete any sub-second sleeping

		for _ in range(int(seconds)):
			for _ in range(10):
				if not self.__running:
					return
				self.socket_io.sleep(0.1)
