package write

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatHTTPContentKO(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("This is only plain text."))
	r.Header.Add("Content-Type", "text/plain")

	rr := httptest.NewRecorder()
	CatHTTP(rr, r)

	resp := rr.Result()
	respBody := bufio.NewScanner(resp.Body)
	respBody.Scan()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "Something wrong happend and your cat could not be created 😿", respBody.Text())
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
