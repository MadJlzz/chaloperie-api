package write

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
