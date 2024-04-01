package model

type Job struct {
	Id            string `json:"id"`
	MaxRetryCount int    `json:"maxRetryCount"`
	CallbackURL   string `json:"callbackURL"`
	JobType       string `json:"jobType"`
	Config        string `json:"config"`
}

func (Job) TableName() string {
	return "job"
}
