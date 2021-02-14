package write

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

const ProjectId = "madproject-271618"
const BucketName = "madproject-chaloperie"
const FirestoreRootCollection = "chaloperie"

var storageCli *storage.Client
var firestoreCli *firestore.Client

type catCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	PictureURL  string `json:"pictureURL"`
}

type catPersistent struct {
	Title string
	Description string
	Datetime string
	PictureURL string
}

func init() {
	// Removing flags to logger to avoid logging with timestamps and
	// make Google Cloud Logging confused.
	log.SetFlags(0)
}

// CatHTTP is an HTTP Cloud Function that takes get
// cats details and store them to Cloud Firestore.
func CatHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := context.Background()

	err := initClientStorage(ctx)
	if err != nil {
		log.Printf(`{"severity": "error", "message": "could not initialize GCS client", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, "Google encounters problem so your cat could not be created 😿", http.StatusInternalServerError)
		return
	}

	err = initClientFirestore(ctx)
	if err != nil {
		log.Printf(`{"severity": "error", "message": "could not initialize Google Cloud Firestore client", "logging.googleapis.com/trace": "%v"}`+"\n", err)
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
	ref, err := uploadCatPicture(ctx, BucketName, uuid.New().String(), cat.PictureURL)
	if err != nil {
		log.Printf(`{"severity": "error", "message": "saving to Google Cloud Storage failed", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, "Something wrong happened and your cat could not be created 😿", http.StatusInternalServerError)
	}

	// Saving cat to Cloud Firestore
	log.Printf("Saving cat [%s] to Cloud Firestore!\n", ref)
	docRef := firestoreCli.Collection(FirestoreRootCollection).NewDoc()
	_, err = docRef.Create(ctx, catPersistent{
		Title: cat.Title,
		Description: cat.Description,
		Datetime: time.Now().Format(time.RFC3339),
		PictureURL: ref,
	})
	if err != nil {
		log.Printf(`{"severity": "error", "message": "saving to Google Cloud Firestore failed", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, "Something wrong happened and your cat could not be created 😿", http.StatusInternalServerError)
	}
}

func validateInputs(cat catCreateRequest) error {
	if len(cat.Title) < 5 || len(cat.Title) > 30 {
		return errors.New("cat title should have a minimum length of 5 and is capped at 30 characters maximum 😼")
	}
	if len(cat.PictureURL) <= 0 {
		return errors.New("A cat image should be provided 😼")
	}
	return nil
}

func initClientStorage(ctx context.Context) error {
	var err error
	if storageCli == nil {
		storageCli, err = storage.NewClient(ctx)
	}
	return err
}

func initClientFirestore(ctx context.Context) error {
	var err error
	if firestoreCli == nil {
		firestoreCli, err = firestore.NewClient(ctx, ProjectId)
	}
	return err
}

func uploadCatPicture(ctx context.Context, bucket, filename, data string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	imgBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("could not decode base64 data. got: %v", err)
	}
	imgBuf := bytes.NewBuffer(imgBytes)

	bucketObj := storageCli.Bucket(bucket).Object(filename)
	wc := bucketObj.NewWriter(ctx)
	if _, err := io.Copy(wc, imgBuf); err != nil {
		return "", fmt.Errorf("error when copying image buffer to writer. got: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("closing writer resulted in an error. got: %v", err)
	}
	log.Printf(`{"severity": "info", "message": "image [%s] has been uploaded"`, filename)

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, filename), nil
}
