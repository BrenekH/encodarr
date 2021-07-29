import axios from "axios";
import React from "react";
import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import FormControl from "react-bootstrap/FormControl";
import InputGroup from "react-bootstrap/InputGroup";
import Modal from "react-bootstrap/Modal";
import Table from "react-bootstrap/Table";

import "./LibrariesTab.css";
import "../spacers.css";

import addLibraryIcon from "./addLibraryIcon.svg";
import TerminalIcon from "./shared/TerminalIcon";

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
	target_video_codec: string,
	create_stereo_audio: boolean,
	skip_hdr: boolean,
	use_hardware: boolean,
	hardware_codec: string,
	hw_device: string,

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
			target_video_codec: "HEVC",
			create_stereo_audio: true,
			skip_hdr: true,
			use_hardware: false,
			hardware_codec: "",
			hw_device: "",

			showEditModal: false,
			showQueueModal: false,
		};
	}

	componentDidMount() {
		this.getLibraryData();
	}

	getLibraryData() {
		axios.get(`/api/web/v1/library/${this.props.id}`).then((response) => {
			const cmd_decider_settings = JSON.parse(response.data.command_decider_settings);

			this.setState({
				folder: response.data.folder,
				priority: response.data.priority,
				fs_check_interval: response.data.fs_check_interval,
				path_masks: response.data.path_masks.join(","),
				queue: response.data.queue.Items,

				target_video_codec: cmd_decider_settings.target_video_codec,
				create_stereo_audio: cmd_decider_settings.create_stereo_audio,
				skip_hdr: cmd_decider_settings.skip_hdr,

				use_hardware: response.data.pipeline.use_hardware,
				hardware_codec: response.data.pipeline.hardware_codec,
				hw_device: response.data.pipeline.hw_device,
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
				<p className="text-center">Target Video Codec: {this.state.target_video_codec}</p>
				<p className="text-center">Create Stereo Audio Track: {(this.state.create_stereo_audio) ? "True" : "False"}</p>
				<p className="text-center">Skip HDR Files: {(this.state.skip_hdr) ? "True" : "False"}</p>
				{(this.state.use_hardware) ? <p className="text-center">Hardware Codec: {this.state.hardware_codec}</p> : null }
				{(this.state.use_hardware) ? <p className="text-center">Hardware Device: {this.state.hw_device}</p> : null }
				{(this.state.path_masks.length !== 0) ? <p className="text-center">Path Masks: {this.state.path_masks}</p> : null }
				<Button variant="secondary" onClick={() => {this.setState({showQueueModal: true})}}>Queue</Button>
				<Button variant="primary" onClick={() => {this.setState({showEditModal: true})}}>Edit</Button>
			</Card>
			{(this.state.showEditModal) ? (<EditLibraryModal
				show={true}
				closeHandler={() => { this.setState({showEditModal: false}); this.getLibraryData(); }}
				id={this.props.id}
				folder={this.state.folder}
				priority={this.state.priority}
				fs_check_interval={this.state.fs_check_interval}
				path_masks={this.state.path_masks}
				target_video_codec={this.state.target_video_codec}
				create_stereo_audio={this.state.create_stereo_audio}
				skip_hdr={this.state.skip_hdr}
				use_hardware={this.state.use_hardware}
				hardware_codec={this.state.hardware_codec}
				hw_device={this.state.hw_device}
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
	target_video_codec: string,
	create_stereo_audio: boolean,
	skip_hdr: boolean,
	use_hardware: boolean,
	hardware_codec: string,
	hw_device: string,
}

class CreateLibraryModal extends React.Component<ICreateLibraryModalProps, ICreateLibraryModalState> {
	constructor(props: ICreateLibraryModalProps) {
		super(props);

		this.state = {
			folder: "",
			priority: "",
			fs_check_interval: "",
			path_masks: "",
			target_video_codec: "HEVC",
			create_stereo_audio: true,
			skip_hdr: true,
			use_hardware: false,
			hardware_codec: "",
			hw_device: "",
		}

		this.submitLib = this.submitLib.bind(this);
	}

	submitLib(): void {
		let data = {
			folder: this.state.folder,
			priority: parseInt(this.state.priority),
			fs_check_interval: this.state.fs_check_interval,
			path_masks: this.state.path_masks.split(","),
			pipeline: {
				target_video_codec: this.state.target_video_codec,
				create_stereo_audio: this.state.create_stereo_audio,
				skip_hdr: this.state.skip_hdr,
				use_hardware: this.state.use_hardware,
				hardware_codec: this.state.hardware_codec,
				hw_device: this.state.hw_device,
			},
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

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Target Video Codec</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input no-box-shadow"
							as="select"
							custom
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ target_video_codec: event.target.value }); }}
							value={this.state.target_video_codec}
						>
							<option value="HEVC">H.265 (HEVC)</option>
							<option value="AVC">H.264 (AVC)</option>
							<option value="VP9">VP9</option>
						</FormControl>
					</InputGroup>

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Use Hardware</InputGroup.Text></InputGroup.Prepend>
						<InputGroup.Checkbox
							aria-label="Use Hardware Checkbox"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ use_hardware: event.target.checked }); }}
							checked={this.state.use_hardware}
						/>
					</InputGroup>

					{(this.state.use_hardware) ? <h6>WARNING: Hardware encoding is untested and highly experimental. Use at your own risk. <a href="https://github.com/BrenekH/encodarr/wiki/Hardware-Encoding" target="_blank" rel="noreferrer">More info.</a></h6> : null}

					{(this.state.use_hardware) ? <InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Hardware Codec</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder=""
							aria-label="Hardware Codec"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ hardware_codec: event.target.value }); }}
							value={this.state.hardware_codec}
						/>
					</InputGroup> : null}

					{(this.state.use_hardware) ? <InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Hardware Device</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder="/dev/dri/renderD128"
							aria-label="Hardware Device"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ hw_device: event.target.value }); }}
							value={this.state.hw_device}
						/>
					</InputGroup> : null}

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Create Stereo Audio Track</InputGroup.Text></InputGroup.Prepend>
						<InputGroup.Checkbox
							aria-label="Create Stereo Audio Track Checkbox"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ create_stereo_audio: event.target.checked }); }}
							checked={this.state.create_stereo_audio}
						/>
					</InputGroup>

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Skip HDR</InputGroup.Text></InputGroup.Prepend>
						<InputGroup.Checkbox
							aria-label="Skip HDR Checkbox"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ skip_hdr: event.target.checked }); }}
							checked={this.state.skip_hdr}
						/>
					</InputGroup>

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
	target_video_codec: string,
	create_stereo_audio: boolean,
	skip_hdr: boolean,
	use_hardware: boolean,
	hardware_codec: string,
	hw_device: string,
}

