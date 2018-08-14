package xnats

import (
	"context"
)

type metadataContextKeyType string
type metadata map[string]string

const (
	metadataContextKey metadataContextKeyType = "xnats.metadata"
)

func Metadata(ctx context.Context) metadata {
	meta := ctx.Value(metadataContextKey)
	if meta == nil {
		return nil
	}
	return meta.(metadata)
}

func WithMetaValue(ctx context.Context, key, value string) context.Context {
	meta := Metadata(ctx)
	if meta == nil { // create a new map if not exist
		meta = make(map[string]string)
	}
	meta[key] = value // upsert the value
	return context.WithValue(ctx, metadataContextKey, meta)
}

func MetaValue(ctx context.Context, key string) string {
	meta := Metadata(ctx)
	if meta == nil {
		return ""
	}
	return meta[key]
}
