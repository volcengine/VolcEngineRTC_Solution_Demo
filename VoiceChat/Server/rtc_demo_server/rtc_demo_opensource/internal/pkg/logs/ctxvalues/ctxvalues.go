package ctxvalues

import "context"

const CtxLogID = "LogID"

func SetLogID(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, CtxLogID, logID)
}

func LogID(ctx context.Context) (string, bool) {
	return getStringFromContext(ctx, CtxLogID)
}

func getStringFromContext(ctx context.Context, key string) (string, bool) {
	if ctx == nil {
		return "", false
	}

	v := ctx.Value(key)
	if v == nil {
		return "", false
	}

	switch v := v.(type) {
	case string:
		return v, true
	case *string:
		if v == nil {
			return "", false
		}
		return *v, true
	}
	return "", false
}
