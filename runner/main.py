from logging import DEBUG, INFO, getLogger, WARNING, StreamHandler, FileHandler, Formatter
from os import getenv as os_getenv
from sys import argv

from encodarr_runner import JobRunner

# Setup logging for main.py
# Logging related env var setup
log_level = DEBUG if os_getenv("ENCODARR_DEBUG") == "True" or "--debug" in argv else INFO

temp = os_getenv("ENCODARR_LOG_FILE")

if temp != None:
	log_file = temp
else:
	log_file = "/config/runner.log"

# TODO: Add command-line argument for log file location

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

file_handler = FileHandler(log_file)
file_handler.setLevel(log_level)
file_format = Formatter("%(asctime)s|%(name)s|%(levelname)s|%(lineno)d|%(message)s")
file_handler.setFormatter(file_format)

root_logger = getLogger()
root_logger.addHandler(file_handler)
root_logger.setLevel(log_level)

if __name__ == "__main__":
	controller_ip = os_getenv("ENCODARR_RUNNER_CONTROLLER_IP")
	controller_port = os_getenv("ENCODARR_RUNNER_CONTROLLER_PORT")

	runner_name = os_getenv("ENCODARR_RUNNER_NAME")

	if controller_ip == None:
		controller_ip = "localhost"
	if controller_port == None:
		controller_port = "8123"

	if "--name" in argv:
		flag_index = argv.index("--name")
		try:
			runner_name = argv[flag_index + 1]
		except IndexError:
			raise RuntimeError("--name must be followed by another argument")

	if "--controller-ip" in argv:
		flag_index = argv.index("--controller-ip")
		try:
			controller_ip = argv[flag_index + 1]
		except IndexError:
			raise RuntimeError("--controller-ip must be followed by another argument")

	if "--controller-port" in argv:
		flag_index = argv.index("--controller-port")
		try:
			controller_port = argv[flag_index + 1]
		except IndexError:
			raise RuntimeError("--controller-port must be followed by another argument")

	runner = JobRunner(controller_ip=f"{controller_ip}:{controller_port}", runner_name=runner_name)

	logger.info("Starting Encodarr Runner")
	runner.run()
