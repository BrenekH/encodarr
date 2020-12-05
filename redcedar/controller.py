import time
from collections import deque
from flask_socketio import SocketIO
from pathlib import Path
from pymediainfo import MediaInfo
from typing import List
from uuid import uuid4

from .runner import JobRunner

class JobController:
	def __init__(self, socket_io: SocketIO, path_to_search: Path=Path.cwd(), config_dir: Path=Path("/config")) -> None:
		self.socket_io = socket_io
		self.__path_to_search = path_to_search

		self.__file_system_check_offset = 15 * 60 # 15 minutes in seconds

		self.__last_file_system_check = 0

		self.__job_queue = deque()

		self.__job_history = deque()

		self.current_job_status = {"uuid": None, "file": None, "percentage": None,
									"elapsed_time": None, "estimated_time": None,
									"current_fps": None, "average_fps": None}

		self.runner = None

		self.__running = False

	def start(self) -> None:
		self.runner = JobRunner(self.socket_io)

		self.__running = True
		self.__run()

	def stop(self) -> None:
		self.__running = False
		self.runner.stop()

	def get_job_history(self):
		return list(self.__job_history)

	def get_job_queue(self):
		return list(self.__job_queue)

	def __run(self) -> None:
		while self.__running:
			# If last time file system was checked is greater than wait time: parse through all files
			if time.time() - self.__file_system_check_offset > self.__last_file_system_check:
				self.__last_file_system_check = time.time()
				video_file_paths = self.get_video_file_paths()
				for video_file in video_file_paths:
					media_info = MediaInfo.parse(str(video_file))
					if not len(media_info.video_tracks) > 0:
						continue

					if media_info.video_tracks[0].color_primaries == "BT.2020" or "Plex Versions" in str(video_file): # Is the file HDR or 'optimized' by Plex
						continue

					is_hevc = media_info.video_tracks[0].format == "HEVC"
					has_stereo = len([audio_track for audio_track in media_info.audio_tracks if audio_track.channel_s == 2]) > 0
					# Is not HEVC (we assume all HEVC is progressive) and video track scan type is not 'Progressive'
					is_interlaced = not is_hevc and media_info.video_tracks[0].scan_type != "Progressive"

					if is_hevc and has_stereo and not is_interlaced:
						continue

					if len([job for job in self.get_job_queue() if job["file"] == str(video_file)]):
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
					# print(f"Added to job queue: {to_append['file']}")
					self.socket_io.sleep(0.05)

			if self.runner.active == False and len(self.__job_queue) > 0:
				job_to_send = self.__job_queue.popleft()
				# print(f"Sending new job: {to_append['file']}")
				self.runner.new_job(job_to_send)

			for job in self.runner.completed_jobs():
				self.__job_history.appendleft(job)

			self.socket_io.sleep(0.075)

	def get_video_file_paths(self) -> List[Path]:
		video_file_types = [".m4v", ".mp4", ".mkv", ".avi", ".mov", ".webm", ".ogg", ".m4p", ".wmv", ".qt"]
		return [x for x in self.__path_to_search.glob("**/*") if x.is_file() and x.suffix in video_file_types]
