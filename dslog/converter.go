package dslog

import (
	"encoding"
	"fmt"
	"reflect"

	"golang.org/x/exp/slog"
)

type Converter func(loggerAttr []slog.Attr, record *slog.Record) map[string]interface{}

func DefaultConverter(loggerAttr []slog.Attr, record *slog.Record) map[string]interface{} {
	log := make(map[string]interface{}, len(loggerAttr)+3)
	log["timestamp"] = record.Time.Unix()
	log["level"] = record.Level.String()
	log["message"] = record.Message

	recordValues := make(map[string]interface{}, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		k, v := attrToValue(attr)
		recordValues[k] = v
		return true
	})

	loggerValues := attrsToValue(loggerAttr)
	if v, ok := loggerValues["context"]; ok {
		val := v.(map[string]interface{})
		for k, v := range recordValues {
			val[k] = v
		}
		log["context"] = val
	} else {
		log["context"] = recordValues
	}

	for k, v := range loggerValues {
		if k == "context" {
			continue
		}
		log[k] = v
	}

	return log
}

func attrsToValue(attrs []slog.Attr) map[string]interface{} {
	log := map[string]interface{}{}

	for i := range attrs {
		k, v := attrToValue(attrs[i])
		log[k] = v
	}

	return log
}

func attrToValue(attr slog.Attr) (string, interface{}) {
	k := attr.Key
	v := attr.Value
	kind := v.Kind()

	switch kind {
	case slog.KindAny:
		if k == "error" {
			if err, ok := v.Any().(error); ok {
				return k, buildExceptions(err)
			}
		}

		return k, v.Any()
	case slog.KindLogValuer:
		return k, v.Any()
	case slog.KindGroup:
		return k, attrsToValue(v.Group())
	case slog.KindInt64:
		return k, v.Int64()
	case slog.KindUint64:
		return k, v.Uint64()
	case slog.KindFloat64:
		return k, v.Float64()
	case slog.KindString:
		return k, v.String()
	case slog.KindBool:
		return k, v.Bool()
	case slog.KindDuration:
		return k, v.Duration()
	case slog.KindTime:
		return k, v.Time()
	default:
		return k, anyValueToString(v)
	}
}

func anyValueToString(v slog.Value) string {
	if tm, ok := v.Any().(encoding.TextMarshaler); ok {
		data, err := tm.MarshalText()
		if err != nil {
			return ""
		}

		return string(data)
	}

	return fmt.Sprintf("%+v", v.Any())
}

func buildExceptions(err error) map[string]interface{} {
	return map[string]interface{}{
		"kind":  reflect.TypeOf(err).String(),
		"error": err.Error(),
		"stack": nil, // @TODO
	}
}
