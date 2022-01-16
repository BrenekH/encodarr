import TerminalIconSvg from "./TerminalIconSvg";

export interface ITerminalIconProps {
	title: string,
	height?: string,
	width?: string,
}

export default function TerminalIcon(props: ITerminalIconProps) {
	return (<div style={{ display: "inline" }} title={props.title}>
		<TerminalIconSvg height={props.height} width={props.width} />
	</div>);
}
