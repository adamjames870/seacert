package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adamjames870/seacert/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := storage.Config{
		BucketName:      os.Getenv("R2_BUCKET_NAME"),
		Endpoint:        os.Getenv("R2_ENDPOINT"),
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
	}

	if cfg.BucketName == "" || cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		log.Fatal("R2 storage configuration is incomplete")
	}

	r2, err := storage.NewR2Storage(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to initialize R2 storage: %v", err)
	}

	ctx := context.Background()
	key := "test-delete-" + time.Now().Format("20060102150405") + ".txt"

	fmt.Printf("Generating upload URL for: %s\n", key)
	uploadURL, err := r2.GetPresignedUploadURL(ctx, key, "text/plain", 5*time.Minute)
	if err != nil {
		log.Fatalf("failed to get upload URL: %v", err)
	}
	fmt.Printf("Upload URL: %s\n", uploadURL)
	fmt.Println("Please upload a file using this URL (e.g. via curl -X PUT --data 'hello' URL)")

	fmt.Println("Press Enter after you've uploaded the file...")
	var input string
	fmt.Scanln(&input)

	fmt.Printf("Deleting object: %s\n", key)
	err = r2.DeleteObject(ctx, key)
	if err != nil {
		log.Fatalf("failed to delete object: %v", err)
	}
	fmt.Println("DeleteObject call returned success!")
}
