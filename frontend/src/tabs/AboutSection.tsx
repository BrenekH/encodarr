import axios from "axios";
import React from "react";

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

			<p>License: This project is licensed under the GNU Public License v2</p>

			<p>Controller Version: {this.state.controller_version}</p>
			<p>Supported Web API Versions: {this.state.web_api_versions.join(", ")}</p>
			<p>Supported Runner API Versions: {this.state.runner_api_versions.join(", ")}</p>

			<p>GitHub Repository: <a href="https://github.com/BrenekH/encodarr">https://github.com/BrenekH/encodarr</a></p>
		</>);
	}
}
