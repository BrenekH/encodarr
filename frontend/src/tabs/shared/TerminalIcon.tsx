import TerminalIconSvg from "./TerminalIconSvg";

export interface ITerminalIconProps {
	title: string,
}

export default function TerminalIcon(props: ITerminalIconProps) {
	return (<div style={{ display: "inline" }} title={props.title}>
		<TerminalIconSvg />
	</div>);
}
