$(document).ready(function() {
	// Connect to the socket server.
	var socket = io.connect(`http://${document.domain}:${location.port}/updates`);

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

	socket.on("current_jobs_update", function(json_obj) {
		let HTMLString = "";
		let looped = false; // This var is used to tell if the no running jobs message should be displayed or not

		for (const uuid in json_obj) {
			looped = true;
			HTMLString += renderRunningJobCard(uuid,
				json_obj[uuid].file,
				json_obj[uuid].runner_name,
				json_obj[uuid].status.stage,
				json_obj[uuid].status.percentage,
				json_obj[uuid].status.job_elapsed_time,
				json_obj[uuid].status.fps,
				json_obj[uuid].status.stage_elapsed_time,
				json_obj[uuid].status.stage_estimated_time_remaining);
		}

		if (!looped) {
			HTMLString = `<h5 style="text-align: center;">No running jobs</h5>`
		}

		$("#running-jobs").html(HTMLString)
	});
});

$('a[data-toggle="tab"]').on("shown.bs.tab", function (e) {
	var target = $(e.target).attr("href");
	if (target == "#queue") {
		updateQueue();
	} else if (target == "#history") {
		updateHistory();
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
			finalHTMLString += renderQueueEntry(i, queue[i-1].filename, queue[i-1].video_op, queue[i-1].audio_op);
		}
		$("#queue-content").html(finalHTMLString);
		enableTooltips();
	}).catch(function (error) {
		console.log(`Request to /api/v1/queue failed with error: ${error}`);
	});
}

function renderQueueEntry(entryNumber, filePath, videoOp, audioOp) {
	let videoHTML = "";
	if (videoOp) {
		videoHTML = `<img class="playButtonImage queue-icon" src="/static/svg/play_button.svg" alt="Play Button" height="20px" data-bs-toggle="tooltip" data-bs-placement="top" title="File will be encoded to HEVC">`
	}

	let audioHTML = "";
	if (audioOp) {
		audioHTML = `<img class="queue-icon" src="/static/svg/headphones.svg" alt="Headphones" height="20px" data-bs-toggle="tooltip" data-bs-placement="top" title="An additional stereo audio track will be created">`
	}
	return `<tr><th scope="row">${entryNumber}</th><td>${filePath}</td><td><div class="queue-icon-container">${videoHTML}${audioHTML}</div></td></tr>\n`;
}

function updateHistory() {
	axios.get("/api/v1/history").then(function (response) {
		let history = response.data.history;
		if (history === undefined) {
			console.error("Response from /api/v1/history returned undefined for data.history")
			return
		}
		let finalHTMLString = "";
		for (let i = 0; i < history.length; i++) {
			let obj = history[i];
			finalHTMLString += renderHistoryEntry(obj.datetime_completed, obj.file);
		}
		$("#history-content").html(finalHTMLString);
	}).catch(function (error) {
		console.log(`Request to /api/v1/history failed with error: ${error}`);
	});
}

function renderHistoryEntry(dateTimeString, filePath) {
	return `<tr><td>${dateTimeString}</td><td>${filePath}</td></tr>`;
}

function enableTooltips() {
	var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
	var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
		return new bootstrap.Tooltip(tooltipTriggerEl)
	})
}

function renderRunningJobCard(uuid, filename, runnerName, stageValue, progress, jobElapsedTime, fps, stageElapsedTime, stageEstimatedTimeRemaining) {
	return `
<div class="card" id="${uuid}-job-card" style="padding: 1rem;">
	<div class="card-header text-center" style="padding-bottom: .25rem;">
		<h5 id="${uuid}-current-file">${filename}</h5>
		<h6 id="${uuid}-current-stage">Stage: ${stageValue}</h6>
		<h6 id="${uuid}-runner-name">Runner: ${runnerName}</h6>
	</div>
	<div class="progress" style="margin-bottom: 1rem; margin-top: 1rem; height: 2rem;">
		<div class="progress-bar progress-bar-striped progress-bar-animated" id="${uuid}-progress-bar" role="progressbar" style="width: ${progress}%; font-size: 0.9rem;" aria-valuenow="${progress}" aria-valuemin="0" aria-valuemax="100">${progress}%</div>
	</div>
	<div class="row">
		<div class="col">
			<h6 class="job-elapsed-time-label text-right">Job Elapsed Time:</h6>
		</div>
		<div class="col job-elapsed-time">
			<p id="${uuid}-job-elapsed-time">${jobElapsedTime}</p>
		</div>
		<div class="col">
			<h6 class="fps-label text-right">FPS:</h6>
		</div>
		<div class="col fps">
			<p id="${uuid}-fps">${fps}</p>
		</div>
	</div>
	<div class="row">
		<div class="col">
			<h6 class="stage-elapsed-time-label text-right">Stage Elapsed Time:</h6>
		</div>
		<div class="col stage-elapsed-time">
			<p id="${uuid}-stage-elapsed-time">${stageElapsedTime}</p>
		</div>
		<div class="col">
			<h6 class="stage-estimated-time-remaining-label text-right">Stage Estimated Time Remaining:</h6>
		</div>
		<div class="col stage-estimated-time-remaining">
			<p id="${uuid}-stage-estimated-time-remaining">${stageEstimatedTimeRemaining}</p>
		</div>
	</div>
</div>
<div class="smol-spacer"></div>
`;
}
