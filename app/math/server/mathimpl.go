package server

import (
	"context"
	"log"
	"math"

	"github.com/sknv/micronats/app/lib/xnats"
	"github.com/sknv/micronats/app/lib/xnats/status"
	"github.com/sknv/micronats/app/math/rpc"
)

type MathImpl struct{}

func (*MathImpl) Circle(ctx context.Context, args *rpc.CircleArgs) (*rpc.CircleReply, error) {
	if args.Radius <= 0 {
		return nil, status.Error(status.InvalidArgument, "radius must be a positive number")
	}

	log.Printf("[INFO] circle meta foo: %s", xnats.MetaValue(ctx, "foo")) // access sample metadata

	return &rpc.CircleReply{
		Length: 2 * math.Pi * args.Radius,
		Square: math.Pi * args.Radius * args.Radius,
	}, nil
}

func (*MathImpl) Rect(_ context.Context, args *rpc.RectArgs) (*rpc.RectReply, error) {
	if args.Width <= 0 || args.Height <= 0 {
		return nil, status.Error(status.InvalidArgument, "width and height must be positive numbers")
	}

	return &rpc.RectReply{
		Perimeter: 2*args.Width + 2*args.Height,
		Square:    args.Width * args.Height,
	}, nil
}
