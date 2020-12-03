import time
from pathlib import Path
from pymediainfo import MediaInfo
from typing import List
from uuid import uuid4

from .runner import JobRunner

class JobController:
	def __init__(self, path_to_search: Path=Path.cwd()) -> None:
		self.__path_to_search = path_to_search

		self.__file_system_check_offset = 15 * 60 # 15 minutes in seconds

		self.__last_file_system_check = 0
		self.__discoverer = None
		self.__dispatcher = None
		self.__verifier = None

		self.__job_queue = []

		self.current_job_status = {"uuid": None, "current_file": None, "percentage": None,
									"elapsed_time": None, "estimated_time": None,
									"current_fps": None, "average_fps": None}

		self.__runner = None

		self.__running = False

	def start(self) -> None:
		self.__runner = JobRunner()

		self.__running = True
		self.__run()

	def stop(self) -> None:
		self.__running = False

	def __run(self) -> None:
		while self.__running:
			# If last time file system was checked is greater than wait time: parse through all files
			if time.time() - self.__file_system_check_offset > self.__last_file_system_check:
				self.__last_file_system_check = time.time()
				video_file_paths = self.get_video_file_paths()
				# TODO: Verify that each video is not HDR, is not already in the Queue, and is not missing either HEVC or Stereo Audio
				for video_file in video_file_paths:
					media_info = MediaInfo.parse(str(video_file))
					if media_info.video_tracks[0].color_primaries == "BT.2020": # Is file HDR
						continue

					is_hevc = media_info.video_tracks[0].format == "HEVC"
					has_stereo = len([audio_track for audio_track in media_info.audio_tracks if audio_track.channel_s == 2]) > 0
					# Is not HEVC (we assume all HEVC is progressive) and video track scan type is not 'Progressive'
					is_interlaced = not is_hevc and media_info.video_tracks[0].scan_type != "Progressive"

					if is_hevc and has_stereo and not is_interlaced:
						continue

					if len([job for job in self.__job_queue if job["file"] == str(video_file)]):
						continue

					self.__job_queue.append({
						"uuid": str(uuid4()),
						"file": str(video_file),
						"is_hevc": is_hevc,
						"has_stereo": has_stereo,
						"is_interlaced": is_interlaced
					})

			if self.__runner.active == False and len(self.__job_queue) > 0:
				job_to_send = self.__job_queue.pop(0)
				self.__runner.new_job(job_to_send)

			self.current_job_status = self.__runner.get_job_status()

	def get_video_file_paths(self) -> List[Path]:
		video_file_types = [".m4v", ".mp4", ".mkv", ".avi", ".mov", ".webm", ".ogg", ".m4p", ".wmv", ".qt"]
		return [x for x in self.__path_to_search.glob("**/*") if x.is_file() and x.suffix in video_file_types]