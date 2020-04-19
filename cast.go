package req

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func (r *Resp) isPathShadowedInDeepMap(path []string, m map[string]interface{}) string {
	var parentVal interface{}
	for i := 1; i < len(path); i++ {
		parentVal = r.searchMap(m, path[0:i])
		if parentVal == nil {
			// not found, no need to add more path elements
			return ""
		}
		switch parentVal.(type) {
		case map[interface{}]interface{}:
			continue
		case map[string]interface{}:
			continue
		default:
			// parentVal is a regular value which shadows "path"
			return strings.Join(path[0:i], ".")
		}
	}
	return ""
}

func (r *Resp) searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next
		}

		switch next.(type) {
		case map[interface{}]interface{}:
			return r.searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			return r.searchMap(next.(map[string]interface{}), path[1:])
		default:
			return nil
		}
	}

	return nil
}

func (r *Resp) find(lcaseKey string) interface{} {
	var (
		val    interface{}
		path   = strings.Split(lcaseKey, ".")
		nested = len(path) > 1
	)

	val = r.searchMap(r.forMap(), path)

	if val != nil {
		return val
	}

	if nested && r.isPathShadowedInDeepMap(path, r.forMap()) != "" {
		return nil
	}

	return nil
}

func (r *Resp) forMap() map[string]interface{} {
	if r.respMap != nil {
		return r.respMap
	}

	json.Unmarshal(r.Bytes(), &r.respMap)

	return r.respMap
}

func (r *Resp) get(key string) interface{} {
	lcaseKey := strings.ToLower(key)

	val := r.find(lcaseKey)

	return val
}

func (r *Resp) GetString(key string) string {
	return cast.ToString(r.get(key))
}

func (r *Resp) GetBool(key string) bool {
	return cast.ToBool(r.get(key))
}

func (r *Resp) GetInt(key string) int {
	return cast.ToInt(r.get(key))
}

func (r *Resp) GetInt64(key string) int64 {
	return cast.ToInt64(r.get(key))
}

func (r *Resp) GetFloat64(key string) float64 {
	return cast.ToFloat64(r.get(key))
}

func (r *Resp) GetTime(key string) time.Time {
	return cast.ToTime(r.get(key))
}

func (r *Resp) GetDuration(key string) time.Duration {
	return cast.ToDuration(r.get(key))
}

func (r *Resp) GetStringSlice(key string) []string {
	return cast.ToStringSlice(r.get(key))
}

func (r *Resp) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(r.get(key))
}

func (r *Resp) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(r.get(key))
}

func (r *Resp) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(r.get(key))
}
