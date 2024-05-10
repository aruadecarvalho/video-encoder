package repositories_test

import (
	"testing"
	"video-encoder/application/repositories"
	"video-encoder/domain"
	"video-encoder/framework/database"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
)

func InsertJobInDb(db *gorm.DB, t *testing.T) (*domain.Job, *domain.Video) {
	video := InsertVideoInDb(db)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	return job, video
}

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	job, video := InsertJobInDb(db, t)

	repoJob := repositories.JobRepositoryDb{Db: db}
	j, err := repoJob.Find(job.ID)

	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestJobRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	job, _ := InsertJobInDb(db, t)

	job.Status = "Completed"

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Update(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.Status, job.Status)
}
