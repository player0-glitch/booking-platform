package core

import "net/http"

// let the core pass around the auth middleware
type Middleware func(http.Handler) http.Handler
