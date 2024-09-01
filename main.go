package main

import (
	"bufio"
	"fmt"
	"golisp/parser"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		runeReader := bufio.NewReader(os.Stdin)

		for {
			tokenizer := parser.NewTokenizerv2(runeReader)

			n := parser.Parse(tokenizer)
			fmt.Printf("%+v\n", n)
			n.NodePprint()
			fmt.Println()

			fmt.Println(parser.Eval(n))
		}
		time.Sleep(time.Millisecond * 100)
	}()

	wg.Wait()

}
