package service

import (
	"time"

	"github.com/sahandhnj/ml-deployment-benchmarks/v3/db"
	"github.com/sahandhnj/ml-deployment-benchmarks/v3/types"
)

type ReqService struct {
	DBHandler *db.DBStore
}

func NewReqService(db *db.DBStore) *ReqService {
	return &ReqService{
		DBHandler: db,
	}
}

func (rs *ReqService) Add(t time.Time, responseTime int64) error {
	req := &types.Req{
		ID:           rs.DBHandler.ReqService.GetNextIdentifier(),
		Time:         t,
		ResponseTime: responseTime,
	}

	rs.DBHandler.ReqService.CreateReq(req)
	return nil
}

type ReqDataResp struct {
	Count   int   `json:"count"`
	Average int64 `json:"average"`
}

func (rs *ReqService) Stat() ReqDataResp {
	requests, _ := rs.DBHandler.ReqService.Reqs()
	var sum int64 = 0
	for _, r := range requests {
		sum = sum + r.ResponseTime
	}

	count := len(requests)
	average := sum / int64(count)

	resp := ReqDataResp{
		Count:   count,
		Average: average,
	}

	return resp
}
