package mapper

import (
	"testing"

	"github.com/antsrp/fio_service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	person = domain.Person{
		Id: 1,
		PersonCommon: domain.PersonCommon{
			Name:     "some name",
			Surname:  "some surname",
			Patronym: "some patronym",
		},
		Age:         25,
		Gender:      "male",
		Nationality: "some nationality",
	}
	personJsonized = []byte(`{
"id": 1,
"name": "some name",
"surname": "some surname",
"patronym": "some patronym",
"age": 25,
"gender": "male",
"nationality": "some nationality"
}`)
	personJsonizedWithIndent = []byte(`{
	"id": 1,
	"name": "some name",
	"surname": "some surname",
	"patronym": "some patronym",
	"age": 25,
	"gender": "male",
	"nationality": "some nationality"
}`)
)

func TestPersonToJSON(t *testing.T) {
	expected := personJsonized

	got, err := ToJSON[domain.Person](person, &Indent{})

	require.NoError(t, err, "there should be no error when executing this request")
	assert.Equal(t, expected, got, "person %v, expected value %v, got %v", person, expected, got)
}

func TestPersonToJSONIndent(t *testing.T) {

	expected := personJsonizedWithIndent

	got, err := ToJSON[domain.Person](person, NewIndent("", "\t"))

	require.NoError(t, err, "there should be no error when executing this request")
	assert.Equal(t, expected, got, "person %v, expected value %v, got %v", person, expected, got)
}

func TestPersonFromJSON(t *testing.T) {

	data := personJsonized
	expected := person

	got, err := FromJSON[domain.Person](data)

	require.NoError(t, err, "there should be no error when executing this request")
	assert.Equal(t, expected, *got, "person %v, expected value %v, got %v", person, expected, *got)
}

func TestPersonFromJSONIndent(t *testing.T) {

	data := personJsonizedWithIndent
	expected := person

	got, err := FromJSON[domain.Person](data)

	require.NoError(t, err, "there should be no error when executing this request")
	assert.Equal(t, expected, *got, "person %v, expected value %v, got %v", person, expected, *got)
}
