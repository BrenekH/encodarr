import axios from "axios";
import React from "react";

import "./AboutSection.css";

interface IAboutSectionState {
	controller_version: string,
	web_api_versions: Array<string>,
	runner_api_versions: Array<string>,
}

export default class AboutSection extends React.Component<{}, IAboutSectionState> {
	constructor (props: any) {
		super(props);

		this.state = {
			controller_version: "Could not contact a ring",
			web_api_versions: [],
			runner_api_versions: [],
		};
	}

	componentDidMount() {
		// Get api supported versions
		axios.get("/api").then((response) => {
			this.setState({
				web_api_versions: response.data.web.versions,
				runner_api_versions: response.data.runner.versions,
			});
		}).catch((error) => {
			console.error(`Request to /api failed with error: ${error}`);
		});

		// Get controller version
		axios.get("/version").then((response) => {
			this.setState({
				controller_version: response.data,
			});
		}).catch((error) => {
			console.error(`Request to /api failed with error: ${error}`);
		});
	}

	render() {
		return (<>
			<h5>About Encodarr</h5>

			<p><b>License:</b> This project is licensed under the Mozilla Public License 2.0 a copy of which can be found <a href="https://github.com/BrenekH/encodarr/blob/master/LICENSE" target="_blank" rel="noreferrer">here</a></p>

			<p><b>Controller Version:</b> {this.state.controller_version}</p>
			<p className="list-title"><b>Supported API Versions:</b></p>
			<ul className="api-list">
				<li><b>Web:</b> {this.state.web_api_versions.join(", ")}</li>
				<li><b>Runner:</b> {this.state.runner_api_versions.join(", ")}</li>
			</ul>

			<p><b>GitHub Repository:</b> <a href="https://github.com/BrenekH/encodarr" target="_blank" rel="noreferrer">https://github.com/BrenekH/encodarr</a></p>
		</>);
	}
}
