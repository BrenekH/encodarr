$(document).ready(function(){
	//connect to the socket server.
	var socket = io.connect("http://" + document.domain + ":" + location.port + "/websocket");
	var numbers_received = [];

	var completed_files = [];
	var skipped_files = [];

	var already_connected = false;

	//receive details from server
	socket.on("newnumber", function(msg) {
		console.log("Received number" + msg.number);
		//maintain a list of two numbers
		if (numbers_received.length >= 2){
			numbers_received.shift()
		}            
		numbers_received.push(msg.number);
		numbers_string = "";
		for (var i = 0; i < numbers_received.length; i++){
			numbers_string = numbers_string + "<p>" + numbers_received[i].toString() + "</p>";
		}
		$("#log").html(numbers_string);
	});

	socket.on("current_file_status_update", function(json_obj) {
		$("#total-time").html(json_obj.total_time);
		$("#percent-complete").html(json_obj.percentage + " Complete");
		$("#completion-estimate").html(json_obj.current_etr);
		$("#current-time-running").html(json_obj.current_time);
		$("#fps").html(json_obj.current_fps);
		$("#avg-fps").html(json_obj.avg_fps);
	});

	socket.on("connect_info", function(json_obj) {
		if (already_connected) { return; } 
		$("#file-counter").html(json_obj.file_count);
		$("#current-file").html(json_obj.current_file_path);
		// TODO: Parse, save, and display json_obj.completed_files and json_obj.skipped_files
		already_connected = true;
	});

	socket.on("current_file_update", function(json_obj) {
		console.log("current_file_update");
		console.log(json_obj);
		$("#current-file").html(json_obj.file_path);
		$("#file-counter").html(json_obj.file_count);
	});

	socket.on("file_complete", function(json_obj) {
		console.log("file_complete");
		console.log(json_obj);
	});

	socket.on("file_skip", function(json_obj) {
		console.log("file_skip");
		console.log(json_obj);
    });

});
