package domain_test

import (
	"testing"
	"time"
	"video-encoder/domain"

	"github.com/stretchr/testify/require"

	uuid "github.com/satori/go.uuid"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)

}

func TestVideoIValidation(t *testing.T) {
	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.ResourceID = "resourceId"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Nil(t, err)
}
