package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"sort"
	"strings"
)

type UTCFormatter struct {
	logrus.Formatter
}

func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}

type DetailFormatter struct {
	TimestampFormat string
}

func NewDetailFormatter(timestampFormat string) *DetailFormatter {
	if len(timestampFormat) == 0 {
		timestampFormat = "Jan _2 15:04:05.000"
	}
	return &DetailFormatter{TimestampFormat: timestampFormat}
}

func CallerPrettyfier(f *runtime.Frame) string {
	s := strings.LastIndex(f.Function, ".") + 1
	return fmt.Sprintf("caller={%s:%d::%s}", path.Base(f.File), f.Line, f.Function[s:])
	//return fmt.Sprintf("[%s:%d::%s]", path.Base(f.File), f.Line, f.Function[s:])
}

func CallerPrettyfierEx(f *runtime.Frame) (string, string) {
	s := strings.LastIndex(f.Function, ".") + 1
	return f.Function[s:], fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
}

func (f *DetailFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(entry.Time.UTC().Format(f.TimestampFormat))
	b.WriteString(" [")
	b.WriteString(entry.Level.String())
	b.WriteString("]")

	if entry.Message != "" {
		b.WriteString(" - ")
		b.WriteString(entry.Message)
	}

	b.WriteString(" || ")

	// Sort keys of entry.Data
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteByte('{')
		_, _ = fmt.Fprint(b, entry.Data[key])
		b.WriteString("}, ")
	}

	if entry.Caller != nil {
		b.WriteString(" ")
		b.WriteString(CallerPrettyfier(entry.Caller))
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}
