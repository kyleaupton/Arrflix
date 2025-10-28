package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/kyleaupton/snaggle/ops/internal/config"
	"github.com/kyleaupton/snaggle/ops/internal/docker"
	"github.com/kyleaupton/snaggle/ops/internal/reconciler"
)

func main() {
	log.Println("Starting Snaggle Ops Reconciler")

	// Load configuration
	cfg := config.Load()

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	// Create reconciler
	rec := reconciler.New(dockerClient, &cfg)

	// Set up signal handling
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Run reconciler
	if err := rec.Run(ctx); err != nil {
		log.Fatalf("Reconciler failed: %v", err)
	}

	log.Println("Snaggle Ops Reconciler stopped")
}
