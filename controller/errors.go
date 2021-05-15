package controller

import "errors"

// ErrEmptyQueue represents when the operation cannot be completed because the queue is empty
var ErrEmptyQueue error = errors.New("queue is empty")

var ErrClosed = errors.New("attempted operation on closed struct")
