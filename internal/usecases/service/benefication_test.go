package service

import (
	"log"
	"testing"

	"github.com/antsrp/fio_service/internal/domain/benefication"
	"github.com/antsrp/fio_service/internal/infrastructure/http/web"
	iweb "github.com/antsrp/fio_service/internal/interfaces/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var webConnection iweb.Connector

func TestMain(m *testing.M) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Can't create zap logger: ", err)
	}

	webConnection = web.CreateNewConnection(logger)

	m.Run()
}

func Test1(t *testing.T) {
	name := "Alexey"
	data, err := beneficate(webConnection, name)
	expected := benefication.DataCacheValue{
		Age:    46,
		Gender: "male",
		Country: []benefication.Country{
			{ID: "RU", Probability: 0.445},
			{ID: "UA", Probability: 0.213},
			{ID: "BY", Probability: 0.106},
			{ID: "KZ", Probability: 0.055},
			{ID: "IL", Probability: 0.038},
		},
	}

	require.NoError(t, err, "there should be no error when executing this request")
	assert.Equal(t, expected, data, "name %s, expected value %v, got %v", name, expected, data)
}

func Test2(t *testing.T) {
	name := "Dmitriy"
	data, err := beneficate(webConnection, name)
	expected := benefication.DataCacheValue{
		Age:    42,
		Gender: "male",
		Country: []benefication.Country{
			{ID: "UA", Probability: 0.419},
			{ID: "RU", Probability: 0.291},
			{ID: "KZ", Probability: 0.097},
			{ID: "BY", Probability: 0.069},
			{ID: "IL", Probability: 0.019},
		},
	}

	require.NoError(t, err, "there should be no error when executing this request")
	assert.Equal(t, expected, data, "name %s, expected value %v, got %v", name, expected, data)
}

func Test3(t *testing.T) {
	name := "SomeNonExistingName"
	data, err := beneficate(webConnection, name)
	expected := benefication.DataCacheValue{
		Country: []benefication.Country{},
	}

	require.NoError(t, err, "there should be no error when executing this request")
	assert.Equal(t, expected, data, "name %s, expected value %v, got %v", name, expected, data)
}
