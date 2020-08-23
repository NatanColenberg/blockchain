package blockchain

import (
	"fmt"
	"errors"
	"encoding/json"

	"github.com/NatanColenberg/blockchain/database"
)

// Block represents a single block in a linked list
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// CreateBlock creats a new Block
func CreateBlock(data string, prevHash []byte) (*Block, error) {
	
	// Create Block
	block :=
		&Block{
			Data:     []byte(data),
			PrevHash: prevHash,
		}

	// Mine Block
	pow := NewProof(block)
	nonce, hash := pow.Run()

	// Update With mined data
	block.Hash = hash
	block.Nonce = nonce

	return block, nil
}

// Serialize a single Block for storing in DB
func (b *Block) Serialize() ([]byte, error) {

	blockData, marshalErr := json.Marshal(b)
	if marshalErr != nil {
		msg := fmt.Sprintf("Failed to Marshal Block (Error = %s)", marshalErr.Error())
		return nil, errors.New(msg)
	}

	return blockData, nil
}

// Deserializer a single Block for fetching from DB
func Deserializer(data []byte) (*Block, error) {
	var block Block

	unmarshalErr := json.Unmarshal(data, &block)
	if unmarshalErr != nil {
		msg := fmt.Sprintf("Failed to Unmarshal Block (Error = %s)", unmarshalErr.Error())
		return nil, errors.New(msg)
	}

	return &block, nil
}

// Prev iterates thought the chain of Blocks
func (b *Block) Prev() (*Block, error) {

	prevHash := string(b.PrevHash)

	// Validate End of Chain
	if prevHash == "" {
		return nil, errors.New("End of Chain")
	}

	// Get Prev Block from DB
	prevBlockData, err := database.GetBlock(b.PrevHash)
	if err != nil {
		msg := fmt.Sprintf("Failed to Get Block from DB (Error = %s)", err.Error())
		return nil, errors.New(msg)
	}

	// Deserializer Prev Block
	prevBlock, err := Deserializer([]byte(prevBlockData))

	return prevBlock, nil
}