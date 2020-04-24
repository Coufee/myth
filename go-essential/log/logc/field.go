package logc

import (
	"math"
	"time"

	"myth/go-essential/log/logc/internal/core"
)

// Field represents a map of entry level data used for structured logging.
// type Field map[string]interface{}
type Field = core.Field

// KVString construct Field with string value.
func KVString(key string, value string) Field {
	return Field{Key: key, Type: core.StringType, StringVal: value}
}

// KVInt construct Field with int value.
func KVInt(key string, value int) Field {
	return Field{Key: key, Type: core.IntTpye, Int64Val: int64(value)}
}

// KVInt64 construct Field with int64 value.
func KVInt64(key string, value int64) Field {
	return Field{Key: key, Type: core.Int64Type, Int64Val: value}
}

// KVUint construct Field with uint value.
func KVUint(key string, value uint) Field {
	return Field{Key: key, Type: core.UintType, Int64Val: int64(value)}
}

// KVUint64 construct Field with uint64 value.
func KVUint64(key string, value uint64) Field {
	return Field{Key: key, Type: core.Uint64Type, Int64Val: int64(value)}
}

// KVFloat32 construct Field with float32 value.
func KVFloat32(key string, value float32) Field {
	return Field{Key: key, Type: core.Float32Type, Int64Val: int64(math.Float32bits(value))}
}

// KVFloat64 construct Field with float64 value.
func KVFloat64(key string, value float64) Field {
	return Field{Key: key, Type: core.Float64Type, Int64Val: int64(math.Float64bits(value))}
}

// KVDuration construct Field with Duration value.
func KVDuration(key string, value time.Duration) Field {
	return Field{Key: key, Type: core.DurationType, Int64Val: int64(value)}
}

// KV return a log kv for logging field.
// NOTE: use KV{type name} can avoid object alloc and get better performance. []~(￣▽￣)~*干杯
func KV(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
