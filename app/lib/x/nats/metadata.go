package nats

import (
	"context"
)

type Metadata map[string]string

type ctxMeta int

const (
	ctxKeyMeta ctxMeta = 0
)

func ContextWithMetaValue(ctx context.Context, key, value string) context.Context {
	meta := MetadataFrom(ctx)
	if meta == nil { // create a new map if not exist
		meta = make(Metadata)
	}
	meta[key] = value // upsert the value
	return context.WithValue(ctx, ctxKeyMeta, meta)
}

func MetadataFrom(ctx context.Context) Metadata {
	meta := ctx.Value(ctxKeyMeta)
	if meta == nil {
		return nil
	}
	return meta.(Metadata)
}

func MetaValueFrom(ctx context.Context, key string) string {
	meta := MetadataFrom(ctx)
	if meta == nil {
		return ""
	}
	return meta[key]
}
