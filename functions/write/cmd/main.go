package main

import (
	"context"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/MadJlzz/chaloperie-api/functions/write"
)

func main() {
	ctx := context.Background()
	if err := funcframework.RegisterHTTPFunctionContext(ctx, "/api/write", write.CatHTTP); err != nil {
		log.Fatalf("funcframework.RegisterHTTPFunctionContext: %v\n", err)
	}
	// Use PORT environment variable, or default to 8080.
	port := "8090"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../credentials/chaloperie-writer.json")
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
