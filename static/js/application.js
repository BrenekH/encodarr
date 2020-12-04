$(document).ready(function() {
	// Connect to the socket server.
	var socket = io.connect("http://" + document.domain + ":" + location.port + "/updates");

	socket.on("current_job_status_update", function(json_obj) {
		$("#current-stage").html(`Stage: ${json_obj.stage}`);

		setProgressBar(json_obj.percentage);
		$("#job-elapsed-time").html(json_obj.job_elapsed_time);
		$("#stage-estimated-time-remaining").html(json_obj.stage_estimated_time_remaining);
		$("#fps").html(json_obj.fps);
		$("#stage-elapsed-time").html(json_obj.stage_elapsed_time);
	});

	socket.on("current_job_update", function(json_obj) {
		$("#current-file").html(json_obj.file);
	});
});

function setProgressBar(progress) {
	let progressBar = document.getElementById("progress-bar");
	progressBar.textContent = `${progress}%`;
	progressBar.setAttribute("aria-valuenow", progress);
	progressBar.style.width = `${progress}%`;
}
