import axios from "axios";
import React from "react";
import Button from "react-bootstrap/Button";
import FormControl from "react-bootstrap/FormControl";
import InputGroup from "react-bootstrap/InputGroup";
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

interface ICreateLibraryModalProps {
	show: Boolean,
	closeHandler: any,
}

interface ICreateLibraryModalState {
	folder: string,
}

class CreateLibraryModal extends React.Component<ICreateLibraryModalProps, ICreateLibraryModalState> {
	constructor(props: ICreateLibraryModalProps) {
		super(props);

		this.state = {
			folder: "",
		}
	}

	render(): React.ReactNode {
		return (<div>
			<Modal show={this.props.show} onHide={this.props.closeHandler}>
				<Modal.Header closeButton>
					<Modal.Title>Create New Library</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Folder</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder="/home/user/lib1"
							aria-label="folder"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ folder: event.target.value }); }}
							value={this.state.folder}
						/>
					</InputGroup>
					<p>priority</p>
					<p>fs_check_interval</p>
					<p>Plugin Pipeline</p>
					<p>path_masks</p>
				</Modal.Body>
				<Modal.Footer>
					<Button variant="secondary" onClick={this.props.closeHandler}>Close</Button>
				</Modal.Footer>
			</Modal>
		</div>);
	}
}
