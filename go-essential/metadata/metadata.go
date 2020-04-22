// Package metadata is a way of defining message headers
package metadata

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type metadataKey struct{}

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]interface{}

func (md Metadata) Len() int {
	return len(md)
}

func (md Metadata) Get(key string) (interface{}, bool) {
	// attempt to get as is
	val, ok := md[key]
	if ok {
		return val, ok
	}

	// attempt to get lower case
	val, ok = md[strings.Title(key)]
	return val, ok
}

func (md Metadata) Set(key, val string) {
	md[key] = val
}

func (md Metadata) Delete(key string) {
	// delete key as-is
	delete(md, key)
	// delete also Title key
	delete(md, strings.Title(key))
}

// New creates an MD from a given key-value map.
func New(m map[string]interface{}) Metadata {
	md := Metadata{}
	for k, val := range m {
		md[k] = val
	}
	return md
}

// Join joins any number of mds into a single MD.
// The order of values for each key is determined by the order in which
// the mds containing those values are presented to Join.
func Join(mds ...Metadata) Metadata {
	out := Metadata{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = v
		}
	}
	return out
}

// Copy makes a copy of the metadata
func Copy(md Metadata) Metadata {
	cmd := make(Metadata, len(md))
	for k, v := range md {
		cmd[k] = v
	}
	return cmd
}

// Pairs returns an MD formed by the mapping of key, value ...
// Pairs panics if len(kv) is odd.
func Pairs(kv ...interface{}) Metadata {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := Metadata{}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = s.(string)
			continue
		}
		md[key] = s
	}
	return md
}

// Delete key from metadata
func Delete(ctx context.Context, k string) context.Context {
	return Set(ctx, k, "")
}

// Set add key with val to metadata
func Set(ctx context.Context, k string, v interface{}) context.Context {
	md, ok := FromContext(ctx)
	if !ok {
		md = make(Metadata)
	}
	if v == "" {
		delete(md, k)
	} else {
		md[k] = v
	}
	return context.WithValue(ctx, metadataKey{}, md)
}

// Get returns a single value from metadata in the context
func Get(ctx context.Context, key string) (interface{}, bool) {
	md, ok := FromContext(ctx)
	if !ok {
		return "", ok
	}
	// attempt to get as is
	val, ok := md[key]
	if ok {
		return val, ok
	}

	// attempt to get lower case
	val, ok = md[strings.Title(key)]

	return val, ok
}

// FromContext returns metadata from the given context
func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return nil, ok
	}

	// capitalise all values
	newMD := make(Metadata, len(md))
	for k, v := range md {
		newMD[strings.Title(k)] = v
	}

	return newMD, ok
}

// NewContext creates a new context with the given metadata
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataKey{}, md)
}

// MergeContext merges metadata to existing metadata, overwriting if specified
func MergeContext(ctx context.Context, patchMd Metadata, overwrite bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	md, _ := ctx.Value(metadataKey{}).(Metadata)
	cmd := make(Metadata, len(md))
	for k, v := range md {
		cmd[k] = v
	}
	for k, v := range patchMd {
		if _, ok := cmd[k]; ok && !overwrite {
			// skip
		} else if v != "" {
			cmd[k] = v
		} else {
			delete(cmd, k)
		}
	}
	return context.WithValue(ctx, metadataKey{}, cmd)
}

// Value get value from metadata in context return nil if not found
func Value(ctx context.Context, key string) interface{} {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return nil
	}
	return md[key]
}

// Bool get boolean from metadata in context use strconv.Parse.
func Bool(ctx context.Context, key string) bool {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return false
	}

	switch md[key].(type) {
	case bool:
		return md[key].(bool)
	case string:
		ok, _ = strconv.ParseBool(md[key].(string))
		return ok
	default:
		return false
	}
}

// String get string value from metadata in context
func String(ctx context.Context, key string) string {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return ""
	}
	str, _ := md[key].(string)
	return str
}

// Int32 get int64 value from metadata in context
func Int32(ctx context.Context, key string) int32 {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return 0
	}
	i32, _ := md[key].(int32)
	return i32
}

// Int64 get int64 value from metadata in context
func Int64(ctx context.Context, key string) int64 {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return 0
	}
	i64, _ := md[key].(int64)
	return i64
}
