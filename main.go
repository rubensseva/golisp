package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"golisp/golisp"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		runeReader := bufio.NewReader(os.Stdin)

		for {
			tokenizer := golisp.NewTokenizer(runeReader)

			n := golisp.Parse(tokenizer)
			fmt.Printf("%+v\n", n)
			n.NodePprint()
			fmt.Println()

			fmt.Println(golisp.Eval(n))
		}
		time.Sleep(time.Millisecond * 100)
	}()

	wg.Wait()

}
