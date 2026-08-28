package contentexport

import (
	"encoding/json"
	"strings"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
)

// extInfoRefs 收集 extInfo JSON 中所有字符串值里的图片引用（用于登记打包）。
func extInfoRefs(raw *contract.JSONRaw) []string {
	if raw == nil || len(*raw) == 0 {
		return nil
	}
	var value any
	if err := json.Unmarshal([]byte(*raw), &value); err != nil {
		return nil
	}
	var refs []string
	collectURLStrings(value, &refs)
	return refs
}

// rewriteExtInfo 递归改写 extInfo 中字符串值里的图片引用（只改值、不改键）。
// 解析失败或空值时原样返回。
func rewriteExtInfo(raw *contract.JSONRaw, r *Resolver, depth int) *contract.JSONRaw {
	if raw == nil || len(*raw) == 0 {
		return raw
	}
	var value any
	if err := json.Unmarshal([]byte(*raw), &value); err != nil {
		return raw
	}
	rewritten := walkRewrite(value, r, depth)
	out, err := json.Marshal(rewritten)
	if err != nil {
		return raw
	}
	result := contract.JSONRaw(out)
	return &result
}

func collectURLStrings(value any, out *[]string) {
	switch v := value.(type) {
	case string:
		if isImageURLLike(v) {
			*out = append(*out, v)
		}
	case map[string]any:
		for _, item := range v {
			collectURLStrings(item, out)
		}
	case []any:
		for _, item := range v {
			collectURLStrings(item, out)
		}
	}
}

func walkRewrite(value any, r *Resolver, depth int) any {
	switch v := value.(type) {
	case string:
		if isImageURLLike(v) {
			return rewriteValue(v, r, depth)
		}
		return v
	case map[string]any:
		for key, item := range v {
			v[key] = walkRewrite(item, r, depth)
		}
		return v
	case []any:
		for i, item := range v {
			v[i] = walkRewrite(item, r, depth)
		}
		return v
	default:
		return v
	}
}

func isImageURLLike(value string) bool {
	return strings.HasPrefix(value, "http://") ||
		strings.HasPrefix(value, "https://") ||
		strings.HasPrefix(value, "/uploads/")
}
