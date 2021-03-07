import axios from "axios";
import React from "react";
import Button from "react-bootstrap/Button";
import InputGroup from "react-bootstrap/InputGroup";
import FormControl from "react-bootstrap/FormControl";

import "./SettingsTab.css";
import "../spacers.css";

interface IInputValues {
	fsCheckInterval: string,
	smallerFilesCheck: Boolean,
	healthCheckInterval: string,
	unresponsiveRunnerTimeout: string,
	logVerbosity: string
}

interface ISettingsTabState {
	inputValues: IInputValues,
	showSavedIndicator: Boolean,
}

export class SettingsTab extends React.Component<any, ISettingsTabState> {
	constructor(props: any) {
		super(props);
		this.state = {
			inputValues: {
				"fsCheckInterval": "",
				"smallerFilesCheck": false,
				"healthCheckInterval": "",
				"unresponsiveRunnerTimeout": "",
				"logVerbosity": "",
			},
			showSavedIndicator: false,
		}

		this.handleClick = this.handleClick.bind(this);
	}

	componentDidMount(): void {
		this.updateSettings();
	}

	createChangeHandler(id: string, checkChecked = false) {
		let f = (event: React.ChangeEvent<HTMLInputElement>) => {
			const currentValues: any = Object.assign({}, this.state.inputValues);
			if (checkChecked) {
				currentValues[id] = event.target.checked;
			} else {
				currentValues[id] = event.target.value;
			}
			this.setState({
				inputValues: currentValues,
			});
		};
		f.bind(this);
		return f;
	}

	handleClick(): void {
		axios.put("/api/web/v1/settings", {
			"FileSystemCheckInterval": this.state.inputValues.fsCheckInterval,
			"SmallerFiles": this.state.inputValues.smallerFilesCheck,
			"HealthCheckInterval": this.state.inputValues.healthCheckInterval,
			"HealthCheckTimeout": this.state.inputValues.unresponsiveRunnerTimeout,
			"LogVerbosity": this.state.inputValues.logVerbosity,
		}).then((response) => {
			if (response.status >= 200 && response.status <= 299) {
				this.setState({
					showSavedIndicator: true,
				});
			} else {
				console.error(response);
			}

			this.updateSettings(); // The settings are updated after saving to clear any invalid values
		});
	}

	updateSettings(): void {
		axios.get("/api/web/v1/settings").then((response) => {
			this.setState({
				inputValues: {
					"fsCheckInterval": response.data.FileSystemCheckInterval,
					"smallerFilesCheck": response.data.SmallerFiles,
					"healthCheckInterval": response.data.HealthCheckInterval,
					"unresponsiveRunnerTimeout": response.data.HealthCheckTimeout,
					"logVerbosity": response.data.LogVerbosity,
				}
			});
		});
	}

	render(): React.ReactNode {
		const savedIndicator = (this.state.showSavedIndicator) ? <SavedIndicator /> : null;
		if (this.state.showSavedIndicator) {
			setTimeout(() => {this.setState({
				showSavedIndicator: false,
			});}, 5000);
		}

		return (<div>
			<h5>General</h5>

			<InputGroup className="mb-3">
				<InputGroup.Prepend>
					<InputGroup.Text>File System Scan Interval</InputGroup.Text>
				</InputGroup.Prepend>
				<FormControl
					className="dark-text-input"
					placeholder="0h0m0s"
					aria-label="fs-check-interval"
					aria-describedby="basic-addon1"
					onChange={this.createChangeHandler("fsCheckInterval")}
					value={this.state.inputValues["fsCheckInterval"]}
				/>
			</InputGroup>

			<InputGroup className="mb-3">
				<InputGroup.Prepend>
					<InputGroup.Text>Prefer Smaller Files</InputGroup.Text>
					<InputGroup.Checkbox
						aria-label="Prefer Smaller Files Checkbox"
						onChange={this.createChangeHandler("smallerFilesCheck", true)}
						value={this.state.inputValues["smallerFilesCheck"]}
					/>
				</InputGroup.Prepend>
				<FormControl
					className="dark-text-input"
					value="Only keep encoded files if they are smaller than the original."
					disabled
					readOnly
				/>
			</InputGroup>

			<div className="spacer"></div>

			<h5>Runner Health</h5>

			<InputGroup className="mb-3">
				<InputGroup.Prepend>
					<InputGroup.Text>Runner Health Check Interval</InputGroup.Text>
				</InputGroup.Prepend>
				<FormControl
					className="dark-text-input"
					placeholder="0h0m0s"
					aria-label="health-check-interval"
					aria-describedby="basic-addon1"
					onChange={this.createChangeHandler("healthCheckInterval")}
					value={this.state.inputValues["healthCheckInterval"]}
				/>
			</InputGroup>

			<InputGroup className="mb-3">
				<InputGroup.Prepend>
					<InputGroup.Text>Unresponsive Runner Timeout</InputGroup.Text>
				</InputGroup.Prepend>
				<FormControl
					className="dark-text-input"
					placeholder="0h0m0s"
					aria-label="unresponsive-runner-timeout"
					aria-describedby="basic-addon1"
					onChange={this.createChangeHandler("unresponsiveRunnerTimeout")}
					value={this.state.inputValues["unresponsiveRunnerTimeout"]}
				/>
			</InputGroup>

			<div className="spacer"></div>

			<h5>Logging</h5>

			<InputGroup className="mb-3">
				<InputGroup.Prepend>
					<InputGroup.Text>Log Verbosity</InputGroup.Text>
				</InputGroup.Prepend>
				<FormControl
					className="dark-text-input no-box-shadow"
					as="select"
					custom
					onChange={this.createChangeHandler("logVerbosity")}
					value={this.state.inputValues["logVerbosity"]}
				>
					<option value="TRACE">Trace</option>
					<option value="DEBUG">Debug</option>
					<option value="INFO">Info</option>
					<option value="WARNING">Warning</option>
					<option value="ERROR">Error</option>
					<option value="CRITICAL">Critical</option>
				</FormControl>
			</InputGroup>

			<div className="smol-spacer"></div>

			<Button variant="light" onClick={this.handleClick}>Save</Button>
			{savedIndicator}
		</div>);
	}
}

function SavedIndicator() {
	return <p className="pop-in-out" style={{display: "inline"}}>Saved!</p>;
}
