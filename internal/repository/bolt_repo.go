package repository

import (
	"errors"
	"strconv"

	"github.com/boltdb/bolt"
)

type BucketName string

const (
	AccessTokens  BucketName = "AccessTokens"
	RequestTokens BucketName = "RequestTokens"
)

type TokenStorage struct {
	db *bolt.DB
}

func NewTokenStorage(db *bolt.DB) *TokenStorage {
	return &TokenStorage{db: db}
}

func (s *TokenStorage) SaveAccessToken(chatID int64, token string) error {
	return s.SaveToken(chatID, token, AccessTokens)
}

func (s *TokenStorage) SaveRequestToken(chatID int64, token string) error {
	return s.SaveToken(chatID, token, RequestTokens)
}

func (s *TokenStorage) GetAccessToken(chatID int64) (string, error) {
	return s.GetToken(chatID, AccessTokens)
}

func (s *TokenStorage) GetRequestToken(chatID int64) (string, error) {
	return s.GetToken(chatID, RequestTokens)
}

// Utilitary methods for bolt
func (s *TokenStorage) SaveToken(chatID int64, token string, bucket BucketName) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatID), []byte(token))
	})
}

func (s *TokenStorage) GetToken(chatID int64, bucket BucketName) (string, error) {
	var result string

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		result = string(b.Get(intToBytes(chatID)))
		return nil
	})

	if result == "" {
		return "", errors.New("not found")
	}

	return result, err
}

func intToBytes(input int64) []byte {
	return []byte(strconv.FormatInt(input, 10))
}

// Initialize db and create buckets
func InitBolt(DBFilename string) (*bolt.DB, error) {
	db, err := bolt.Open(DBFilename, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(RequestTokens))
		return err
	}); err != nil {
		return nil, err
	}

	return db, nil
}
