var currentTab = "running";

$(document).ready(function () {
	updateCurrentTab();

	window.setInterval(function(){
		updateCurrentTab();
	}, 2000);
});

$('a[data-toggle="tab"]').on("shown.bs.tab", function (e) {
	var target = $(e.target).attr("href");
	if (target == "#queue") {
		updateQueue();
		currentTab = "queue";
	} else if (target == "#history") {
		updateHistory();
		currentTab = "history";
	} else if (target == "#running-jobs") {
		updateRunning();
		currentTab = "running";
	} else if (target == "#settings") {
		updateSettings();
		currentTab = "settings";
	}
});

function updateCurrentTab() {
	switch (currentTab) {
		case "running":
			updateRunning();
			break;
		case "queue":
			updateQueue();
			break;
		case "history":
			updateHistory();
			break;
		default:
			break;
	}
}

function setProgressBar(progress) {
	let progressBar = document.getElementById("progress-bar");
	progressBar.textContent = `${progress}%`;
	progressBar.setAttribute("aria-valuenow", progress);
	progressBar.style.width = `${progress}%`;
}

function updateRunning() {
	axios.get("/api/web/v1/running").then(function (response) {
		let jobs = response.data.jobs;
		if (jobs === undefined) {
			console.error("Response from /api/web/v1/running returned undefined for data.jobs");
		}
		jobs.sort((a, b) => {
			if (parseFloat(a.status.percentage) > parseFloat(b.status.percentage)) {
				return -1;
			}
			return 1;
		})

		let HTMLString = "";
		let looped = false;

		for (let i = 0; i < jobs.length; i++) {
			looped = true;
			let cJob = jobs[i];
			HTMLString += renderRunningJobCard(
				cJob.job.uuid,
				cJob.job.path,
				cJob.runner_name,
				cJob.status.stage,
				cJob.status.percentage,
				cJob.status.job_elapsed_time,
				cJob.status.fps,
				cJob.status.stage_elapsed_time,
				cJob.status.stage_estimated_time_remaining,
				cJob.job.parameters.hevc,
				cJob.job.parameters.stereo);
		}

		if (!looped) {
			HTMLString = `<h5 style="text-align: center;">No running jobs</h5>`;
		}

		disableTooltips();
		$("#running-jobs").html(HTMLString)
		enableTooltips();
	}).catch(function (error) {
		console.error(`Request to /api/web/v1/running failed with error: ${error}`);
	});
}

function updateQueue() {
	axios.get("/api/web/v1/queue").then(function (response) {
		let queue = response.data.queue; // List of files in queue order
		if (queue === undefined) {
			console.error("Response from /api/web/v1/queue returned undefined for data.queue");
			return;
		}
		let finalHTMLString = "";
		for (let i = 1; i <= queue.length; i++) {
			finalHTMLString += renderQueueEntry(i, queue[i-1].path, queue[i-1].parameters.hevc, queue[i-1].parameters.stereo);
		}
		finalHTMLString += `\n<div class="smol-spacer"></div>`;
		disableTooltips();
		$("#queue-content").html(finalHTMLString);
		enableTooltips();
	}).catch(function (error) {
		console.error(`Request to /api/web/v1/queue failed with error: ${error}`);
	});
}

function renderQueueEntry(entryNumber, filePath, videoOp, audioOp) {
	let videoHTML = "";
	if (videoOp) {
		videoHTML = `<img class="playButtonImage queue-icon" src="/resources/svg/play_button.svg" alt="Play Button" height="20px" data-bs-toggle="tooltip" data-bs-placement="top" title="File will be encoded to HEVC">`;
	}

	let audioHTML = "";
	if (audioOp) {
		audioHTML = `<img class="queue-icon" src="/resources/svg/headphones.svg" alt="Headphones" height="20px" data-bs-toggle="tooltip" data-bs-placement="top" title="An additional stereo audio track will be created">`;
	}
	return `<tr><th scope="row">${entryNumber}</th><td>${filePath}</td><td><div class="queue-icon-container">${videoHTML}${audioHTML}</div></td></tr>\n`;
}

