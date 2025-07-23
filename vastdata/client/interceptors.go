// Copyright (c) HashiCorp, Inc.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	vast_client "github.com/vast-data/go-vast-client"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
)

type contextKey string

const requestIDKey contextKey = "@request_id"

// reqIDCounter holds a globally shared atomic counter used to generate unique request IDs.
// It is initialized with a random value in the lower half of uint32 range to avoid predictable sequences.
var reqIDCounter = rand.Uint32() % (uint32(math.MaxUint32/2) + 1) // result in [0, max]

// BeforeRequestFnCallback logs the HTTP request being sent.
// It reads and optionally compacts the body (if present) for structured logging,
// and includes the request ID from context (if available).
// For more details see: https://github.com/vast-data/go-vast-client
func BeforeRequestFnCallback(ctx context.Context, _ *http.Request, verb, url string, body io.Reader) error {
	var logMsg strings.Builder
	uid, _ := ctx.Value(requestIDKey).(string)
	logMsg.WriteString(fmt.Sprintf(": ➤  start: req_id=%s - [%s] %s", uid, verb, url))

	if body != nil {
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			tflog.Error(ctx, "failed to read request body", map[string]any{
				"error": err.Error(),
			})
			return err
		}

		trimmed := bytes.TrimSpace(bodyBytes)
		if len(trimmed) > 0 && !bytes.Equal(trimmed, []byte("null")) {
			var compact bytes.Buffer
			if err := json.Compact(&compact, trimmed); err == nil {
				logMsg.WriteString(fmt.Sprintf(" - body: %s", compact.String()))
			} else {
				logMsg.WriteString(fmt.Sprintf(" - body (raw): %s", string(trimmed)))
			}
		}
	}

	tflog.Info(ctx, logMsg.String())
	return nil
}

// AfterRequestFnCallback logs the response received from the HTTP request.
// It uses the response's PrettyTable method to render a formatted table,
// and includes the request ID from context.
// For more details see: https://github.com/vast-data/go-vast-client
func AfterRequestFnCallback(ctx context.Context, response vast_client.Renderable) (vast_client.Renderable, error) {
	uid, _ := ctx.Value(requestIDKey).(string)
	tflog.Info(ctx, fmt.Sprintf("%s end: req_id=%s | ", response.PrettyTable(), uid))
	return response, nil
}

// ContextWithRequestID returns a new context containing a generated request ID.
// The ID is a hex string based on a global atomic counter to ensure uniqueness
// across concurrent requests within the same process.
func ContextWithRequestID(ctx context.Context) context.Context {
	newID := atomic.AddUint32(&reqIDCounter, 1)
	return context.WithValue(ctx, requestIDKey, fmt.Sprintf("0x%08x", newID))
}
