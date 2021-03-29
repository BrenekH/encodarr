import axios from "axios";
import React from "react";
import Button from "react-bootstrap/Button";
import Modal from "react-bootstrap/Modal";

import "./LibrariesTab.css";

import addLibraryIcon from "./addLibraryIcon.svg";

interface ILibrariesTabState {
	libraries: Array<Number>
	waitingOnServer: Boolean
	showCreateLibModal: Boolean
}

export class LibrariesTab extends React.Component<{}, ILibrariesTabState> {
	timerID: ReturnType<typeof setInterval>;

	constructor(props: any) {
		super(props);
		this.state = {
			libraries: [],
			waitingOnServer: true,
			showCreateLibModal: false,
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
		return (<>
			<img className="add-lib-ico" src={addLibraryIcon} alt="" height="20px" onClick={() => { this.setState({showCreateLibModal: true}); }} />
			<CreateLibraryModal show={this.state.showCreateLibModal} closeHandler={() => { this.setState({showCreateLibModal: false}); }} />
			<LibraryCard />
		</>);
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

interface CreateLibraryModalProps {
	show: Boolean,
	closeHandler: any,
}

function CreateLibraryModal(props: CreateLibraryModalProps) {
	return (<div>
		<Modal show={props.show} onHide={props.closeHandler}>
			<Modal.Header closeButton>
				<Modal.Title>Create New Library</Modal.Title>
			</Modal.Header>
			<Modal.Body>Hello, World!</Modal.Body>
			<Modal.Footer>
				<Button variant="secondary" onClick={props.closeHandler}>Close</Button>
			</Modal.Footer>
		</Modal>
	</div>);
}
