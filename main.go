package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/KaviiSuri/pokergg/deck"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		d := deck.New()
		fmt.Println(d)
	}
}
