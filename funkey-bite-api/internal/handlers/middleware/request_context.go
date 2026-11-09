package middleware

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/requestctx"
)

// maxPeekBodyBytes bounds how much of a request body middleware reads.
const maxPeekBodyBytes = 1 << 20

var errBodyTooLarge = errors.New("request body too large")

// RequestContext returns the request's context carrying who is making it
// (admin, signed-in customer or guest, with IP and request ID) and its
// idempotency key, for services to record in audit entries and idempotency
// records. Call it in a handler, after authentication has run.
func RequestContext(c *gin.Context) context.Context {
	return requestctx.WithActor(c.Request.Context(), actorFromGin(c))
}

func actorFromGin(c *gin.Context) requestctx.Actor {
	actor := requestctx.Actor{
		Type:      requestctx.ActorCustomer,
		IP:        c.ClientIP(),
		RequestID: GetRequestID(c),
	}
	if userID, ok := c.Get("user_id"); ok {
		if id, ok := userID.(int); ok {
			actor.ID = strconv.Itoa(id)
		}
		if isAdmin, ok := c.Get("is_admin"); ok && isAdmin == true {
			actor.Type = requestctx.ActorAdmin
		}
	}
	return actor
}

// peekBody reads the request body and puts it back so the handler can still
// bind it.
func peekBody(c *gin.Context) ([]byte, error) {
	if c.Request.Body == nil {
		return nil, nil
	}
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxPeekBodyBytes+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxPeekBodyBytes {
		return nil, errBodyTooLarge
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}
