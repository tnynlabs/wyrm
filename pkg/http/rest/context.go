package rest

// UserCtxKey should be used to get/set the authenticated user
// instance in the request context if it exists.
// Example: ctx.Value(UserCtxKey{}).
type UserCtxKey struct{}
