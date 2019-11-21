package cmd

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// initializes the basic configuration for the log wrapper
func initLog() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize zap logger")
	}
	return logger, nil
}
