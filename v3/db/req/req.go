package req

import (
	"github.com/boltdb/bolt"
	"github.com/sahandhnj/ml-deployment-benchmarks/v3/types"
	"github.com/sahandhnj/ml-deployment-benchmarks/v3/util"
)

const (
	BucketName = "req"
)

type Service struct {
	db *bolt.DB
}

func NewService(db *bolt.DB) (*Service, error) {
	err := util.CreateBucket(db, BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		db: db,
	}, nil
}

func (s *Service) Req(ID int) (*types.Req, error) {
	var req types.Req
	identifier := util.Itob(int(ID))

	err := util.GetObject(s.db, BucketName, identifier, &req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func (s *Service) Reqs() ([]types.Req, error) {
	var reqs = make([]types.Req, 0)

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var req types.Req
			err := util.UnmarshalJsonObject(v, &req)
			if err != nil {
				return err
			}
			reqs = append(reqs, req)
		}

		return nil
	})

	return reqs, err
}

func (s *Service) GetNextIdentifier() int {
	return util.GetNextIdentifier(s.db, BucketName)
}

func (s *Service) CreateReq(req *types.Req) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		err := bucket.SetSequence(uint64(req.ID))
		if err != nil {
			return err
		}

		data, err := util.MarshalJsonObject(req)
		if err != nil {
			return err
		}

		return bucket.Put(util.Itob(int(req.ID)), data)
	})
}

func (s *Service) UpdateReq(ID int, req *types.Req) error {
	identifier := util.Itob(int(ID))
	return util.UpdateObject(s.db, BucketName, identifier, req)
}

func (s *Service) DeleteReq(ID int) error {
	identifier := util.Itob(int(ID))
	return util.DeleteObject(s.db, BucketName, identifier)
}
