package cron

import "time"

type Store interface {
	Get(key string) int64
	Set(key string)
}

type Job struct {
	interval  int64
	fromZero  bool
	lastRun   int64
	nextRun   int64
	weeks     []string
	ats       []string
	callbacks []func()
}

type Scheduler struct {
	jobs []*Job
}

func NewScheduler() *Scheduler {

}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(time.Second * 1)
	// 每個 Tick。
	go func() {
		for t := range ticker.C {
			// 每個工作。
			go func() {
				for _, j := range s.jobs {
					// 每個回呼函式。
					go func(j *Job) {
						for _, c := range j.callbacks {
							c()
						}
					}(j)
				}
			}()
		}
	}()
}

func (s *Scheduler) Clear() {

}

func (j *Job) Do(func()) {

}

func (j *Job) FromZero() (job *Job) {

}

func (j *Job) Second() (job *Job) {

}

func (j *Job) Seconds() (job *Job) {

}

func (j *Job) Minute() (job *Job) {

}

func (j *Job) Minutes() (job *Job) {

}

func (j *Job) Hour() (job *Job) {

}

func (j *Job) Hours() (job *Job) {

}

func (j *Job) Day() (job *Job) {

}

func (j *Job) Days() (job *Job) {

}

func (j *Job) Week() (job *Job) {

}

func (j *Job) Weeks() (job *Job) {

}

func (j *Job) Monday() (job *Job) {

}

func (j *Job) Tuesday() (job *Job) {

}

func (j *Job) Wednesday() (job *Job) {

}

func (j *Job) Thursday() (job *Job) {

}

func (j *Job) Friday() (job *Job) {

}

func (j *Job) Saturday() (job *Job) {

}

func (j *Job) Sunday() (job *Job) {

}

func (j *Job) Month() (job *Job) {

}

func (j *Job) Months() (job *Job) {

}

func (j *Job) Year() (job *Job) {

}

func (j *Job) Years() (job *Job) {

}

func (j *Job) At(time string) (job *Job) {

}
