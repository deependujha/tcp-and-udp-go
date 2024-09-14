package main

import (
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: tcp-server <port>")
		os.Exit(1)
		return
	}
	port := os.Args[1]
	println("Starting TCP server on port", port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		println("Error resolving TCP address:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		println("Error listening on TCP address:", err)
		return
	}

	// infinite loop (until interrupted)
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accepting connection:", err)
			continue
		}
		println("Accepted connection from", conn.RemoteAddr())

		handleConnection(conn)
	}
}

// handleConnection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// read from connection
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		println("Error reading from connection:", err)
		return
	}
	// get ip address of client
	ip := conn.RemoteAddr().String()
	println("Received", n, "bytes from", ip)

	// write to stdout
	os.Stdout.Write(buf[:n])

	time.Sleep(10 * time.Second) // simulate delay

	// for two connections at the same time, one will be delayed for 20 seconds (as no goroutines are running)
	// main thread itself is blocked until the first connection is closed

	// write to connection
	_, err = conn.Write([]byte("Echo server\n"))
	if err != nil {
		println("Error writing to connection:", err)
		return
	}
}
