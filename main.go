package main

import (
	"bytes"
	"fmt"
	"golisp/parser"
	"golisp/channel"
	"log"
	"sync"
	"time"
)

func main() {
	b := bytes.NewReader([]byte("(println (+ 1 2 (/ 1 2) (- 2 (* 4 3 2)) 5))"))


	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		cout := channel.NewPeekableChannel(1024)
		err  := parser.Tokenize(b, cout)
		if err != nil {
			log.Fatalln(err)
		}

		n := parser.Parse(cout)
		fmt.Printf("%+v\n", n)
		n.NodePprint()
		fmt.Println()
		time.Sleep(time.Second * 10)
	}()

	wg.Wait()

}
