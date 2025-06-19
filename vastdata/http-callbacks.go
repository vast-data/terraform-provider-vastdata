package vastdata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"strings"
)

func BeforeRequestFnCallback(ctx context.Context, r *http.Request, verb, url string, body io.Reader) error {
	var logMsg strings.Builder

	logMsg.WriteString(fmt.Sprintf(":start ======= [%s] %s", verb, url))

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
				logMsg.WriteString(fmt.Sprintf("\nBody: %s", compact.String()))
			} else {
				logMsg.WriteString(fmt.Sprintf("\nBody (raw): %s", string(trimmed)))
			}
		}
	}

	tflog.Info(ctx, logMsg.String())
	return nil
}

func AfterRequestFnCallback(ctx context.Context, response Renderable) (Renderable, error) {
	tflog.Info(ctx, fmt.Sprintf("%s :end ============", response.PrettyTable()))
	return response, nil
}
