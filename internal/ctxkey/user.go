package ctxkey

import "context"

type contextKey string

const userIDKey contextKey = "userID"

func SetUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
