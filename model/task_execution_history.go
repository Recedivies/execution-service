package model

import "time"

type TaskExecutionHistory struct {
	Id            string    `json:"id"`
	RetryCount    int       `json:"retryCount"`
	Status        string    `json:"status"`
	JobId         string    `json:"jobId"`
	ExecutionTime time.Time `json:"executionTime"`
}

func (TaskExecutionHistory) TableName() string {
	return "task_execution_history"
}
