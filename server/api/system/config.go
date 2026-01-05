package system

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// GetValue 获取配置
func GetValue(app core.App, key string) string {
	record, err := app.FindFirstRecordByData("config", "key", key)

	if err != nil {
		return ""
	}

	value, err := record.Get("value").(types.JSONRaw).Value()

	if err != nil {
		return ""
	}

	// Normalize different possible underlying types into a plain string.
	switch v := value.(type) {
	case string:
		// If the string itself is a JSON encoded string (e.g. "\"2117609272\""), try to unmarshal it.
		var s string
		if err := json.Unmarshal([]byte(v), &s); err == nil {
			return s
		}
		// Fallback to strconv.Unquote for quoted literals
		if unq, err := strconv.Unquote(v); err == nil {
			return unq
		}
		return v
	case []byte:
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			return s
		}
		return string(v)
	default:
		// For numbers, booleans, etc.
		return fmt.Sprint(v)
	}
}

// SetValue 设置配置
func SetValue(app core.App, key string, value string) error {
	record, err := app.FindFirstRecordByData("config", "key", key)

	if err != nil {
		return err
	}

	record.Set("value", value)

	return app.Save(record)
}
