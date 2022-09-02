package consensus

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	msg  []byte
	prev []byte
	r    []byte
	s    []byte
}

func (b *Block) GetCryptoHash() []byte {
	hash := sha256.Sum256(append(b.msg, b.prev...))
	return hash[:]
}

type Chain struct {
	mu          sync.Mutex
	blocks      []*Block
	hash0       []byte
	hk          []byte
	p           []byte
	q           []byte
	g           []byte
	difficulty  int
	subscribers map[uint]chan int
	SubID       uint
}

func NewChain(hk []byte, p []byte, q []byte, g []byte, difficulty int) *Chain {
	msg := [32]byte{}
	prev := [32]byte{}
	r := Randgen(&q)
	s := Randgen(&q)
	genesis := &Block{msg[:], prev[:], r, s}
	cryptoHash := genesis.GetCryptoHash()
	var hash0 []byte
	ChameleonHash(&hk, &p, &q, &g, &cryptoHash, &r, &s, &hash0)
	blocks := []*Block{genesis}
	return &Chain{blocks: blocks, hash0: hash0, hk: hk, p: p, q: q, g: g, difficulty: difficulty, subscribers: make(map[uint]chan int), SubID: 0}
}

func (c *Chain) GetLast() *Block {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.blocks[len(c.blocks)-1]
}

func (c *Chain) GetFirst() *Block {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.blocks[0]
}

func (c *Chain) GetBlocks() []*Block {
	return c.blocks
}

func (c *Chain) AddBock(b *Block) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	//checking if b.prev contains the hash of the current chain head
	last := c.blocks[len(c.blocks)-1]
	cryptoHash := last.GetCryptoHash()
	if !bytes.Equal(b.prev, cryptoHash) {
		return fmt.Errorf("outdated block")
	}
	//checking if the chameleon hash of the cryptographic hash of b is close enough to the genesis chameleon hash
	cryptoHash = b.GetCryptoHash()
	var hash1 []byte
	ChameleonHash(&c.hk, &c.p, &c.q, &c.g, &cryptoHash, &b.r, &b.s, &hash1)
	if !bytes.Equal(c.hash0[0:c.difficulty], hash1[0:c.difficulty]) {
		return fmt.Errorf("invalid block")
	}
	//appending block and notifying subscribers
	c.blocks = append(c.blocks, b)
	//fmt.Println(len(c.blocks))
	for _, ch := range c.subscribers {
		if len(ch) < cap(ch) {
			ch <- 1
		}
	}
	return nil
}

func (c *Chain) Subscribe() (chan int, uint) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch := make(chan int, 1)
	ch <- 1
	id := c.SubID
	c.SubID++
	c.subscribers[id] = ch
	return ch, id
}

func (c *Chain) Unsubscribe(id uint) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.subscribers, id)
	return
}
