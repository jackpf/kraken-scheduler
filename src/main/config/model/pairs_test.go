package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestJsonPair struct {
	Pair Pair `json:"pair"`
}

func TestPair_Name(t *testing.T) {
	pair := Pairs["ADAETH"]
	assert.Equal(t, "ADAETH", pair.Name())
}

func TestPair_UnmarshalJSON(t *testing.T) {
	var pair TestJsonPair

	err := json.Unmarshal([]byte(`{"pair": "ADAETH"}`), &pair)

	assert.NoError(t, err)
	assert.Equal(t, Pairs["ADAETH"], pair.Pair)
}

func TestPair_UnmarshalJSON_Invalid(t *testing.T) {
	var pair TestJsonPair

	err := json.Unmarshal([]byte(`{"pair": "foobar"}`), &pair)

	assert.Errorf(t, err, "pair foobar is not valid")
}