interface IEditLibraryModalState {
	folder: string,
	priority: string,
	fs_check_interval: string,
	path_masks: string,
	target_video_codec: string,
	create_stereo_audio: boolean,
	skip_hdr: boolean,
	use_hardware: boolean,
	hardware_codec: string,
	hw_device: string,
}

class EditLibraryModal extends React.Component<IEditLibraryModalProps, IEditLibraryModalState> {
	constructor(props: IEditLibraryModalProps) {
		super(props);

		this.state = {
			folder: props.folder,
			priority: props.priority,
			fs_check_interval: props.fs_check_interval,
			path_masks: props.path_masks,
			target_video_codec: props.target_video_codec,
			create_stereo_audio: props.create_stereo_audio,
			skip_hdr: props.skip_hdr,
			use_hardware: props.use_hardware,
			hardware_codec: props.hardware_codec,
			hw_device: props.hw_device,
		}

		this.putChanges = this.putChanges.bind(this);
		this.deleteLibrary = this.deleteLibrary.bind(this);
	}

	putChanges(): void {
		let data = {
			folder: this.state.folder,
			priority: parseInt(this.state.priority),
			fs_check_interval: this.state.fs_check_interval,
			path_masks: this.state.path_masks.split(","),
			command_decider_settings: JSON.stringify({
				target_video_codec: this.state.target_video_codec,
				create_stereo_audio: this.state.create_stereo_audio,
				skip_hdr: this.state.skip_hdr,
				use_hardware: this.state.use_hardware,
				hardware_codec: this.state.hardware_codec,
				hw_device: this.state.hw_device,
			}),
		};
		axios.put(`/api/web/v1/library/${this.props.id}`, data).then(() => {
			this.props.closeHandler();
		}).catch((error) => {
			console.error(`/api/web/v1/library/${this.props.id} failed with error: ${error}`)
		});
	}

