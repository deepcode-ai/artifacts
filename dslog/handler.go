package dslog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"

	"golang.org/x/exp/slog"
)

type Option struct {
	Level slog.Leveler

	Writer io.Writer

	Converter Converter

	Humanize bool
}

type DSHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (o Option) NewDSHandler() *DSHandler {
	if o.Level == nil {
		o.Level = LevelInfo
	}
	if o.Writer == nil {
		o.Writer = os.Stderr
	}

	return &DSHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

func (h *DSHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *DSHandler) Handle(ctx context.Context, record slog.Record) error {
	converter := DefaultConverter
	if h.option.Converter != nil {
		converter = h.option.Converter
	}

	message := converter(h.attrs, &record)

	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(b)+1))
	buf.Write(b)
	buf.WriteByte('\n')

	_, err = buf.WriteTo(h.option.Writer)
	return err
}

func (h *DSHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DSHandler{
		option: h.option,
		attrs:  append(h.attrs, attrs...),
	}
}

func (h *DSHandler) WithGroup(name string) slog.Handler {
	return &DSHandler{
		option: h.option,
		groups: append(h.groups, name),
	}
}
