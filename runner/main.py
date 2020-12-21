from logging import DEBUG, INFO, getLogger, WARNING, StreamHandler, FileHandler, Formatter
from os import getenv as os_getenv

from redcedar_runner import JobRunner

# Setup logging for main.py
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

log_level = DEBUG if os_getenv("REDCEDAR_DEBUG") == "True" else INFO

file_handler = FileHandler("/config/log.log")
file_handler.setLevel(log_level)
file_format = Formatter("%(asctime)s|%(name)s|%(levelname)s|%(lineno)d|%(message)s")
file_handler.setFormatter(file_format)

root_logger = getLogger()
root_logger.addHandler(file_handler)
root_logger.setLevel(log_level)

if __name__ == "__main__":
	controller_ip = os_getenv("REDCEDAR_RUNNER_CONTROLLER_IP")
	controller_port = os_getenv("REDCEDAR_RUNNER_CONTROLLER_PORT")

	if controller_ip == None:
		controller_ip = "localhost"
	if controller_port == None:
		controller_port = "5000"

	runner = JobRunner(controller_ip=f"{controller_ip}:{controller_port}")

	logger.info("Starting RedCedarRunner")
	runner.run()