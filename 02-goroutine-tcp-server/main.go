package main

import (
	"math"
	"net"
	"os"
)

// CPU-intensive function to find prime numbers
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Function to simulate CPU-bound task
func cpuIntensiveTask() {
	count := 0
	for i := 2; i < 1000000; i++ { // Adjust the range for more/less intensity
		if isPrime(i) {
			count++
		}
	}
}

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

	req_count := 0
	// infinite loop (until interrupted)
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accepting connection:", err)
			continue
		}
		println("Accepted connection from", conn.RemoteAddr())

		// handle connection in a new goroutine
		// so we can continue listening for new connections
		// (multithreaded server)
		req_count++
		go handleConnection(conn, req_count)
	}
}

// handleConnection
func handleConnection(conn net.Conn, req_count int) {
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
	println("Worker", req_count, "is handling connection")
	println("Received", n, "bytes from", ip)

	// write to stdout
	os.Stdout.Write(buf[:n])

	// time.Sleep(10 * time.Second) // simulate delay
	cpuIntensiveTask()

	// for two connections at the same time, both will be delayed for 10 seconds
	// also, main thread is not blocked as there are two goroutines running

	// write to connection
	_, err = conn.Write([]byte("Echo server\n"))
	if err != nil {
		println("Error writing to connection:", err)
		return
	}
}
