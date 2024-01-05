package middlewares

import "context"

func ContextWrapAll(ctx context.Context, valueMap map[string]any) context.Context {
	for k, v := range valueMap {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}
