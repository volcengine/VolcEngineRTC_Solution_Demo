package util

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs/ctxvalues"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs/logid"
)

// EnsureID ensures log id context
func EnsureID(ctx context.Context) context.Context {
	if len(RetrieveID(ctx)) == 0 {
		return SetID(ctx, logid.GenLogID())
	}
	return ctx
}

// SetID sets log id to context
func SetID(ctx context.Context, id string) context.Context {
	return ctxvalues.SetLogID(ctx, id)
}

// RetrieveID retrieves log if from context
func RetrieveID(ctx context.Context) string {
	logID, _ := ctxvalues.LogID(ctx)
	return logID
}
