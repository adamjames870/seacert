package main

import (
	"context"
	"log"
	"time"

	"github.com/adamjames870/seacert/internal/notifications"
	_ "github.com/lib/pq"
)

func main() {

	var state notificationState
	err := run(&state)
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}

}

func run(state *notificationState) error {

	errState := LoadState(state)
	if errState != nil {
		return errState
	}

	ctx := context.Background()

	generator := notifications.NewGenerator(state.Repo)

	count7Day, err7Day := generator.GenerateNoCertificates7Day(ctx)
	logGenerationResult("NoCertificates7Day", count7Day, err7Day)

	count1Month, err1Month := generator.GenerateNoCertificates1Month(ctx)
	logGenerationResult("NoCertificates1Month", count1Month, err1Month)

	countExp, errExp := generator.GenerateCertificateExpirySummary(ctx, time.Now())
	logGenerationResult("CertificateExpirySummary", countExp, errExp)

	processor := notifications.NewProcessor(state.Repo)

	results, errProcess := processor.ProcessPendingNotifications(ctx, 1000)
	if errProcess != nil {
		log.Printf("Error processing notifications: %v", errProcess)
	} else {
		log.Printf(
			"Notifications found: %d | Completed: %d | Failed: %d",
			results.Found,
			results.Completed,
			results.Failed,
		)
	}

	return nil
}

func logGenerationResult(name string, count int, err error) {
	if err != nil {
		log.Printf("%s: generated %d notification(s), error: %v", name, count, err)
		return
	}

	log.Printf("%s: generated %d notification(s)", name, count)
}
