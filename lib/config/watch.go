package config

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"

	"go-micro.dev/v4/config"
)

type Watcher struct {
	log *slog.Logger
	w   config.Watcher
}

func NewWatcher(log *slog.Logger, path ...string) *Watcher {
	log = log.With("file", cfg.file, "path", strings.Join(path, "."))

	if _, file, line, ok := runtime.Caller(1); ok {
		caller := fmt.Sprintf("%s/%s:%d", filepath.Base(filepath.Dir(file)), filepath.Base(file), line)
		log = log.With("caller", caller)
	}

	w, err := cfg.Watch(path...)
	if err != nil {
		panic(err)
	}

	return &Watcher{
		log: log,
		w:   w,
	}
}

func (cw *Watcher) Watch(ctx context.Context, callback func(val Value)) error {
	go func() {
		<-ctx.Done()
		cw.w.Stop()
		cw.log.Debug("stopped config watcher")
	}()

	go func() {
		cw.log.Debug("started config watcher")

		for {
			v, err := cw.w.Next()
			if err != nil {
				return
			}

			if string(v.Bytes()) != "null" {
				cw.log.Info("watched config updated", "value", v)
				callback(v)
			}
		}
	}()

	return nil
}
