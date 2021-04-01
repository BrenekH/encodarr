import axios from "axios";
import React from "react";
import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import FormControl from "react-bootstrap/FormControl";
import InputGroup from "react-bootstrap/InputGroup";
import Modal from "react-bootstrap/Modal";
import Table from "react-bootstrap/Table";

import { AudioImage } from "./shared/AudioImage";
import { VideoImage } from "./shared/VideoImage";

import "./LibrariesTab.css";
import "../spacers.css";

import addLibraryIcon from "./addLibraryIcon.svg";

interface ILibrariesTabState {
	libraries: Array<number>
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
		const libsList = this.state.libraries.map((v) => {
			return (<div key={v}><LibraryCard id={v} /><div className="smol-spacer"></div></div>);
		});

		return (<>
			<img className="add-lib-ico" src={addLibraryIcon} alt="" height="20px" onClick={() => { this.setState({showCreateLibModal: true}); }} />
			<CreateLibraryModal show={this.state.showCreateLibModal} closeHandler={() => { this.setState({showCreateLibModal: false}); }} />
			<div className="smol-spacer"></div>
			{libsList}
		</>);
	}
}

interface ILibraryCardProps {
	id: number,
}

interface ILibraryCardState {
	folder: string,
	priority: string,
	fs_check_interval: string,
	path_masks: string,
	queue: Array<IQueuedJob>,

	showEditModal: Boolean,
	showQueueModal: Boolean,
}

class LibraryCard extends React.Component<ILibraryCardProps, ILibraryCardState> {
	constructor(props: ILibraryCardProps) {
		super(props);

		this.state = {
			folder: "",
			priority: "",
			fs_check_interval: "",
			path_masks: "",
			queue: [],

			showEditModal: false,
			showQueueModal: false,
		};
	}

	componentDidMount() {
		this.getLibraryData();
	}

	getLibraryData() {
		axios.get(`/api/web/v1/library/${this.props.id}`).then((response) => {
			this.setState({
				folder: response.data.folder,
				priority: response.data.priority,
				fs_check_interval: response.data.fs_check_interval,
				path_masks: response.data.path_masks.join(","),
				queue: response.data.queue.Items,
			});
		}).catch((error) => {
			console.error(`Request to /api/web/v1/library/${this.props.id} failed with error: ${error}`)
		});
	}

	render() {
		return (
		<>
			<Card>
				<Card.Header className="text-center"><h5>{this.state.folder}</h5></Card.Header>
				<p className="text-center">Priority: {this.state.priority}</p>
				<p className="text-center">File System Check Interval: {this.state.fs_check_interval}</p>
				<p className="text-center">Path Masks: {this.state.path_masks}</p>
				<Button variant="secondary" onClick={() => {this.setState({showQueueModal: true})}}>Queue</Button>
				<Button variant="primary" onClick={() => {this.setState({showEditModal: true})}}>Edit</Button>
			</Card>
			{(this.state.showEditModal) ? (<EditLibraryModal
				show={true}
				closeHandler={() => { this.setState({showEditModal: false}); }}
				id={this.props.id}
				folder={this.state.folder}
				priority={this.state.priority}
				fs_check_interval={this.state.fs_check_interval}
				path_masks={this.state.path_masks}
			/>) : null}

			{(this.state.showQueueModal) ? (<QueueModal
				show={true}
				closeHandler={() => { this.setState({showQueueModal: false}); }}
				queue={this.state.queue}
			/>) : null}
		</>);
	}
}

interface ICreateLibraryModalProps {
	show: Boolean,
	closeHandler: any,
}

interface ICreateLibraryModalState {
	folder: string,
	priority: string,
	fs_check_interval: string,
	path_masks: string,
}

class CreateLibraryModal extends React.Component<ICreateLibraryModalProps, ICreateLibraryModalState> {
	constructor(props: ICreateLibraryModalProps) {
		super(props);

		this.state = {
			folder: "",
			priority: "",
			fs_check_interval: "",
			path_masks: "",
		}

		this.submitLib = this.submitLib.bind(this);
	}

