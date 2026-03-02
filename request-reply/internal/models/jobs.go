package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type Job struct{
	JobID string `pg:"job_id,unique"`
	Status string `pg:"status"`
	CreatedAt time.Time `pg:"created_at"`
}

func InsertJob(pg *pg.DB, job *Job)error{
	_,err:=pg.Model(job).Insert()
	return err
}