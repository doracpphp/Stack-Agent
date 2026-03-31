package task

import (
	"fmt"
	"time"

	appctx "stack-agent/context"
)

// GetCurrentTimeHandler returns the current date and time as a string.
//
// Params:
//
//	format  string  - Go time format string (default: "2006-01-02 15:04:05")
//	timezone string - IANA timezone name, e.g. "Asia/Tokyo" (default: local)
type GetCurrentTimeHandler struct{}

func (h *GetCurrentTimeHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	format := paramStringRaw(params, "format")
	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	tz := paramStringRaw(params, "timezone")
	now := time.Now()
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			return fmt.Errorf("get_current_time: invalid timezone %q: %w", tz, err)
		}
		now = now.In(loc)
	}

	result := now.Format(format)
	fmt.Printf("  → Current time: %s\n", result)
	ctx.Set(taskID, &appctx.Content{Text: result})
	return nil
}
