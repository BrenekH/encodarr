import { SvgImage } from "./SvgImage";

import terminalIcon from "./terminalIcon.svg";

export interface ITerminalIconProps {
	title: string,
}

export default function TerminalIcon(props: ITerminalIconProps) {
	return (<SvgImage
		location={terminalIcon}
		alt="Terminal"
		title={props.title}
	/>);
}
