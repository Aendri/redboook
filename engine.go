package app

import (
	"context"
	"fmt"

	"github.com/gopherslab/redbook/internal/config"
	engineModel "github.com/gopherslab/redbook/internal/engines/engine/model"
	"github.com/gopherslab/redbook/pkg/log"
)

type engines []engineModel.Engine

func buildEngines(
	_ log.Logger,
	_ *config.Config,
	_ *services,
) engines {
	eng := make(engines, 0)
	return eng
}

func (e engines) Start(ctx context.Context, log log.Logger) error {
	for _, engine := range e {
		if err := engine.Start(ctx, log); err != nil {
			return fmt.Errorf("error starting engine(%v):%w", engine.GetName(), err)
		}
	}
	return nil
}

func (e engines) Stop(ctx context.Context, log log.Logger) error {
	var lastErr error
	for _, engine := range e {
		if err := engine.Stop(ctx, log); err != nil {
			lastErr = fmt.Errorf("error stopping engine(%v):%w", engine.GetName(), err)
		}
	}
	return lastErr
}
