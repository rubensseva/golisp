package main

import (
	"bytes"
	"fmt"
	"golisp/channel"
	"golisp/parser"
	"log"
	"sync"
	"time"
)

func main() {
	b := bytes.NewReader([]byte("(+ 1 2)"))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		cout := channel.NewPeekableChannel(1024)
		err := parser.Tokenize(b, cout)
		if err != nil {
			log.Fatalln(err)
		}

		n := parser.Parse(cout)
		fmt.Printf("%+v\n", n)
		n.NodePprint()
		fmt.Println()

		fmt.Println(parser.Eval(n))
		time.Sleep(time.Millisecond * 100)
	}()

	wg.Wait()

}
