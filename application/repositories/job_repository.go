package repositories

import (
	"fmt"
	"video-encoder/domain"

	"github.com/jinzhu/gorm"
)

type JobRepository interface {
	Insert(*domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(id *domain.Job) (*domain.Job, error)
}

type JobRepositoryDb struct {
	Db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepositoryDb {
	return &JobRepositoryDb{Db: db}
}

func (repo JobRepositoryDb) Insert(video *domain.Job) (*domain.Job, error) {
	err := repo.Db.Create(video).Error

	if err != nil {
		return nil, err
	}

	return video, nil
}

func (repo JobRepositoryDb) Find(id string) (*domain.Job, error) {
	var job domain.Job
	err := repo.Db.Preload("Video").First(&job, "id = ?", id).Error

	if err != nil || job.ID == "" {
		return nil, fmt.Errorf("job does not exist")
	}

	return &job, nil
}

func (repo JobRepositoryDb) Update(job *domain.Job) (*domain.Job, error) {
	err := repo.Db.Update(&job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}
