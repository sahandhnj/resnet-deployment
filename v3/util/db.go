package util

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

const (
	ErrNotFound = "Object was not found"
)

func GetError(errorMessage string) error {
	return errors.New(errorMessage)
}

// Itob returns an 8-byte big endian representation of v.
// This function is typically used for encoding integer IDs to byte slices
// so that they can be used as BoltDB keys.
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func CreateBucket(db *bolt.DB, bucketName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return nil
	})
}

func GetObject(db *bolt.DB, bucketName string, key []byte, object interface{}) error {
	var data []byte

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		value := bucket.Get(key)
		if value == nil {
			return GetError(ErrNotFound)
		}

		data = make([]byte, len(value))
		copy(data, value)

		return nil
	})

	if err != nil {
		return err
	}

	return UnmarshalJsonObject(data, object)
}

func UpdateObject(db *bolt.DB, bucketName string, key []byte, object interface{}) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		data, err := MarshalJsonObject(object)
		if err != nil {
			return err
		}

		err = bucket.Put(key, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func DeleteObject(db *bolt.DB, bucketName string, key []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Delete(key)
	})
}

func GetNextIdentifier(db *bolt.DB, bucketName string) int {
	var identifier int

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		id := bucket.Sequence()
		identifier = int(id)
		return nil
	})

	identifier++
	return identifier
}

func MarshalJsonObject(object interface{}) ([]byte, error) {
	return json.Marshal(object)
}

func UnmarshalJsonObject(data []byte, object interface{}) error {
	return json.Unmarshal(data, object)
}
