package consensus

import (
	"bytes"
	"crypto/rand"
	"fmt"
)

func BruteForceFinder(c *Chain, stop chan int) {
	newBlockReady, id := c.Subscribe()
	var lastBlock, b *Block
L:
	for {
		select {
		case <-stop:
			break L
		case <-newBlockReady:
			lastBlock = c.GetLast()
			payload := make([]byte, 32)
			rand.Read(payload)
			cryptoHash := lastBlock.GetCryptoHash()
			b = &Block{msg: payload[:32], prev: cryptoHash, r: Randgen(&c.q), s: Randgen(&c.q)}
		default:
			cryptoHash := b.GetCryptoHash()
			var hash1 []byte
			ChameleonHash(&c.hk, &c.p, &c.q, &c.g, &cryptoHash, &b.r, &b.s, &hash1)
			if bytes.Equal(c.hash0[0:c.difficulty], hash1[0:c.difficulty]) {
				err := c.AddBock(b)
				if err != nil {
					//fmt.Println(err)
				}
			} else {
				b.r = Randgen(&c.q)
				b.s = Randgen(&c.q)
			}
		}
	}
	c.Unsubscribe(id)

}

func SmartFinder(c *Chain, stop chan int, tk []byte) {
	newBlockReady, id := c.Subscribe()
	firstBlock := c.GetFirst()
	r0 := firstBlock.r
	s0 := firstBlock.s
	cryptoHash0 := firstBlock.GetCryptoHash()
	var lastBlock, b *Block
L:
	for {
		select {
		case <-stop:
			break L
		case <-newBlockReady:
			lastBlock = c.GetLast()
			payload := make([]byte, 32)
			rand.Read(payload)
			cryptoHash := lastBlock.GetCryptoHash()
			b = &Block{payload[:32], cryptoHash, Randgen(&c.q), Randgen(&c.q)}
			var r1, s1 []byte
			cryptoHash1 := b.GetCryptoHash()
			GenerateCollision(&c.hk, &tk, &c.p, &c.q, &c.g, &cryptoHash0, &cryptoHash1, &r0, &s0, &r1, &s1)
			b.r = r1
			b.s = s1
			c.AddBock(b)
		default:
			fmt.Println("Should not be here")
		}
	}
	c.Unsubscribe(id)

}
