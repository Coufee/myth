package logc

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"myth/go-essential/log/logc/internal/filewriter"
)

type FileHandler struct {
	render Render
	fws    [_FileTotalIdx]*filewriter.FileWriter
}

func NewFile(dir string, bufferSize, rotateSize int64, maxLogFile int) *FileHandler {
	// new info writer
	newWriter := func(name string) *filewriter.FileWriter {
		var options []filewriter.Option
		if rotateSize > 0 {
			options = append(options, filewriter.MaxSize(rotateSize))
		}
		if maxLogFile > 0 {
			options = append(options, filewriter.MaxFile(maxLogFile))
		}
		w, err := filewriter.New(filepath.Join(dir, name), options...)
		if err != nil {
			panic(err)
		}
		return w
	}

	handler := &FileHandler{
		render: newPatternRender(defaultFilePattern, handlerTypeFileLog),
	}

	for idx, name := range _FileNames {
		handler.fws[idx] = newWriter(name)
	}

	return handler
}

func (h *FileHandler) Log(ctx context.Context, lv LogLevel, args ...Field) {
	d := toMap(args...)
	// add extra fields
	addExtraField(ctx, d)
	d[_time] = time.Now().Format(_timeFormat)
	var w io.Writer
	switch lv {
	case _LogLevelWarn:
		w = h.fws[_FileWarnIdx]
	case _LogLevelError:
		w = h.fws[_FileErrorIdx]
	default:
		w = h.fws[_FileInfoIdx]
	}

	h.render.Render(w, lv, d)
	w.Write([]byte("\n"))
}

func (h *FileHandler) Close() error {
	for _, fw := range h.fws {
		// ignored error
		fw.Close()
	}
	return nil
}

func (h *FileHandler) SetFormat(format string) {
	h.render = newPatternRender(format, handlerTypeFileLog)
}
