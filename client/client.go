package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
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
	go RunSendingWorker(conn, ctx, cancel)
	go RunReceivingWorker(conn, ctx, cancel)

	select {
	case <-ctx.Done():
		println("======== BYE ========")
	}
}

func FromReaderToWriter(reader *bufio.Reader, writer *bufio.Writer, ctx context.Context, cancel context.CancelFunc, workerType string) {
	for {
		select {
		case <-ctx.Done():
			println("done")
			return

		default:
			input, err := reader.ReadString('\n')
			if err == io.EOF {
				fmt.Printf("Connection closed by %v \n", workerType)
				cancel()
				return
			}
			_, err = writer.WriteString(input)
			if err != nil {
				cancel()
				return
			}

			err = writer.Flush()
			if err != nil {
				cancel()
				return
			}
		}
	}
}

func RunReceivingWorker(conn net.Conn, ctx context.Context, cancel context.CancelFunc) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(os.Stdout)
	FromReaderToWriter(reader, writer, ctx, cancel, HOST)
}

func RunSendingWorker(conn net.Conn, ctx context.Context, cancel context.CancelFunc) {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)
	FromReaderToWriter(reader, writer, ctx, cancel, CLIENT)
}
