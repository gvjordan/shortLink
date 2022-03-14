package main

import (
	"time"
)

type rateUser struct {
	Count int
	Time  int
}

type rateLimit struct {
	current map[string]*rateUser
}

func (r *rateLimit) inc(key string) {
	if _, ok := r.current[key]; ok {
		r.current[key].Count++
	} else {
		r.current[key] = &rateUser{
			Count: 1,
			Time:  int(time.Now().Unix()),
		}
	}
}

func (r *rateLimit) expire(key string, expire int) {
	time.Sleep(time.Duration(expire) * time.Second)
	if _, ok := r.current[key]; ok {
		if r.current[key].Count == 0 || r.current[key].Count == 1 {
			delete(r.current, key)
		} else {
			r.current[key].Count--
		}
	}
}

func newRateLimit() *rateLimit {
	return &rateLimit{
		current: make(map[string]*rateUser),
	}
}

func setRateLimit(limit int, expire int) {
	limitByIP = limit
	limitByIPExpire = expire
}
