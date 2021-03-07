import axios from "axios";
import React from "react";
import Card from "react-bootstrap/Card";
import Table from "react-bootstrap/Table"

interface IJobHistory {
	datetime_completed: string,
	file: string,
}

interface IHistoryTabState {
	jobs: Array<IJobHistory>,
	waitingOnServer: Boolean,
}

export class HistoryTab extends React.Component<{}, IHistoryTabState> {
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
		axios.get("/api/web/v1/history").then((response) => {
			let history = response.data.history;
			if (history === undefined) {
				console.error("Response from /api/web/v1/history returned undefined for data.history");
				return;
			}
			history.reverse();
			this.setState({
				jobs: history,
				waitingOnServer: false,
			});
		}).catch((error) => {
			console.error(`Request to /api/web/v1/history failed with error: ${error}`);
		});
	}

	render(): React.ReactNode {
		const jobEntries = this.state.jobs.map((v: IJobHistory, i: number) => {
			return <TableEntry key={i} datetime={v.datetime_completed} file={v.file}/>;
		});

		const waitingForServerEntry = <tr><th scope="row">-</th><td>Waiting on server</td></tr>;

		return (<Card>
			<Table hover size="sm">
				<thead>
					<tr>
						<th scope="col">Time Completed</th>
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
	datetime: string,
	file: string,
}

function TableEntry(props: ITableEntryProps) {
	return <tr><td>{props.datetime}</td><td>{props.file}</td></tr>;
}
