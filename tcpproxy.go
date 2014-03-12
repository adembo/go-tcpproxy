package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <local_port> <server_host> <server_port>\n", os.Args[0])
		os.Exit(2)
	}
	fmt.Printf("Invoked with arguments: %s\n", os.Args[1:])
	listen_port := os.Args[1]
	server_host := os.Args[2]
	server_port := os.Args[3]

	hostport := net.JoinHostPort("", listen_port)
	ln, err := net.Listen("tcp", hostport)
	if err != nil {
		fmt.Printf("Error in listen: %s\n", err.Error())
		return
	}
	fmt.Printf("Listening on %s\n", hostport)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error in accept: %s\n", err.Error())
			continue
		}
		go handle_connection(server_host, server_port, conn)
	}
}

func handle_connection(server_host, server_port string, in net.Conn) {
	fmt.Printf("Accepted a connection\n")

	// Connect to remote host.
	out, err := net.Dial("tcp", net.JoinHostPort(server_host, server_port))
	if err != nil {
		fmt.Printf("Error in dial: %s\n", err.Error())
		in.Close()
		return
	}

	// Create channels and kick-off reader and writer.
	ch1 := make(chan int64)
	ch2 := make(chan int64)
	go handle_io(ch1, in, out)
	go handle_io(ch2, out, in)
	total := <-ch1
	<-ch2

	// Close everything.
	out.Close()
	in.Close()
	fmt.Printf("Done with connection: %d bytes written\n", total)
}

func handle_io(c chan int64, first, second net.Conn) {
	var total int64
	for {
		bytes, err := io.Copy(first, second)
		if err != nil {
			fmt.Printf("Error in copy: %s\n", err.Error())
			break
		}
		if bytes == 0 {
			break
		}
		total += bytes
	}
	c <- total
}
