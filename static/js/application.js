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

$('a[data-toggle="tab"]').on("shown.bs.tab", function (e) {
	var target = $(e.target).attr("href");
	if (target == "#queue") {
		updateQueue();
	}
});

function setProgressBar(progress) {
	let progressBar = document.getElementById("progress-bar");
	progressBar.textContent = `${progress}%`;
	progressBar.setAttribute("aria-valuenow", progress);
	progressBar.style.width = `${progress}%`;
}

function updateQueue() {
	axios.get("/api/v1/queue").then(function (response) {
		let queue = response.data.queue; // List of files in queue order
		if (queue === undefined) {
			console.error("Response from /api/v1/queue returned undefined for data.queue")
			return
		}
		let finalHTMLString = "";
		for (let i = 1; i <= queue.length; i++) {
			finalHTMLString += renderQueueEntry(i, queue[i-1]);
		}
		$("#queue-content").html(finalHTMLString);
	}).catch(function (error) {
		console.log(`Request to /api/v1/queue failed with error: ${error}`);
	});
}

function renderQueueEntry(entryNumber, filePath) {
	return `<tr><th scope="row">${entryNumber}</th><td>${filePath}</td></tr>\n`;
}
