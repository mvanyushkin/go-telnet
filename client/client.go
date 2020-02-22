package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

func RunClient(addr string, duration time.Duration) {
	fmt.Printf("connecting to %v ... \n", addr)
	conn, err := net.DialTimeout("tcp", addr, duration)
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}

	fmt.Print("ok.\n")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go RunSendingWorker(conn, ctx, &wg)
	go RunReceivingWorker(conn, ctx, &wg)
	wg.Wait()
	cancel()
}

func FromReaderToWriter(reader *bufio.Reader, writer *bufio.Writer, ctx context.Context, wg *sync.WaitGroup, workerType string) {
	for {
		select {
		case <-ctx.Done():
			println("done")
			return

		default:
			input, err := reader.ReadString('\n')
			if err == io.EOF {
				fmt.Printf("Connection closed by %v \n", workerType)
				wg.Done()
				return
			}
			_, err = writer.WriteString(input)
			if err != nil {
				wg.Done()
				return
			}

			err = writer.Flush()
			if err != nil {
				wg.Done()
				return
			}
		}
	}
}

func RunReceivingWorker(conn net.Conn, ctx context.Context, wg *sync.WaitGroup) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(os.Stdout)
	FromReaderToWriter(reader, writer, ctx, wg, HOST)
}

func RunSendingWorker(conn net.Conn, ctx context.Context, wg *sync.WaitGroup) {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)
	FromReaderToWriter(reader, writer, ctx, wg, CLIENT)
}
