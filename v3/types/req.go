package types

import (
	"time"
)

type Req struct {
	ID           int       `json:"id" yaml:"id"`
	Time         time.Time `json:"time" yaml:"time"`
	ResponseTime int64     `json:"response_time" yaml:"response_time"`
}
