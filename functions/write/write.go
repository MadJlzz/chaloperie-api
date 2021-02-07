package write

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

var client *storage.Client

type catCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	PictureURL  string `json:"pictureUrl"`
}

func init() {
	// Removing flags to logger to avoid logging with timestamps and
	// make Google Cloud Logging confused.
	log.SetFlags(0)
}

// CatHTTP is an HTTP Cloud Function that takes get
// cats details and store them to Cloud Firestore.
func CatHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	err := initClientStorage(ctx)
	if err != nil {
		log.Printf(`{"severity": "error", "message": "could not initialize GCS client", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, "Google encounters problem so your cat could not be created 😿", http.StatusInternalServerError)
		return
	}

	var cat catCreateRequest
	if err = json.NewDecoder(r.Body).Decode(&cat); err != nil {
		log.Printf(`{"severity": "error", "message": "could not decode request body as json", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, "Something wrong happened and your cat could not be created 😿", http.StatusInternalServerError)
		return
	}

	// Validate the inputs to be sure that we have everything.
	// If the validation fails, we return an error to the frontend.
	err = validateInputs(cat)
	if err != nil {
		log.Printf(`{"severity": "error", "message": "validation of cat creation request failed", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Saving cat picture to Cloud Storage
	uploadCatPicture(ctx, "madproject-chaloperie", uuid.New().String())

	// Saving cat to Cloud Firestore
	log.Println("Saving cat to Cloud Firestore!")
}

func validateInputs(cat catCreateRequest) error {
	if len(cat.Title) < 5 || len(cat.Title) > 30 {
		return errors.New("Cat title should have a minimum length of 5 and is capped at 30 characters maximum 😼")
	}
	if len(cat.PictureURL) <= 0 {
		return errors.New("A cat image should be provided 😼")
	}
	return nil
}

func initClientStorage(ctx context.Context) error {
	var err error
	if client == nil {
		client, err = storage.NewClient(ctx)
	}
	return err
}

func uploadCatPicture(ctx context.Context, bucket, filename string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	log.Println(bucket)
	log.Println(filename)

	return nil
}
