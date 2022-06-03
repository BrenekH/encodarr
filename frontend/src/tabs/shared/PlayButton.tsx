import { SVGProps } from "react"

const PlayButton = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={48}
    height={48}
    viewBox="0 0 12.7 12.7"
    {...props}
  >
    <path
      style={{
        fill: "#0f0",
        fillRule: "evenodd",
        strokeWidth: 0.264583,
      }}
      d="M11.628 6.348 1.302 12.309V.386z"
      transform="matrix(1.22077 0 0 1.05866 -1.524 -.39)"
    />
  </svg>
)

export default PlayButton
