from logging import DEBUG, INFO, getLogger, WARNING, StreamHandler, FileHandler, Formatter
from os import getenv as os_getenv
from random import randint

from redcedar_runner import JobRunner

# Setup logging for main.py
# Logging related env var setup
log_level = DEBUG if os_getenv("REDCEDAR_DEBUG") == "True" else INFO

temp = os_getenv("REDCEDAR_LOG_FILE")

if temp != None:
	log_file = temp
else:
	log_file = "/config/runner.log"

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

# Non-logging related env var setup
raw_runner_name = os_getenv("REDCEDAR_RUNNER_NAME")

if raw_runner_name == None:
	runner_name = f"Runner-{str(randint(1, 999)).rjust(3, '0')}"
	logger.warning(f"No runner name was set. Using {runner_name}.")
else:
	runner_name = raw_runner_name

if __name__ == "__main__":
	controller_ip = os_getenv("REDCEDAR_RUNNER_CONTROLLER_IP")
	controller_port = os_getenv("REDCEDAR_RUNNER_CONTROLLER_PORT")

	if controller_ip == None:
		controller_ip = "localhost"
	if controller_port == None:
		controller_port = "5000"

	runner = JobRunner(controller_ip=f"{controller_ip}:{controller_port}", runner_name=runner_name)

	logger.info("Starting RedCedarRunner")
	runner.run()
