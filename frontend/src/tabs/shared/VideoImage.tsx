import { SvgImage } from "./SvgImage";

import playButton from "./play_button.svg";

export function VideoImage() {
	return (<SvgImage
		location={playButton}
		alt="Play Button"
		title="File will be encoded to HEVC"
	/>);
}
