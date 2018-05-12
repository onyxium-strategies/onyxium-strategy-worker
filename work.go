package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
)

type WorkRequest struct {
	ID       int
	Strategy *models.Strategy
}
