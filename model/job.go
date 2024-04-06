package model

type Job struct {
	Id            string `json:"id"`
	MaxRetryCount int    `json:"maxRetryCount"`
	CallbackUrl   string `json:"callbackUrl"`
	JobType       string `json:"jobType"`
	Config        string `json:"config"`
}

func (Job) TableName() string {
	return "job"
}
