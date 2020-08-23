package blockchain

import (
	"errors"
	"fmt"

	"github.com/NatanColenberg/blockchain/database"
	"github.com/go-redis/redis/v8"
)

// BlockChain holds the Last-Hash of the last Block stored
// as a key to find that block in DB and a Database Client Connection
type BlockChain struct {
	LastHash []byte
	Database *redis.Client
}

// InitBlockChain initializes a new chain of blocks
func InitBlockChain() (*BlockChain, error) {

	// Connect to Database
	rdb, rdbConnErr := database.Connect("localhost:6379", "", 0)
	if rdbConnErr != nil {
		msg := fmt.Sprintf("Failed to connect to database (Error = %s)", rdbConnErr.Error())
		return nil, errors.New(msg)
	}

	// Get the Last-Hash
	lh, lhErr := database.GetLastHash()

	// Check if this is a new chain (lh does not exist)
	if lhErr == redis.Nil {

		// Create the first Block
		genesis, genesisErr := CreateBlock("Genesis", []byte{})
		if genesisErr != nil {
			msg := fmt.Sprintf("Failed to create first block (Error = %s)", genesisErr.Error())
			return nil, errors.New(msg)
		}

		// Serialize first Block
		genesisSerialized, serializedErr := genesis.Serialize()
		if serializedErr != nil {
			msg := fmt.Sprintf("Failed to serialize first block (Error = %s)", serializedErr.Error())
			return nil, errors.New(msg)
		}

		// Store first Block in DB
		setGenesisErr := database.SetBlock(genesis.Hash, genesisSerialized)
		if setGenesisErr != nil {
			msg := fmt.Sprintf("Failed to store first block in DB (Error = %s)", setGenesisErr.Error())
			return nil, errors.New(msg)
		}

	} else if lhErr != nil {
		msg := fmt.Sprintf("Failed to retrieve Last-Hash from DB (Error = %s)", lhErr.Error())
		return nil, errors.New(msg)
	}

	blockchain := BlockChain{lh, rdb}
	return &blockchain, nil
}

// AddBlock adds a new Block to the end of the chain
func (chain *BlockChain) AddBlock(data string) error {

	// Get Last-Hash from DB
	lastHash, lastHashErr := database.GetLastHash()
	if lastHashErr != nil {
		msg := fmt.Sprintf("Failed to retrieve Last-Hash from DB (Error = %s)", lastHashErr.Error())
		return errors.New(msg)
	}

	// Create New Block
	newBlock, newBlockErr := CreateBlock(data, lastHash)
	if newBlockErr != nil {
		msg := fmt.Sprintf("Failed to create new block (Error = %s)", newBlockErr.Error())
		return errors.New(msg)
	}

	// Serialize the New Block
	blockSerialized, serializedErr := newBlock.Serialize()
	if serializedErr != nil {
		msg := fmt.Sprintf("Failed to serialize new block (Error = %s)", serializedErr.Error())
		return errors.New(msg)
	}

	// Store New Block in DB
	storeBlockErr := database.SetBlock(newBlock.Hash, blockSerialized)
	if storeBlockErr != nil {
		msg := fmt.Sprintf("Failed to store new block in DB (Error = %s)", storeBlockErr.Error())
		return errors.New(msg)
	}

	// Update the Last-Hash
	chain.LastHash = newBlock.Hash

	return nil
}

// GetLastBlock retrieves that last block stored in the chain
func (chain *BlockChain) GetLastBlock() (*Block, error) {

	// Get the Last Block
	lastBlockData, lastBlockDataErr := database.GetLastBlock()
	if lastBlockDataErr != nil {
		msg := fmt.Sprintf("Failed to fatche last block from DB (Error = %s)", lastBlockDataErr.Error())
		return nil, errors.New(msg)
	}

	// Deserialize Last Block
	lastBlock, deserializeErr := Deserializer(lastBlockData)
	if deserializeErr != nil {
		msg := fmt.Sprintf("Failed to deserialize last block (Error = %s)", deserializeErr.Error())
		return nil, errors.New(msg)
	}
	return lastBlock, nil
}
