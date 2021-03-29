import axios from "axios";
import React from "react";

import "./LibrariesTab.css";

interface ILibrariesTabState {
	libraries: Array<Number>
	waitingOnServer: Boolean
}

export class LibrariesTab extends React.Component<{}, ILibrariesTabState> {
	timerID: ReturnType<typeof setInterval>;

	constructor(props: any) {
		super(props);
		this.state = {
			libraries: [],
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
		axios.get("/api/web/v1/libraries").then((response) => {
			if (response.status === 200) {
				this.setState({
					libraries: response.data.IDs
				});
			}
		}).catch((error) => {
			console.error(`Request to /api/web/v1/libraries failed with error: ${error}`);
		});
	}

	render(): React.ReactNode {
		return <LibraryCard />;
	}
}

function LibraryCard() {
	let id = 1;
	axios.get(`/api/web/v1/library/${id}`).then((response) => {
		// TODO: Add to state (to load)
	}).catch((error) => {
		console.error(`Request to /api/web/v1/library/${id} failed with error: ${error}`)
	});

	return (<div></div>);
}
