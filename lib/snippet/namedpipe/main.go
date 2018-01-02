package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pagarme/warp-pipe/lib/namedpipe"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <filename>\n", os.Args[0])
		return
	}

	var (
		filename      = os.Args[1]
		pipe          = namedpipe.New(filename)
		ctx, cancel   = context.WithCancel(context.Background())
		lineCh, errCh = pipe.Loop(ctx)
		sigCh         = make(chan os.Signal, 1)
		line          string
		err           error
	)
	defer cancel()

	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	for ok := true; ok; {
		select {
		case line, ok = <-lineCh:
			fmt.Println(line)
		case err, ok = <-errCh:
			fmt.Println(err)
			return
		case sig := <-sigCh:
			fmt.Printf("signal: %s(%d)\n", sig, sig)
			switch sig {
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGTERM:
				return
			}
		}
	}
}
