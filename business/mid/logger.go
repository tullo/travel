package mid

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dgraph-io/travel/foundation/web"
)

// Logger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			// Call the next handler.
			err := handler(ctx, w, r)

			log.Printf("%s : (%d) : %s %s -> %s (%s)",
				v.TraceID, v.StatusCode,
				r.Method, r.URL.Path,
				r.RemoteAddr, time.Since(v.Now),
			)

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return m
}
