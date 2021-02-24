package server

func inferMIMETypeFromExt(ext string) string {
	switch ext {
	case "mp4":
		return "video/mp4"
	case "m4v":
		return "video/m4v"
	case "avi":
		return "video/x-msvideo"
	case "mov", "qt":
		return "video/quicktime"
	case "wmv":
		return "video/x-ms-wmv"
	case "mkv":
		return "video/x-matroska"
	case "ogg":
		return "application/ogg"
	case "webm":
		return "video/webm"
	case "m4p":
		return "application/octet-stream"
	default:
		return "application/octet-stream"
	}
}
