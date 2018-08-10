package rpc

import (
	"context"
)

const (
	CircleSubject = "/rpc/math/circle"
	RectSubject   = "/rpc/math/rect"
)

type Math interface {
	Circle(context.Context, *CircleArgs) (*CircleReply, error)
	Rect(context.Context, *RectArgs) (*RectReply, error)
}

type (
	CircleArgs struct {
		Radius float64 `json:"radius,omitempty"`
	}

	CircleReply struct {
		Length float64 `json:"length,omitempty"`
		Square float64 `json:"square,omitempty"`
	}
)

type (
	RectArgs struct {
		Width  float64 `json:"width,omitempty"`
		Height float64 `json:"height,omitempty"`
	}

	RectReply struct {
		Perimeter float64 `json:"perimeter,omitempty"`
		Square    float64 `json:"square,omitempty"`
	}
)
