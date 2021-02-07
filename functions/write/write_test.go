package write

import (
	"bufio"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

// TODO: Think about how it could be possible (or not)
// to mock the Google Cloud Storage client.
// After some research, it seems very difficult.

// Setup tests to disable logger.
func TestMain(m *testing.M) {
	// Overriding global client with a testing one.
	client, _ = storage.NewClient(context.Background(), option.WithoutAuthentication())

	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestCatHTTPContentKO(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("This is only plain text."))
	r.Header.Add("Content-Type", "text/plain")

	rr := httptest.NewRecorder()
	CatHTTP(rr, r)

	resp := rr.Result()
	respBody := bufio.NewScanner(resp.Body)
	respBody.Scan()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "Something wrong happened and your cat could not be created 😿", respBody.Text())
}

func TestCatHTTPContentOK(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"title": "Shadow does HAI", "pictureUrl": "https://google.com"}`))
	r.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	CatHTTP(rr, r)

	resp := rr.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestValidateInputsTitleToSmall(t *testing.T) {
	c := catCreateRequest{}
	if err := validateInputs(c); err == nil {
		t.Error("empty cat title should have thrown an error")
	} else {
		assert.Equal(t, err.Error(), "Cat title should have a minimum length of 5 and is capped at 30 characters maximum 😼")
	}
}

func TestValidateInputsTitleToBig(t *testing.T) {
	c := catCreateRequest{
		Title: "Shadow does HAI",
	}
	if err := validateInputs(c); err == nil {
		t.Error("empty pictureURL should have thrown an error")
	} else {
		assert.Equal(t, err.Error(), "A cat image should be provided 😼")
	}
}

func TestValidateInputsOk(t *testing.T) {
	c := catCreateRequest{
		Title:      "Shadow does HAI",
		PictureURL: "https://my.image.com",
	}
	if err := validateInputs(c); err != nil {
		t.Error("validateInputs should return no error")
	}
}