function updateHistory() {
	axios.get("/api/web/v1/history").then(function (response) {
		let history = response.data.history;
		if (history === undefined) {
			console.error("Response from /api/web/v1/history returned undefined for data.history");
			return;
		}
		let finalHTMLString = "";
		for (let i = 1; i <= history.length; i++) {
			let obj = history[history.length-i];
			finalHTMLString += renderHistoryEntry(obj.datetime_completed, obj.file);
		}
		finalHTMLString += `\n<div class="smol-spacer"></div>`;
		$("#history-content").html(finalHTMLString);
	}).catch(function (error) {
		console.error(`Request to /api/web/v1/history failed with error: ${error}`);
	});
}

function renderHistoryEntry(dateTimeString, filePath) {
	return `<tr><td>${dateTimeString}</td><td>${filePath}</td></tr>`;
}

function enableTooltips() {
	$('[data-bs-toggle="tooltip"]').tooltip();
}

function disableTooltips() {
	$('[data-bs-toggle="tooltip"]').tooltip("dispose");
}

function renderRunningJobCard(uuid, filename, runnerName, stageValue, progress, jobElapsedTime, fps, stageElapsedTime, stageEstimatedTimeRemaining, videoOp, audioOp) {
	let videoHTML = "";
	if (videoOp) {
		videoHTML = `<img class="queue-icon" src="/resources/svg/play_button.svg" alt="Play Button" height="20px" data-bs-toggle="tooltip" data-bs-placement="top" title="File will be encoded to HEVC">`
	}

	let audioHTML = "";
	if (audioOp) {
		audioHTML = `<img class="running-stereo-icon queue-icon" src="/resources/svg/headphones.svg" alt="Headphones" height="20px" data-bs-toggle="tooltip" data-bs-placement="top" title="An additional stereo audio track will be created">`
	}

	return `
<div class="card" id="${uuid}-job-card" style="padding: 1rem;">
	<div class="card-header text-center" style="padding-bottom: .25rem;">
		<div class="file-image-container">
			<h5 id="${uuid}-current-file">${filename}</h5>
			<div class="svg-flex-container">
				${videoHTML}
				${audioHTML}
			</div>
		</div>
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

// Settings tab functions
function updateSettings() {
	lockSettings();

	axios.get("/api/web/v1/settings").then(function(response) {
		document.getElementById("fs-check-interval").value = response.data.FileSystemCheckInterval;
		document.getElementById("health-check-interval").value = response.data.HealthCheckInterval;
		document.getElementById("unresponsive-runner-timeout").value = response.data.HealthCheckTimeout;
		document.getElementById("log-verbosity-select").value = response.data.LogVerbosity;

		unlockSettings();
	});
}

document.getElementById("save-settings-btn").onclick = function() {
	axios.put("/api/web/v1/settings", {
		"FileSystemCheckInterval": document.getElementById("fs-check-interval").value,
		"HealthCheckInterval": document.getElementById("health-check-interval").value,
		"HealthCheckTimeout": document.getElementById("unresponsive-runner-timeout").value,
		"LogVerbosity": document.getElementById("log-verbosity-select").value
	}).then(function(response) {
		if (response.status >= 200 && response.status <= 299) {
			document.getElementById("saved-container").innerHTML = `<p class="pop-in-out" style="display:inline;">Saved!</p>`;
		} else {
			console.error(response);
		}

		updateSettings(); // Update settings is used to correct malicious users by resetting invalid values in the UI
	});
};

function lockSettings() {
	document.getElementById("fs-check-interval").disabled = true;
	document.getElementById("health-check-interval").disabled = true;
	document.getElementById("unresponsive-runner-timeout").disabled = true;

	document.getElementById("log-verbosity-select").hidden = true;
}

function unlockSettings() {
	document.getElementById("fs-check-interval").disabled = false;
	document.getElementById("health-check-interval").disabled = false;
	document.getElementById("unresponsive-runner-timeout").disabled = false;

	document.getElementById("log-verbosity-select").hidden = false;
}
