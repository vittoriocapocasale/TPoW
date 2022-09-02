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

}