	submitLib(): void {
		let data = {
			folder: this.state.folder,
			priority: parseInt(this.state.priority),
			fs_check_interval: this.state.fs_check_interval,
			path_masks: this.state.path_masks.split(","),
		};
		axios.post("/api/web/v1/library/new", data).then(() => {
			this.props.closeHandler();
		}).catch((error) => {
			console.error(`/api/web/v1/library/new failed with error: ${error}`)
		});
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

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Priority</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder="0"
							aria-label="priority"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ priority: event.target.value }); }}
							value={this.state.priority}
						/>
					</InputGroup>

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>File System Check Interval</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder="0h0m0s"
							aria-label="fs_check_interval"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ fs_check_interval: event.target.value }); }}
							value={this.state.fs_check_interval}
						/>
					</InputGroup>
					{/* <p>Plugin Pipeline</p> */}

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Path Masks</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder="Plex Versions,private,.m4a"
							aria-label="path_masks"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ path_masks: event.target.value }); }}
							value={this.state.path_masks}
						/>
					</InputGroup>
				</Modal.Body>
				<Modal.Footer>
					<Button variant="secondary" onClick={this.props.closeHandler}>Close</Button>
					<Button variant="primary" onClick={this.submitLib}>Create</Button>
				</Modal.Footer>
			</Modal>
		</div>);
	}
}

interface IEditLibraryModalProps {
	show: Boolean,
	closeHandler: any,
	id: number,
	folder: string,
	priority: string,
	fs_check_interval: string,
	path_masks: string,
}

interface IEditLibraryModalState {
	folder: string,
	priority: string,
	fs_check_interval: string,
	path_masks: string,
}

class EditLibraryModal extends React.Component<IEditLibraryModalProps, IEditLibraryModalState> {
	constructor(props: IEditLibraryModalProps) {
		super(props);

		this.state = {
			folder: props.folder,
			priority: props.priority,
			fs_check_interval: props.fs_check_interval,
			path_masks: props.path_masks,
		}

		this.putChanges = this.putChanges.bind(this);
	}

	putChanges(): void {
		let data = {
			folder: this.state.folder,
			priority: parseInt(this.state.priority),
			fs_check_interval: this.state.fs_check_interval,
			path_masks: this.state.path_masks.split(","),
		};
		axios.put(`/api/web/v1/library/${this.props.id}`, data).then(() => {
			this.props.closeHandler();
		}).catch((error) => {
			console.error(`/api/web/v1/library/${this.props.id} failed with error: ${error}`)
		});
	}

	render(): React.ReactNode {
		return (<div>
			<Modal show={this.props.show} onHide={this.props.closeHandler}>
				<Modal.Header closeButton>
					<Modal.Title>Edit Library</Modal.Title>
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

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Priority</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder="0"
							aria-label="priority"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ priority: event.target.value }); }}
							value={this.state.priority}
						/>
					</InputGroup>

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>File System Check Interval</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder="0h0m0s"
							aria-label="fs_check_interval"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ fs_check_interval: event.target.value }); }}
							value={this.state.fs_check_interval}
						/>
					</InputGroup>
					{/* <p>Plugin Pipeline</p> */}

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Path Masks</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder="Plex Versions,private,.m4a"
							aria-label="path_masks"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ path_masks: event.target.value }); }}
							value={this.state.path_masks}
						/>
					</InputGroup>
				</Modal.Body>
				<Modal.Footer>
					<Button variant="secondary" onClick={this.props.closeHandler}>Close</Button>
					<Button variant="primary" onClick={this.putChanges}>Update</Button>
				</Modal.Footer>
			</Modal>
		</div>);
	}
}

interface IQueueModalProps {
	show: Boolean,
	closeHandler: any,
	queue: Array<IQueuedJob>,
}

class QueueModal extends React.Component<IQueueModalProps> {
	render(): React.ReactNode {
		const qList = this.props.queue.map((v: IQueuedJob, i: number) => {
			return <TableEntry key={v.uuid} index={i+1} path={v.path} videoOperation={v.parameters.hevc} audioOperation={v.parameters.stereo}/>;
		});

		return (<div>
			<Modal show={this.props.show} onHide={this.props.closeHandler} size="lg">
				<Modal.Header closeButton>
					<Modal.Title>Queue</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Table>
						<thead>
							<tr>
								<th scope="col">#</th>
								<th scope="col">File</th>
							</tr>
						</thead>
						<tbody>
							{qList}
						</tbody>
					</Table>
				</Modal.Body>
				<Modal.Footer>
					<Button variant="secondary" onClick={this.props.closeHandler}>Close</Button>
				</Modal.Footer>
			</Modal>
		</div>);
	}
}

interface IQueuedJob {
	uuid: string,
	path: string,
	parameters: {
		hevc: Boolean,
		stereo: Boolean,
	},
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
