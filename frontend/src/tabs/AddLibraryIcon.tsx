import { SVGProps } from "react"

const AddLibraryIcon = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="150mm"
    height="150mm"
    viewBox="0 0 150 150"
    {...props}
  >
    <rect
      style={{
        fill: "#fff",
        fillOpacity: 1,
        strokeWidth: 4.16469,
        stopColor: "#000",
      }}
      width={50}
      height={150}
      x={50}
      y={-150}
      rx={0}
      ry={0}
      transform="scale(1 -1)"
    />
    <rect
      style={{
        fill: "#fff",
        fillOpacity: 1,
        strokeWidth: 4.1647,
        stopColor: "#000",
      }}
      width={50}
      height={150}
      x={50}
      rx={0}
      ry={0}
      transform="matrix(0 1 1 0 0 0)"
    />
  </svg>
)

export default AddLibraryIcon
