package utils_test

import (
	"testing"
	"video-encoder/framework/utils"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `{
		"id": "random-uuid-v4",
		"file_path": "teste.mp4",
		"status": "pending"
		}`
	err := utils.IsJson(json)

	require.Nil(t, err)

	json = `test`
	err = utils.IsJson(json)

	require.Error(t, err)

}
