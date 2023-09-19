package config

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/antsrp/fio_service/internal/infrastructure/broker"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	path, _ := os.Getwd()
	if base := filepath.Base(path); base == "config" {
		dir := filepath.Join(path, "../..")
		os.Chdir(dir)
	}
}

func TestWithLoad(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("can't load env variables: %v", err.Error())
	}

	prefix := "BROKER"
	expected := broker.Settings{
		Type:        "kafka",
		Host:        "localhost",
		Port:        9092,
		GroupID:     "fio_group",
		Topic:       "FIO",
		TopicFailed: "FIO_FAILED",
	}

	got, err := Parse[broker.Settings](prefix)
	require.NoError(t, err, "there should be no error when executing this request")
	assert.Equalf(t, expected, *got, "parse env variables to config with prefix %s: expected %v, got %v", prefix, expected, *got)
}

func TestNonExistingPrefix(t *testing.T) {
	prefix := "PREFix"
	expected := broker.Settings{}

	got, err := Parse[broker.Settings](prefix)
	require.NoError(t, err, "there should be no error when executing this request")
	assert.Equalf(t, expected, *got, "parse env variables to config with prefix %s: expected %v, got %v", prefix, expected, *got)
}
