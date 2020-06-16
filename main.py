from flask_socketio import SocketIO, emit
from flask import Flask, render_template, url_for, copy_current_request_context
from pathlib import Path
from queue import Queue
from random import random
from sys import argv
from threading import Thread, Event
from time import sleep
from redcedar import RedCedar
from redcedar.mock_cedar import MockCedar

app = Flask(__name__)
app.config['SECRET_KEY'] = 'my_secret'
app.config['DEBUG'] = False

#turn the flask app into a socketio app
socketio = SocketIO(app, async_mode=None, logger=True, engineio_logger=True)

#random number Generator Thread
thread = Thread()
thread_stop_event = Event()
no_numbers = False

redcedar_comm_queue = Queue()
redcedar_obj = None

def randomNumberGenerator():
	"""
	Generate a random number every 1 second and emit to a socketio instance (broadcast)
	"""
	while not thread_stop_event.isSet():
		number = round(random()*10, 3)
		print(number)
		socketio.emit('newnumber', {'number': number}, namespace='/websocket')
		socketio.sleep(5)

def run_redcedar():
	global redcedar_obj
	redcedar_obj = RedCedar(socketio, redcedar_comm_queue, Path("/usr/app/tosearch"))
	redcedar_obj.run()

def run_mockcedar():
	global redcedar_obj
	redcedar_obj = MockCedar(socketio, redcedar_comm_queue, Path("/usr/app/tosearch"))
	redcedar_obj.run()

@app.route('/')
def index():
	#only by sending this page first will the client be connected to the socketio instance
	return render_template('index.html')

@socketio.on('connect', namespace='/websocket')
def test_connect():
	redcedar_comm_queue.put("new connection")
	# need visibility of the global thread object
	global thread
	print('Client connected')

	#Start the random number generator thread only if the thread has not been started before.
	if not no_numbers and not thread.isAlive():
		print("Starting Thread")
		thread = socketio.start_background_task(randomNumberGenerator)

@socketio.on('disconnect', namespace='/websocket')
def test_disconnect():
	print('Client disconnected')

if __name__ == '__main__':
	if "mockcedar" in argv:
		print("Running with mock RedCedar background process")
		socketio.start_background_task(run_mockcedar)
	elif "noredcedar" in argv:
		print("Running without RedCedar background process")
	else:
		print("Starting redcedar")
		socketio.start_background_task(run_redcedar)
	if "nonumbers" in argv:
		no_numbers = True
	socketio.run(app, host="0.0.0.0")
	
	if redcedar_obj != None:
		redcedar_obj.stop()
