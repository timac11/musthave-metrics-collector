package middleware

type Middleware struct {
	hashingKey string
}

func NewMiddleware(hashingKey string) *Middleware {
	return &Middleware{hashingKey: hashingKey}
}
