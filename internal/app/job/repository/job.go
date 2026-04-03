// Package repository handles the CRUD operations
package repository

import (
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JobDBItf interface {
	GetJobInfo(job *entity.Job) error
	GetJobList(offset *int, limit *int, job *[]entity.Job) error
	GetCompanyInfo(company *entity.Company) error
}

type JobDB struct {
	db *gorm.DB
}

func NewJobDB(db *gorm.DB) JobDBItf {
	return &JobDB{
		db: db,
	}
}

func (r *JobDB) GetJobInfo(job *entity.Job) error {
	return r.db.Debug().
		Preload(clause.Associations).
		First(job).
		Error
}

func (r *JobDB) GetJobList(offset *int, limit *int, job *[]entity.Job) error {
	return r.db.Debug().
		Preload(clause.Associations).
		Limit(*limit).
		Offset(*offset).
		Find(job).
		Error
}

func (r *JobDB) GetCompanyInfo(company *entity.Company) error {
	return r.db.Debug().
		Model(&entity.Company{}).
		First(company).
		Error
}
