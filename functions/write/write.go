package write

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type catCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	PictureURL  string `json:"pictureUrl"`
}

// CatHTTP is an HTTP Cloud Function that takes get
// cats details and store them to Cloud Firestore.
func CatHTTP(w http.ResponseWriter, r *http.Request) {
	var cat catCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		fmt.Printf(`{"severity": "error", "message": "could not decode request body as json", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, "Something wrong happend and your cat could not be created 😿", http.StatusInternalServerError)
		return
	}

	// Validate the inputs to be sure that we have everything.
	// If the validation fails, we return an error to the frontend.
	err := validateInputs(cat)
	if err != nil {
		fmt.Printf(`{"severity": "error", "message": "validation of cat creation request failed", "logging.googleapis.com/trace": "%v"}`+"\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Saving cat to Cloud Firestore
	fmt.Println("Saving cat to Cloud Firestore!")
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
