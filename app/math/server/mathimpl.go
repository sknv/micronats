package server

import (
	"context"
	"math"

	"github.com/sknv/micronats/app/lib/xnats"
	"github.com/sknv/micronats/app/math/rpc"
)

type MathImpl struct{}

func (*MathImpl) Circle(_ context.Context, args *rpc.CircleArgs) (*rpc.CircleReply, error) {
	if args.Radius <= 0 {
		return nil, xnats.ErrorStatus(xnats.StatusInvalidArgument, "radius must be a positive number")
	}

	return &rpc.CircleReply{
		Length: 2 * math.Pi * args.Radius,
		Square: math.Pi * args.Radius * args.Radius,
	}, nil
}

func (*MathImpl) Rect(_ context.Context, args *rpc.RectArgs) (*rpc.RectReply, error) {
	if args.Width <= 0 || args.Height <= 0 {
		return nil, xnats.ErrorStatus(xnats.StatusInvalidArgument, "width and height must be positive numbers")
	}

	return &rpc.RectReply{
		Perimeter: 2*args.Width + 2*args.Height,
		Square:    args.Width * args.Height,
	}, nil
}
