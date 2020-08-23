package database

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

var ctx context.Context
var client *redis.Client

// Connect is used to create and verify a connection with a Redis DB Server
func Connect(addr, pass string, db int) (*redis.Client, error) {

	// Create Context
	ctx = context.Background()

	// Create Client
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	// Ping Redis Server
	pong, pongErr := client.Ping(ctx).Result()
	if pongErr != nil || pong != "PONG" {
		return nil, errors.New("Failed to PING Redis server")
	}

	return client, nil
}

// SetLastHash Sets the Last-Hash of the last added Block
func SetLastHash(lastHash []byte) error {
	err := client.Set(ctx, "lh", string(lastHash), 0).Err()
	if err != nil {
		return errors.New("Failed to Set Lash-Hash")
	}
	return nil
}
// GetLastHash Gets the Last-Hash of the last added Block
func GetLastHash() ([]byte, error) {
	lh, lhErr := client.Get(ctx, "lh").Result()
	if lhErr == redis.Nil || lhErr != nil {
		return nil, lhErr
	}

	return []byte(lh), nil
}

// SetBlock adds a new Block to the DB
func SetBlock(hash, block []byte) error {
	err := client.Set(ctx, string(hash), string(block), 0).Err()
	if err != nil {
		return errors.New("Failed to Set Block")
	}

	err = SetLastHash(hash)
	if err != nil {
		return errors.New("Failed to Set Lash-Hash")
	}

	return nil
}

// GetBlock retrieves a Block to the DB
func GetBlock(hash []byte) (string, error) {
	block, blockErr := client.Get(ctx, string(hash)).Result()
	if blockErr != nil {
		return "", blockErr
	}

	return block, nil
}


// GetLastBlock retrieves the last stored Block to the DB
func GetLastBlock() ([]byte, error) {
	// Get Last-Hash
	lh, lhErr := GetLastHash()
	if lhErr != nil {
		return nil, lhErr
	}

	// Get Last Block
	block, blockErr := client.Get(ctx, string(lh)).Result()
	if blockErr != nil {
		return nil, blockErr
	}

	return []byte(block), nil
}