	deleteLibrary(): void {
		axios.delete(`/api/web/v1/library/${this.props.id}`).then(() => { this.props.closeHandler(); }).catch((error) => {
			console.error(`/api/web/v1/library/${this.props.id} failed with error: ${error}`);
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
						<InputGroup.Prepend><InputGroup.Text>Target Video Codec</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input no-box-shadow"
							as="select"
							custom
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ target_video_codec: event.target.value }); }}
							value={this.state.target_video_codec}
						>
							<option value="HEVC">H.265 (HEVC)</option>
							<option value="AVC">H.264 (AVC)</option>
							<option value="VP9">VP9</option>
						</FormControl>
					</InputGroup>

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Use Hardware</InputGroup.Text></InputGroup.Prepend>
						<InputGroup.Checkbox
							aria-label="Use Hardware Checkbox"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ use_hardware: event.target.checked }); }}
							checked={this.state.use_hardware}
						/>
					</InputGroup>

					{(this.state.use_hardware) ? <h6>WARNING: Hardware encoding is untested and highly experimental. Use at your own risk. <a href="https://github.com/BrenekH/encodarr/wiki/Hardware-Encoding" target="_blank" rel="noreferrer">More info.</a></h6> : null}

					{(this.state.use_hardware) ? <InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Hardware Codec</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder=""
							aria-label="Hardware Codec"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ hardware_codec: event.target.value }); }}
							value={this.state.hardware_codec}
						/>
					</InputGroup> : null}

					{(this.state.use_hardware) ? <InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Hardware Device</InputGroup.Text></InputGroup.Prepend>
						<FormControl
							className="dark-text-input"
							placeholder="/dev/dri/renderD128"
							aria-label="Hardware Device"
							aria-describedby="basic-addon1"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ hw_device: event.target.value }); }}
							value={this.state.hw_device}
						/>
					</InputGroup> : null}

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Create Stereo Audio Track</InputGroup.Text></InputGroup.Prepend>
						<InputGroup.Checkbox
							aria-label="Create Stereo Audio Track Checkbox"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ create_stereo_audio: event.target.checked }); }}
							checked={this.state.create_stereo_audio}
						/>
					</InputGroup>

					<InputGroup className="mb-3">
						<InputGroup.Prepend><InputGroup.Text>Skip HDR</InputGroup.Text></InputGroup.Prepend>
						<InputGroup.Checkbox
							aria-label="Skip HDR Checkbox"
							onChange={(event: React.ChangeEvent<HTMLInputElement>) => { this.setState({ skip_hdr: event.target.checked }); }}
							checked={this.state.skip_hdr}
						/>
					</InputGroup>

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
					<Button className="delete-button" variant="danger" onClick={this.deleteLibrary}>Delete</Button>
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
		let propsQueue = this.props.queue;
		if (propsQueue === null) {
			propsQueue = [];
		}

		const qList = propsQueue.map((v: IQueuedJob, i: number) => {
			return <TableEntry key={v.uuid} index={i+1} path={v.path} command={v.command.join(" ")}/>;
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
								<th scope="col">Cmd</th>
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
	command: Array<string>,
}

interface ITableEntryProps {
	index: number,
	path: string,
	command: string,
}

function TableEntry(props: ITableEntryProps) {
	return (<tr>
		<th scope="row">{props.index}</th>
		<td>{props.path}</td>
		<td>
			<TerminalIcon title={props.command}/>
		</td>
	</tr>);
}
