package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/pocketbase/pocketbase/tools/router"
)

// rawBodyCtxKey 是用于在 request.Context 中存储原始 body 的私有 key 类型，避免与外部冲突。
type rawBodyCtxKey struct{}

// RawBodyFromContext 返回通过 ParseBody 缓存的原始请求体字节，若不存在则返回 nil。
func RawBodyFromContext(r *http.Request) []byte {
	if r == nil {
		return nil
	}
	v := r.Context().Value(rawBodyCtxKey{})
	if v == nil {
		return nil
	}
	if b, ok := v.([]byte); ok {
		return b
	}
	return nil
}

// ParseBody 读取并解析请求体为泛型类型 T，同时返回原始 body 字节。
// 支持的 Content-Type：
// - application/json
// - application/x-www-form-urlencoded
// - multipart/form-data（仅字段，文件忽略）
// - text/plain（当 T 为 string 或 []byte 时）
// 注意：本函数会消耗 r.Body，但会自动重置，确保后续中间件仍可读取。
func ParseBody[T any](r *http.Request) (out T, raw []byte, err error) {
	if r == nil || r.Body == nil {
		var zero T
		return zero, nil, errors.New("nil request or body")
	}
	// 若上下文已有缓存则直接使用，避免重复读取消耗
	if cached := RawBodyFromContext(r); cached != nil {
		raw = cached
		// 仍旧重置一次，确保调用方再次读取时从头开始；
		// 使用 router.RereadableReadCloser 以实现深层可复读，并兼容上游如 limitedReader 的 Rereader 调用。
		r.Body = &router.RereadableReadCloser{ReadCloser: io.NopCloser(bytes.NewReader(raw))}
		r.ContentLength = int64(len(raw))
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(raw)), nil
		}
	} else {
		// 读取原始 body
		raw, err = io.ReadAll(r.Body)
		if err != nil {
			var zero T
			return zero, raw, err
		}
		// 重置 body 并缓存到 context，允许后续重复读取与访问原始字节
		// 使用 RereadableReadCloser 包裹，保持对 Rereader 的支持，避免破坏如限流等上游包装器的语义。
		r.Body = &router.RereadableReadCloser{ReadCloser: io.NopCloser(bytes.NewReader(raw))}
		r.ContentLength = int64(len(raw))
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(raw)), nil
		}
		// 将原始字节加入 context；使用 WithContext 返回新 *http.Request，再覆盖原对象内容以便调用方不需要修改其引用。
		ctx := context.WithValue(r.Context(), rawBodyCtxKey{}, raw)
		*r = *r.WithContext(ctx)
	}

	// 根据 Content-Type 解析
	ctype := r.Header.Get("Content-Type")
	mediaType, params, _ := mime.ParseMediaType(ctype)

	// JSON 解析
	if mediaType == "application/json" || mediaType == "text/json" || len(raw) == 0 {
		var v T
		if len(raw) == 0 {
			return v, raw, nil
		}
		if err := json.Unmarshal(raw, &v); err != nil {
			return v, raw, err
		}
		return v, raw, nil
	}

	switch mediaType {
	case "application/x-www-form-urlencoded":
		values, err := url.ParseQuery(string(raw))
		if err != nil {
			var zero T
			return zero, raw, err
		}
		// 将 form values 转为扁平 map（取每个 key 的第一个值）
		flat := make(map[string]any, len(values))
		for k, v := range values {
			if len(v) > 0 {
				flat[k] = v[0]
			}
		}
		// 通过 JSON 一次映射到 T
		buf, _ := json.Marshal(flat)
		var vv T
		if err := json.Unmarshal(buf, &vv); err != nil {
			return vv, raw, err
		}
		return vv, raw, nil

	case "multipart/form-data":
		boundary := params["boundary"]
		if boundary == "" {
			var zero T
			return zero, raw, errors.New("multipart boundary not found")
		}
		mr := multipart.NewReader(bytes.NewReader(raw), boundary)
		flat := map[string]any{}
		for {
			part, perr := mr.NextPart()
			if perr == io.EOF {
				break
			}
			if perr != nil {
				var zero T
				return zero, raw, perr
			}
			// 忽略文件，仅收集表单字段
			if part.FileName() != "" {
				// drain but ignore content to move reader ahead
				_, _ = io.Copy(io.Discard, part)
				_ = part.Close()
				continue
			}
			name := part.FormName()
			if name == "" {
				_, _ = io.Copy(io.Discard, part)
				_ = part.Close()
				continue
			}
			data, _ := io.ReadAll(part)
			_ = part.Close()
			flat[name] = string(data)
		}
		buf, _ := json.Marshal(flat)
		var vv T
		if err := json.Unmarshal(buf, &vv); err != nil {
			return vv, raw, err
		}
		return vv, raw, nil

	case "text/plain":
		var t T
		switch ptr := any(&t).(type) {
		case *string:
			*ptr = string(raw)
			return t, raw, nil
		case *[]byte:
			*ptr = append((*ptr)[:0], raw...)
			return t, raw, nil
		default:
			// 尝试 JSON 反序列化作为兜底
			if err := json.Unmarshal(raw, &t); err != nil {
				return t, raw, err
			}
			return t, raw, nil
		}
	default:
		// 未知类型：尝试 JSON 解析；若 T 为 string/[]byte 则直接返回
		var t T
		switch ptr := any(&t).(type) {
		case *string:
			*ptr = string(raw)
			return t, raw, nil
		case *[]byte:
			*ptr = append((*ptr)[:0], raw...)
			return t, raw, nil
		default:
			if len(raw) == 0 {
				return t, raw, nil
			}
			if err := json.Unmarshal(raw, &t); err != nil {
				return t, raw, err
			}
			return t, raw, nil
		}
	}
}
