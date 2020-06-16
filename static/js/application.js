$(document).ready(function(){
	//connect to the socket server.
	var socket = io.connect("http://" + document.domain + ":" + location.port + "/websocket");

	var completed_files = [];
	var skipped_files = [];

	var already_connected = false;

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
		already_connected = true;
		$("#file-counter").html(json_obj.file_count);
		$("#current-file").html(json_obj.current_file_path);
		json_obj.completed_files.forEach(function(file_obj) {
			completed_files.push(file_obj);
		});
		json_obj.skipped_files.forEach(function(file_obj) {
			skipped_files.push(file_obj);
		});
		renderCompletedFiles();
		renderSkippedFiles();
	});

	socket.on("current_file_update", function(json_obj) {
		$("#current-file").html(json_obj.file_path);
		$("#file-counter").html(json_obj.file_count);
	});

	socket.on("file_complete", function(json_obj) {
		completed_files.unshift(json_obj);
		renderCompletedFiles();
	});

	socket.on("file_skip", function(json_obj) {
		skipped_files.unshift(json_obj);
		renderSkippedFiles();
	});
	
	function renderCompletedFiles() {
		var html_string = "";
		completed_files.forEach(function(file_obj) {
			html_string += `
			<div class="row text-center">
			<div class="col completed-file-path">
				<p class="text-break">${file_obj.file_path}</p>
			</div>
			<div class="col completed-avg-fps">${file_obj.avg_fps}</div>
			<div class="col completed-time">${file_obj.time_taken}</div>
			<div class="col completed-timestamp">${file_obj.timestamp}</div>
			</div>
			`;
		});

		$("#completed-log").html(html_string);
	}

	function renderSkippedFiles() {
		var html_string = "";
		skipped_files.forEach(function(file_obj) {
			html_string += `
			<div class="row">
				<div class="col skipped-file-path">
					<p class="text-break">${file_obj.file_path}</p>
				</div>
				<div class="col skipped-reason">
					<p class="text-break">${file_obj.reason}</p>
				</div>
				<div class="col completed-timestamp">
					<p>${file_obj.timestamp}</p>
				</div>
			</div>
			`;
		});
		$("#skipped-log").html(html_string);
	}

});
