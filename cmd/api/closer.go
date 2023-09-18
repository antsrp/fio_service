package main

import (
	"io"

	"go.uber.org/zap"
)

func handleCloser(logger *zap.Logger, resource string, closer io.Closer) {
	if err := closer.Close(); err != nil {
		logger.Sugar().Errorf("Can't close %q: %s", resource, err)
	}
}
