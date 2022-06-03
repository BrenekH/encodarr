import axios from "axios";
import React from "react";
import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import Col from "react-bootstrap/Col";
import Modal from "react-bootstrap/Modal";
import ProgressBar from "react-bootstrap/ProgressBar";
import Row from "react-bootstrap/Row";

import "./RunningTab.css";
import "../spacers.css";

import InfoIIcon from "./InfoIIcon";
import TerminalIcon from "./shared/TerminalIcon";

interface IRunningJob {
	runner_name: string,
	job: {
		uuid: string,
		path: string,
		command: Array<string>,
	},
	status: {
		fps: string,
		percentage: string,
		job_elapsed_time: string,
		stage: string,
		stage_elapsed_time: string,
		stage_estimated_time_remaining: string,
	},
}

interface IRunningTabState {
	jobs: Array<IRunningJob>,
	waitingOnServer: boolean,
	showModal: boolean,
	waitingRunnersText: String,
}

export class RunningTab extends React.Component<{}, IRunningTabState> {
	timerID: ReturnType<typeof setInterval>;

	constructor(props: any) {
		super(props);
		this.state = {
			jobs: [],
			waitingOnServer: true,
			showModal: false,
			waitingRunnersText: "",
		};

		// This is just so Typescript doesn't whine about timerID not being instantiated.
		this.timerID = setTimeout(() => { }, Number.POSITIVE_INFINITY);
		clearInterval(this.timerID);
	}

	componentDidMount() {
		this.tick();
		this.timerID = setInterval(
			() => this.tick(),
			2000 // Two seconds
		);
	}

	componentWillUnmount() {
		clearInterval(this.timerID);
	}

	tick() {
		// Update currently running jobs
		axios.get("/api/web/v1/running").then((response) => {
			let rJobs: Array<IRunningJob> = response.data.jobs;
			if (rJobs === undefined) {
				console.error("Response from /api/web/v1/running returned undefined for data.jobs");
				return;
			}

			rJobs.sort((a, b) => {
				if (parseFloat(a.status.percentage) > parseFloat(b.status.percentage)) {
					return -1;
				}
				return 1;
			});

			this.setState({
				jobs: rJobs,
				waitingOnServer: false,
			});
		}).catch((error) => {
			console.error(`Request to /api/web/v1/running failed with error: ${error}`);
		});

		// Update waiting runners
		axios.get("/api/web/v1/waitingrunners").then((response) => {
			if (response.data.Runners.length === 0) {
				this.setState({
					waitingRunnersText: "No waiting runners",
				});
			} else {
				let runStr = response.data.Runners.toString();
				this.setState({
					waitingRunnersText: runStr,
				});
			}
		}).catch((error) => {
			console.error(`Request to /api/web/v1/waitingrunners failed with error: ${error}`);
		});
	}

	render(): React.ReactNode {
		const handleClose = () => this.setState({ showModal: false });
		const handleShow = () => this.setState({ showModal: true });

		const jobsList = this.state.jobs.map((v) => {
			return (<RunningCard
				key={v.job.uuid}
				fps={v.status.fps}
				uuid={v.job.uuid}
				filename={v.job.path}
				progress={v.status.percentage}
				runnerName={v.runner_name}
				stageValue={v.status.stage}
				jobElapsedTime={v.status.job_elapsed_time}
				stageElapsedTime={v.status.stage_elapsed_time}
				stageEstimatedTimeRemaining={v.status.stage_estimated_time_remaining}
				command={v.job.command.join(" ")}
			/>);
		});

		return (<div>
			<InfoIIcon className="info-i" height="20px" width="20px" onClick={handleShow} />
			{(jobsList.length !== 0) ? jobsList : <h5 className="text-center">No running jobs</h5>}

			<Modal show={this.state.showModal} onHide={handleClose}>
				<Modal.Header closeButton>
					<Modal.Title>Waiting Runners</Modal.Title>
				</Modal.Header>
				<Modal.Body>{this.state.waitingRunnersText}</Modal.Body>
				<Modal.Footer>
					<Button variant="secondary" onClick={handleClose}>Close</Button>
				</Modal.Footer>
			</Modal>
		</div>);
	}
}

interface IRunningCardProps {
	fps: string,
	uuid: string,
	filename: string,
	progress: string,
	runnerName: string,
	stageValue: string,
	jobElapsedTime: string,
	stageElapsedTime: string,
	stageEstimatedTimeRemaining: string,
	command: string,
}

function RunningCard(props: IRunningCardProps) {
	return (<div>
		<Card style={{ padding: '1rem' }}>
			<Card.Header className="text-center">
				<div className="file-image-container">
					<h5>{props.filename}</h5>
					<TerminalIcon height="20px" width="25px" title={props.command} />
				</div>

				<h6>Stage: {props.stageValue}</h6>
				<h6>Runner: {props.runnerName}</h6>
			</Card.Header>

			<ProgressBar className="progress-bar-style" animated now={parseFloat(props.progress)} label={`${props.progress}%`} />

			<Row>
				<Col><h6 className="text-right">Job Elapsed Time:</h6></Col>
				<Col><p>{props.jobElapsedTime}</p></Col>
				<Col><h6 className="text-right">FPS:</h6></Col>
				<Col><p>{props.fps}</p></Col>
			</Row>
			<Row>
				<Col><h6 className="text-right">Stage Elapsed Time:</h6></Col>
				<Col><p>{props.stageElapsedTime}</p></Col>
				<Col><h6 className="text-right">Stage Estimated Time Remaining:</h6></Col>
				<Col><p>{props.stageEstimatedTimeRemaining}</p></Col>
			</Row>
		</Card>
		<div className="smol-spacer"></div>
	</div>);
}
