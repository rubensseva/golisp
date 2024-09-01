package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golisp/channel"
	"golisp/parser"
	"log"
	"os"
	"sync"
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

		runeReader := bufio.NewReader(os.Stdin)

		tokenizer := parser.NewTokenizerv2(runeReader)

		for {
			s, err := tokenizer.Token()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(s)
		}


		// n := parser.Parse(cout)
		// fmt.Printf("%+v\n", n)
		// n.NodePprint()
		// fmt.Println()

		// fmt.Println(parser.Eval(n))
		// time.Sleep(time.Millisecond * 100)
	}()

	wg.Wait()

}
