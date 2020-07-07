from colorama import Fore
from colorama import init as colorama_init
from json import load
from pathlib import Path
from sys import argv

video_file_types = [".m4v", ".mp4", ".mkv", ".avi", ".mov", ".webm", ".ogg", ".m4p", ".wmv", ".qt"]
file_suffix_buffer = 6
colorama_init()

def center(text: str, allowed_space: int) -> str:
	return ""

if len(argv) < 2:
	print(f"{Fore.RED}A file path to a 'completed_videos.json' must be passed{Fore.RESET}")
	exit()

target_path = Path(argv[1])

if not target_path.exists():
	print(f"{Fore.RED}{target_path} does not exist{Fore.RESET}")
	exit()
elif target_path.name != "completed_videos.json":
	print(f"{Fore.RED}{target_path} is not a 'completed_videos.json' file{Fore.RESET}")
	exit()

completed_videos_obj = load(open(target_path))

if "completed" not in completed_videos_obj:
	print(f"{Fore.RED}Could not find 'completed' key in {target_path}{Fore.RESET}")
	exit()
elif completed_videos_obj["completed"] == []:
	print(f"{Fore.YELLOW}{target_path} is empty{Fore.RESET}")
	exit()

completion_registry = {}
# completion_registry sample
{
	"path/to/video/that/exists": ".suffix",
	"path/to/video/that/does/not/exist": None
}

for path_str in completed_videos_obj["completed"]:
	for file_type in video_file_types:
		if Path(path_str + file_type).exists():
			completion_registry[path_str] = file_type
			break
	else:
		# Run if break never called in for loop
		completion_registry[path_str] = None

file_suffix_target_start = len(sorted(completed_videos_obj["completed"], key=len)[-1]) + file_suffix_buffer	# So that the [ .suffix ] are all in a straight line

for video_path in completion_registry:
	video_suffix = completion_registry[video_path]
	
	extra_spaces = file_suffix_target_start - len(video_path)

	if video_suffix == None: print(f"{video_path}{' ' * extra_spaces}[ {Fore.RED}Failed{Fore.RESET} ]")
	else: print(f"{video_path}{' ' * extra_spaces}[ {Fore.GREEN}{video_suffix}{Fore.RESET} ]")
