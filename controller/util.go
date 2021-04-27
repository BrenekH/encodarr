package controller

import "context"

// IsContextFinished returns a boolean indicating whether or not a context.Context is finished.
// This replaces the need to use a select code block.
func IsContextFinished(ctx *context.Context) bool {
	select {
	case <-(*ctx).Done():
		return true
	default:
		return false
	}
}
