import axios from "axios";
import React from "react";
import Card from "react-bootstrap/Card";
import Table from "react-bootstrap/Table"

import { AudioImage } from "./shared/AudioImage";
import { VideoImage } from "./shared/VideoImage";

import "./QueueTab.css";

interface IQueuedJob {
	uuid: string,
	path: string,
	parameters: {
		hevc: Boolean,
		stereo: Boolean,
	}
}

interface IQueueTabState {
	jobs: Array<IQueuedJob>,
	waitingOnServer: Boolean,
}

export class QueueTab extends React.Component<{}, IQueueTabState> {
	timerID: ReturnType<typeof setInterval>;

	constructor(props: any) {
		super(props);
		this.state = {
			jobs: [],
			waitingOnServer: true,
		};

		// This is just so Typescript doesn't whine about timerID not being instantiated.
		this.timerID = setTimeout(() => {}, Number.POSITIVE_INFINITY);
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
		axios.get("/api/web/v1/queue").then((response) => {
			let queue = response.data.queue;
			if (queue === undefined) {
				console.error("Response from /api/web/v1/queue returned undefined for data.queue");
				return;
			}
			this.setState({
				jobs: queue,
				waitingOnServer: false,
			});
		}).catch((error) => {
			console.error(`Request to /api/web/v1/queue failed with error: ${error}`);
		});
	}

	render(): React.ReactNode {
		const jobEntries = this.state.jobs.map((v: IQueuedJob, i: number) => {
			return <TableEntry key={v.uuid} index={i+1} path={v.path} videoOperation={v.parameters.hevc} audioOperation={v.parameters.stereo}/>;
		});

		const waitingForServerEntry = <tr><th scope="row">-</th><td>Waiting on server</td></tr>;

		return (<Card>
			<Table hover size="sm">
				<thead>
					<tr>
						<th scope="col">#</th>
						<th scope="col">File</th>
					</tr>
				</thead>
				<tbody>
					{(!this.state.waitingOnServer) ? jobEntries : waitingForServerEntry}
				</tbody>
			</Table>
		</Card>);
	}
}

interface ITableEntryProps {
	index: number,
	path: string,
	videoOperation: Boolean,
	audioOperation: Boolean,
}

function TableEntry(props: ITableEntryProps) {
	return (<tr>
		<th scope="row">{props.index}</th>
		<td>{props.path}</td>
		<td>
			<div className="queue-icon-container">
				{(props.videoOperation) ? <span className="play-button-image"><VideoImage/></span> : null}
				{(props.audioOperation) ? <AudioImage /> : null}
			</div>
		</td>
	</tr>);
}
