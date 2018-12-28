package core

import (
	"time"

	auth "github.com/sknv/micronats/app/auth/rpc"
)

const (
	RPCTimeout = 30 * time.Second
)

type Broker struct {
	Auth auth.Auth
}
