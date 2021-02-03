package main

import (
	"context"
	"time"
)

// AnalyzeParams is analyze parameters for analysing
type AnalyzeParams struct {
	interval time.Duration
}

func reAnalyze(params AnalyzeParams) {
	/* do something */
}

func periodicAnalyze(ctx context.Context, params AnalyzeParams) {
	ticker := time.NewTicker(params.interval)
	defer ticker.Stop()
	reAnalyze(params)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			reAnalyze(params)
		}
	}
}
