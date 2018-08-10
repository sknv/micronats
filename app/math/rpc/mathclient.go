package rpc

// import (
// 	"context"

// 	"github.com/nats-io/go-nats"

// 	"github.com/sknv/micronats/app/lib/xnats/rpc"
// )

// type MathClient struct {
// 	*rpc.ProtoClient
// }

// func NewMathClient(natsConn *nats.Conn) Math {
// 	return &MathClient{ProtoClient: rpc.NewProtoClient(natsConn)}
// }

// func (c *MathClient) Circle(ctx context.Context, args *CircleArgs) (*CircleReply, error) {
// 	reply := new(CircleReply)
// 	if err := c.Call(ctx, CirclePattern, args, reply); err != nil {
// 		return nil, err
// 	}
// 	return reply, nil
// }

// func (c *MathClient) Rect(ctx context.Context, args *RectArgs) (*RectReply, error) {
// 	reply := new(RectReply)
// 	if err := c.Call(ctx, RectPattern, args, reply); err != nil {
// 		return nil, err
// 	}
// 	return reply, nil
// }
