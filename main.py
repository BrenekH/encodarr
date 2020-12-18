from datetime import datetime
from flask.helpers import send_file
from flask_socketio import SocketIO
from flask import abort, Flask, render_template, request, make_response
from json import dumps
from logging import INFO, getLogger, ERROR, WARNING, StreamHandler, FileHandler, Formatter
from pathlib import Path
from sys import argv
from redcedar import RedCedar
from redcedar.mock_cedar import MockCedar

app = Flask(__name__)
app.config["SECRET_KEY"] = "my_secret"
app.config["DEBUG"] = False

# Get rid of unnecessary logs
getLogger("werkzeug").setLevel(ERROR)
getLogger("socketio.server").setLevel(WARNING)
getLogger("engineio.server").setLevel(WARNING)

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

file_handler = FileHandler("/config/log.log")
file_handler.setLevel(INFO)
file_format = Formatter("%(asctime)s|%(name)s|%(levelname)s|%(lineno)d|%(message)s")
file_handler.setFormatter(file_format)

root_logger = getLogger()
root_logger.addHandler(file_handler)
root_logger.setLevel(INFO)

# Turn the flask app into a socketio app
socketio = SocketIO(app, async_mode=None, logger=False, engineio_logger=False)
redcedar_obj = None

def run_redcedar():
	global redcedar_obj
	redcedar_obj = RedCedar(socketio, Path("/usr/app/tosearch"))
	redcedar_obj.start()

def run_redcedar_cwd():
	global redcedar_obj
	redcedar_obj = RedCedar(socketio)
	redcedar_obj.start()

def run_mockcedar():
	global redcedar_obj
	redcedar_obj = MockCedar(socketio, Path("/usr/app/tosearch"))
	redcedar_obj.run()

@app.route("/")
def index():
	# Only by sending this page first will the client be connected to the socketio instance
	return render_template("index.html")

@app.route("/favicon.ico")
def favicon_ico():
	send_file("static/favicon/favicon.ico")

@app.route("/api/v1/queue", methods=["GET"])
def api_v1_queue():
	if request.method != "GET":
		abort(405)

	if redcedar_obj == None:
		abort(500)

	response_dict = {"queue": []}
	for entry in redcedar_obj.get_job_queue():
		response_dict["queue"].append({"filename": entry["file"],
										"video_op": not entry["is_hevc"],
										"audio_op": not entry["has_stereo"]
										})
	response = make_response(dumps(response_dict))
	response.status_code = 200
	response.headers["content-type"] = "application/json"

	return response

@app.route("/api/v1/history", methods=["GET"])
def api_v1_history():
	if request.method != "GET":
		abort(405)

	if redcedar_obj == None:
		abort(500)

	history_to_send = []
	for job in redcedar_obj.get_job_history():
		job["datetime_completed"] = datetime.utcfromtimestamp(job["datetime_completed"]).strftime("%m-%d-%Y %H:%M:%S")
		history_to_send.append(job)

	response = make_response(dumps({"history": history_to_send}))
	response.status_code = 200
	response.headers["content-type"] = "application/json"

	return response

@socketio.on("connect", namespace="/updates")
def test_connect():
	if redcedar_obj != None:
		redcedar_obj.runner.emit_current_job()
		redcedar_obj.runner.emit_current_job_status()
	logger.info("Client connected")

@socketio.on("disconnect", namespace="/updates")
def test_disconnect():
	logger.info("Client disconnected")

if __name__ == "__main__":
	if "mockcedar" in argv:
		logger.info("Running with mock RedCedar background process")
		socketio.start_background_task(run_mockcedar)
	elif "cwd" in argv:
		logger.info("Running RedCedar in current working directory")
		socketio.start_background_task(run_redcedar_cwd)
	elif "noredcedar" in argv:
		logger.info("Running without RedCedar background process")
	elif "logtree" in argv:
		import logging_tree
		logging_tree.printout()
	else:
		logger.info("Starting RedCedar")
		socketio.start_background_task(run_redcedar)

	debug_mode = False
	if "flaskdebug" in argv:
		debug_mode = True

	socketio.run(app, host="0.0.0.0", debug=debug_mode)

	logger.info("Stopping Project RedCedar")
	if redcedar_obj != None:
		redcedar_obj.stop()
