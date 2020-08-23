package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"
)

// Difficulty is how many bits of '0' 
// the mined hash has to start with
const Difficulty = 16

// ProofOfWork is a struct that helps up mine a Block
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof create a new Proof of Work
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{Block: b, Target: target}
	return pow
}

// InitData initializes the Proof of Work
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.Data,
			pow.Block.PrevHash,
			toHex(int64(nonce)),
			toHex(int64(Difficulty)),
		},
		[]byte{})
	return data
}

// Run mines a single Block
func (pow *ProofOfWork) Run() (int, []byte) {
	var initHash big.Int
	var hash [32]byte

	nonce := 0

	start := time.Now()

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		initHash.SetBytes(hash[:])

		if initHash.Cmp(pow.Target) == -1 {
			elapsed := time.Since(start)
			fmt.Printf(" - %s", elapsed)
			break
		} else {
			nonce++
		}
	}

	fmt.Println()

	return nonce, hash[:]
}

// Validate validates a single mined Block
func (pow *ProofOfWork) Validate() bool {
	var initHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	initHash.SetBytes(hash[:])

	return initHash.Cmp(pow.Target) == -1
}

func toHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
