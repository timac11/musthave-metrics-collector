package middleware

type Middleware struct {
	signingKey string
}

func NewMiddleware(hashingKey string) *Middleware {
	return &Middleware{signingKey: hashingKey}
}
