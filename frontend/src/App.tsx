import React from "react";
import Tab from "react-bootstrap/Tab";
import Nav from "react-bootstrap/Nav";

import { RunningTab } from "./tabs/RunningTab";
import { QueueTab } from "./tabs/QueueTab";
import { HistoryTab } from "./tabs/HistoryTab";
import { SettingsTab } from "./tabs/SettingsTab";
import './spacers.css';

function Title() {
	return <div className="header-content text-center"><h1>Project RedCedar</h1></div>
}

class App extends React.Component {
	handleSelect(eventKey: any) {
		switch (eventKey) {
			case "queue":
				window.history.replaceState(undefined, "", "/queue");
				document.title = "Queue - Project RedCedar";
				break;
			case "history":
				window.history.replaceState(undefined, "", "/history");
				document.title = "History - Project RedCedar";
				break;
			case "settings":
				window.history.replaceState(undefined, "", "/settings");
				document.title = "Settings - Project RedCedar";
				break;
			case "running":
				window.history.replaceState(undefined, "", "/running");
				document.title = "Running - Project RedCedar";
				break;
			default:
				break;
		}
	}

	render() {
		let activeKey: string = "running";
		switch (window.location.pathname) {
			case "/queue":
				activeKey = "queue";
				break;
			case "/history":
				activeKey = "history";
				break;
			case "/settings":
				activeKey = "settings";
				break;
			default:
				break;
		}

		return (<div className="container">
			<Title />
			<Tab.Container id="tab-nav" defaultActiveKey={activeKey} transition={false} onSelect={this.handleSelect}>
				<Nav fill variant="pills">
					<Nav.Item>
						<Nav.Link eventKey="running">Running</Nav.Link>
					</Nav.Item>
					<Nav.Item>
						<Nav.Link eventKey="queue">Queue</Nav.Link>
					</Nav.Item>
					<Nav.Item>
						<Nav.Link eventKey="history">History</Nav.Link>
					</Nav.Item>
					<Nav.Item>
						<Nav.Link eventKey="settings">Settings</Nav.Link>
					</Nav.Item>
				</Nav>

				<div className="spacer"></div>

				<Tab.Content>
					<Tab.Pane eventKey="running" mountOnEnter={true} unmountOnExit={true}>
						<RunningTab />
					</Tab.Pane>
					<Tab.Pane eventKey="queue" mountOnEnter={true} unmountOnExit={true}>
						<QueueTab />
					</Tab.Pane>
					<Tab.Pane eventKey="history" mountOnEnter={true} unmountOnExit={true}>
						<HistoryTab />
					</Tab.Pane>
					<Tab.Pane eventKey="settings" mountOnEnter={true} unmountOnExit={true}>
						<SettingsTab />
					</Tab.Pane>
				</Tab.Content>
			</Tab.Container>

			<div className="smol-spacer"></div>
		</div>);
	}
}

export default App;