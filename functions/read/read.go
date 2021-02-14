package read

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

const ProjectId = "madproject-271618"
const FirestoreRootCollection = "chaloperie"

var firestoreCli *firestore.Client

type CatReadResponse struct {
	Cats []CatPersistent `json:"cats"`
}

type CatPersistent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Datetime    string `json:"datetime"`
	PictureURL  string `json:"pictureURL"`
}

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
	w.Header().Set("Content-Type", "application/json")

	ctx := context.Background()

	err := initClientFirestore(ctx)
	if err != nil {
		log.Printf(`{"severity": "error", "message": "could not initialize Google Cloud Firestore client", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, "Google encounters problem so no cats could be returned 😿", http.StatusInternalServerError)
		return
	}

	log.Println("Reading all cats from Google Cloud Firestore")
	docRefs := firestoreCli.Collection(FirestoreRootCollection).
		OrderBy("Datetime", firestore.Desc).
		Limit(5).
		Documents(ctx)
	docs, err := docRefs.GetAll()
	if err != nil {
		log.Printf(`{"severity": "error", "message": "read cats from Google Cloud Firestore failed", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, "Google encounters problem so no cats could be returned 😿", http.StatusInternalServerError)
		return
	}

	resp := CatReadResponse{}
	for _, ds := range docs {
		data := ds.Data()
		resp.Cats = append(resp.Cats, CatPersistent{
			Title:       data["Title"].(string),
			Description: data["Description"].(string),
			Datetime:    data["Datetime"].(string),
			PictureURL:  data["PictureURL"].(string),
		})
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf(`{"severity": "error", "message": "encoding cat response into json failed", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, "Google encounters problem so no cats could be returned 😿", http.StatusInternalServerError)
		return
	}
}

func initClientFirestore(ctx context.Context) error {
	var err error
	if firestoreCli == nil {
		firestoreCli, err = firestore.NewClient(ctx, ProjectId)
	}
	return err
}
