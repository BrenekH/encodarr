import { SvgImage } from "./SvgImage";

import headphones from "./headphones.svg";

export function AudioImage() {
	return (<SvgImage
		location={headphones}
		alt="Headphones"
		title="An additional stereo audio track will be created"
	/>);
}
