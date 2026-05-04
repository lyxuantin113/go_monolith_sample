package auth

import "context"

func GetUserIDFromContext(ctx context.Context) uint {
	if id, ok := ctx.Value("user_id").(uint); ok {
		return id
	}
	return 0
}
