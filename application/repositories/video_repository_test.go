package repositories_test

import (
	"testing"
	"time"
	"video-encoder/application/repositories"
	"video-encoder/domain"
	"video-encoder/framework/database"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func InsertVideoInDb(db *gorm.DB) *domain.Video {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	return video
}

func TestVideoRepositoryDbInser(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := InsertVideoInDb(db)

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	v, err := repo.Find(video.ID)

	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, v.ID, video.ID)

}
