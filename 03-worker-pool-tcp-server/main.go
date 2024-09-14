package main

import (
	"math"
	"net"
	"os"
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

	// create 100 worker goroutines
	connChan := make(chan net.Conn)
	reqCount := make(chan int)

	for i := 0; i < 100; i++ {
		go tcpWorker(i+1, connChan, reqCount)
	}
	count := 0
	// infinite loop (until interrupted)
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accepting connection:", err)
			continue
		}
		println("Accepted connection from", conn.RemoteAddr())

		// send connection to worker goroutine
		// which ever worker is available will receive the connection and handle it
		connChan <- conn
		reqCount <- count + 1
		count++
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

	// sleep for 10 second
	// time.Sleep(10 * time.Second)
	cpuIntensiveTask()

	// write to connection
	_, err = conn.Write([]byte("Echo server\n"))
	if err != nil {
		println("Error writing to connection:", err)
		return
	}
}

// tcpWorker
func tcpWorker(i int, connChan chan net.Conn, reqCount chan int) {
	for conn := range connChan {
		println("Worker", i, "is handling connection count", <-reqCount)
		handleConnection(conn)
	}
}

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
