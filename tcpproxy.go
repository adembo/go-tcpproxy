package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func log(s string, args ...interface{}) {
	if os.Getenv("LOG") != "" {
		fmt.Printf(s, args)
	}
}

func check_err(e error, s string) {
	if e != nil {
		fmt.Printf(s, e.Error())
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <local_port> <server_host> <server_port>\n", os.Args[0])
		os.Exit(2)
	}
	log("Invoked with arguments: %s\n", os.Args[1:])
	listen_port := os.Args[1]
	server_host := os.Args[2]
	server_port := os.Args[3]

	local_hostport := net.JoinHostPort("", listen_port)
	remote_hostport := net.JoinHostPort(server_host, server_port)

	laddr, err := net.ResolveTCPAddr("tcp", local_hostport)
	check_err(err, "Error in resolve local: %s\n")
	raddr, err := net.ResolveTCPAddr("tcp", remote_hostport)
	check_err(err, "Error in resolve remote: %s\n")

	ln, err := net.ListenTCP("tcp", laddr)
	check_err(err, "Error in listen: %s\n")
	log("Listening on %s\n", local_hostport)
	for {
		conn, err := ln.AcceptTCP()
		check_err(err, "Error in accept: %s\n")
		go handle_connection(raddr, conn)
	}
}

func handle_connection(raddr *net.TCPAddr, in *net.TCPConn) {
	log("Accepted a connection from %s\n", in.RemoteAddr().String())

	// Connect to remote host.
	out, err := net.DialTCP("tcp", nil, raddr)
	check_err(err, "Error in dial: %s\n")

	// Create channels and kick-off reader and writer.
	ch1 := make(chan int64)
	ch2 := make(chan int64)
	go handle_io(ch1, in, out)
	go handle_io(ch2, out, in)
	total := <-ch1
	<-ch2
	log("Done with connection: %d bytes written\n", total)
}

func handle_io(c chan int64, first, second *net.TCPConn) {
	var total int64
	for {
		bytes, err := io.Copy(first, second)
		check_err(err, "Error in copy: %s\n")
		if bytes == 0 {
			second.CloseRead()
			first.CloseWrite()
			break
		}
		total += bytes
	}
	c <- total
}
