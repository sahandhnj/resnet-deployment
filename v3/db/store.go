package db

import (
	"os"
	"os/user"
	"path"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sahandhnj/ml-deployment-benchmarks/v3/db/req"
)

const (
	databaseFileName = "imagedetector.db"
	DBStorePath      = ".imagedetector"
	TIMEOUT_SECONDS  = 1
)

type DBStore struct {
	path       string
	db         *bolt.DB
	ReqService *req.Service
}

func NewDBStore() (*DBStore, error) {
	usr, err := user.Current()
	databaseDir := path.Join(usr.HomeDir, DBStorePath)
	databasePath := path.Join(databaseDir, databaseFileName)

	if _, err := os.Stat(databaseDir); os.IsNotExist(err) {
		err = os.Mkdir(databaseDir, 0700)
		if err != nil {
			return nil, err
		}
	}

	dbStore := &DBStore{
		path: databasePath,
	}

	err = dbStore.Open()
	if err != nil {
		return nil, err
	}

	err = dbStore.initServices()
	if err != nil {
		return nil, err
	}

	return dbStore, nil
}

func (d *DBStore) Open() error {
	db, err := bolt.Open(d.path, 0600, &bolt.Options{Timeout: TIMEOUT_SECONDS * time.Second})
	if err != nil {
		return err
	}
	d.db = db

	return d.initServices()
}

func (d *DBStore) Close() error {
	if d.db != nil {
		return d.db.Close()
	}

	return nil
}

func (d *DBStore) initServices() error {
	reqDBService, err := req.NewService(d.db)
	if err != nil {
		return err
	}

	d.ReqService = reqDBService

	return nil
}
