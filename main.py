from flask_socketio import SocketIO, emit
from flask import Flask, render_template, url_for, copy_current_request_context
from logging import getLogger, ERROR
from pathlib import Path
from random import random
from sys import argv
from threading import Thread, Event
from time import sleep
from redcedar import RedCedar, RedCedarSmart
from redcedar.mock_cedar import MockCedar

app = Flask(__name__)
app.config['SECRET_KEY'] = 'my_secret'
app.config['DEBUG'] = False

getLogger('werkzeug').setLevel(ERROR)

#turn the flask app into a socketio app
socketio = SocketIO(app, async_mode=None, logger=False, engineio_logger=False)
redcedar_obj = None

def run_redcedar():
	global redcedar_obj
	redcedar_obj = RedCedar(socketio, Path("/usr/app/tosearch"))
	redcedar_obj.run()

def run_redcedar_cwd():
	global redcedar_obj
	redcedar_obj = RedCedar(socketio)
	redcedar_obj.run()

def run_redcedar_smart():
	global redcedar_obj
	redcedar_obj = RedCedarSmart(socketio, Path("D:\Videos\RedCedarSmartTestEnv"))
	redcedar_obj.start()

def run_redcedar_smart_cwd():
	global redcedar_obj
	redcedar_obj = RedCedarSmart(socketio)
	redcedar_obj.start()

def run_mockcedar():
	global redcedar_obj
	redcedar_obj = MockCedar(socketio, Path("/usr/app/tosearch"))
	redcedar_obj.run()

@app.route('/')
def index():
	#only by sending this page first will the client be connected to the socketio instance
	return render_template('index.html')

@app.route("/new")
def new_index():
	return render_template("newindex.html")

@socketio.on('connect', namespace='/updates')
def test_connect():
	if redcedar_obj != None:
		redcedar_obj.runner.emit_current_job()
		redcedar_obj.runner.emit_current_job_status()
	print('Client connected')

@socketio.on('disconnect', namespace='/updates')
def test_disconnect():
	print('Client disconnected')

if __name__ == '__main__':
	if "mockcedar" in argv:
		print("Running with mock RedCedar background process")
		socketio.start_background_task(run_mockcedar)
	elif "cwd" in argv:
		print("Running redcedar in current working directory")
		socketio.start_background_task(run_redcedar_cwd)
	elif "smart" in argv:
		print("Running RedCedar in Smart mode")
		socketio.start_background_task(run_redcedar_smart)
	elif "smartcwd" in argv:
		print("Running RedCedar in Smart mode")
		socketio.start_background_task(run_redcedar_smart_cwd)
	elif "noredcedar" in argv:
		print("Running without RedCedar background process")
	else:
		print("Starting redcedar")
		socketio.start_background_task(run_redcedar_smart)
	
	socketio.run(app, host="0.0.0.0")
	
	if redcedar_obj != None:
		redcedar_obj.stop()
