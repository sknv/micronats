package server

import (
	"context"
	"log"

	"github.com/sknv/micronats/app/lib/xnats"
	"github.com/sknv/micronats/app/lib/xnats/status"
	math "github.com/sknv/micronats/app/math/rpc"
)

type MathImpl struct{}

func (*MathImpl) Circle(ctx context.Context, args *math.CircleArgs) (*math.CircleReply, error) {
	if args.Radius <= 0 {
		return nil, status.Error(status.InvalidArgument, "radius must be a positive number")
	}

	log.Printf("[INFO] circle meta foo: %s", xnats.MetaValue(ctx, "foo")) // access sample metadata

	pi := 3.1416 // there is math.Pi constant in the standard lib btw
	return &math.CircleReply{
		Length: 2 * pi * args.Radius,
		Square: pi * args.Radius * args.Radius,
	}, nil
}

func (*MathImpl) Rect(_ context.Context, args *math.RectArgs) (*math.RectReply, error) {
	if args.Width <= 0 || args.Height <= 0 {
		return nil, status.Error(status.InvalidArgument, "width and height must be positive numbers")
	}

	return &math.RectReply{
		Perimeter: 2*args.Width + 2*args.Height,
		Square:    args.Width * args.Height,
	}, nil
}
