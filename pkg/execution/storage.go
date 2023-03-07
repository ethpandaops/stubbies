package execution

import (
	"context"
	"encoding/json"
	"math/big"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	log logrus.FieldLogger

	hashMap     map[string]*Block
	numberMap   map[*big.Int]*Block
	latestBlock *Block

	mu sync.Mutex
}

func newStorage(log logrus.FieldLogger) *Storage {
	return &Storage{
		log: log,

		hashMap:   make(map[string]*Block),
		numberMap: make(map[*big.Int]*Block),
	}
}

func (s *Storage) Start(ctx context.Context) {
	if err := s.startCrons(ctx); err != nil {
		s.log.WithError(err).Fatal("Failed to start crons")
	}
}

func (s *Storage) startCrons(ctx context.Context) error {
	c := gocron.NewScheduler(time.Local)

	if _, err := c.Every("1m").Do(func() {
		s.cleanUp()
	}); err != nil {
		return err
	}

	c.StartAsync()

	return nil
}

func (s *Storage) cleanUp() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.latestBlock == nil {
		return
	}

	// delete blocks that are older than 6 epochs (192 blocks) from latest
	for number, block := range s.numberMap {
		diff := new(big.Int)
		diff = diff.Sub(s.latestBlock.Number, number)

		if diff.Cmp(big.NewInt(192)) > 0 {
			delete(s.numberMap, block.Number)
			delete(s.hashMap, block.payload.BlockHash)
		}
	}
}

func (s *Storage) GetBlockByHash(hash string) *Block {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.hashMap[hash]
}

func (s *Storage) GetBlockByNumber(number *big.Int) *Block {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.numberMap[number]
}

func (s *Storage) GetLatestBlock() *Block {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.latestBlock
}

func (s *Storage) AddBlock(payload *RequestParamsNewPayloadV1, raw *json.RawMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	num := new(big.Int)
	num.SetString(payload.BlockNumber[2:], 16)

	block := &Block{
		Number:  num,
		raw:     raw,
		payload: payload,
	}

	s.log.WithField("number", num).Error("Adding block to storage")

	s.hashMap[block.payload.BlockHash] = block
	s.numberMap[num] = block

	if s.latestBlock == nil || s.latestBlock.Number.Cmp(num) < 0 {
		s.latestBlock = block
	}
}
