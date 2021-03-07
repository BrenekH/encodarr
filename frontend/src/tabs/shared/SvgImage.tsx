import "./SvgImage.css";

interface ISvgImageProps {
	location: string,
	alt: string,
	title: string,
}

export function SvgImage(props: ISvgImageProps) {
	return (<img
		className="queue-icon"
		src={props.location}
		alt={props.alt}
		height="20px"
		title={props.title}></img>);
}
