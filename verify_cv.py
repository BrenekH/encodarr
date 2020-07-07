from colorama import Fore
from colorama import init as colorama_init
from json import load
from pathlib import Path
from sys import argv

colorama_init()
video_file_types = [".m4v", ".mp4", ".mkv", ".avi", ".mov", ".webm", ".ogg", ".m4p", ".wmv", ".qt"]

if len(argv) < 2:
	print(f"{Fore.RED}A file path to a 'completed_videos.json' must be passed{Fore.RESET}")
	exit()

target_path = Path(argv[1])

if not target_path.exists():
	print(f"{Fore.RED}{target_path} does not exist{Fore.RESET}")
	exit()
elif target_path.name != "completed_videos.json":
	print(f"{Fore.RED}{target_path} is not a completed_videos.json{Fore.RESET}")
	exit()

completed_videos_obj = load(open(target_path))
