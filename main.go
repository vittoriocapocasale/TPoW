package main

import (
	"fmt"
	"time"

	"github.com/vittoriocapocasale/trapPoWsim/consensus"
)

func main() {
	difficulty := 2
	peers := 10
	trusted := 1
	var p, q, g, hk, tk []byte
	consensus.Keygen(256, &p, &q, &g, &hk, &tk)
	chain := consensus.NewChain(hk, p, q, g, difficulty)
	channels := []chan int{}
	for i := 0; i < peers; i++ {
		ch := make(chan int, 1)
		channels = append(channels, ch)
		go consensus.BruteForceFinder(chain, ch)
	}
	for i := 0; i < trusted; i++ {
		ch := make(chan int, 1)
		channels = append(channels, ch)
		go consensus.SmartFinder(chain, ch, tk)
	}

	time.Sleep(120 * time.Second)
	for i := 0; i < len(channels); i++ {
		channels[i] <- 1
	}
	fmt.Println("Total Blocks: ", len(chain.GetBlocks()))
	//go BruteForceFinder(chain)
	//go BruteForceFinder(chain)
	//go BruteForceFinder(chain)
	//consensus.BruteForceFinder(chain)
	//consensus.SmartFinder(chain, tk)
}

/*
func main() {
	// Generate the parameters.
	var p, q, g, hk, tk, hash1, hash2, r1, s1, r2, s2, msg1, msg2 []byte

	Keygen(256, &p, &q, &g, &hk, &tk)

	msg1 = []byte("YES")
	msg2 = []byte("oooooo")

	r1 = Fixedgen(&q)
	s1 = Randgen(&q)

	fmt.Printf("CHAMELEON HASH PARAMETERS:"+
		"\np: %s1"+
		"\nq: %s1"+
		"\ng: %s1"+
		"\nhk: %s1"+
		"\ntk: %s1"+
		"\nDONE!", p, q, g, hk, tk)

	// First we generate a chameleon hash.
	ChameleonHash(&hk, &p, &q, &g, &msg1, &r1, &s1, &hash1)

	fmt.Printf("\n\nROUND 1:"+
		"\nmsg1: %s"+
		"\nr1: %s1"+
		"\ns1: %s1"+
		"\nhash1: %x\n",
		msg1, r1, s1, hash1)

	fmt.Printf("\n\nGENERATING COLLISION...\n\n")

	// Now we need to generate a collision.
	start := time.Now()
	for i := 0; i < 10000; i++ {
		GenerateCollision(&hk, &tk, &p, &q, &g, &msg1, &msg2, &r1, &s1, &r2, &s2)
	}
	end := time.Now()
	fmt.Println(end.Sub(start) / 10000)

	ChameleonHash(&hk, &p, &q, &g, &msg2, &r2, &s2, &hash2)

	fmt.Printf("\nROUND 2:"+
		"\nmsg2: %s"+
		"\nr2: %s"+
		"\ns2: %s"+
		"\nhash2: %x\n",
		msg2, r2, s2, hash2)
}
*/
