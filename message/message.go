package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
    "time"

	"github.com/quic-go/quic-go"
)

var (
    shutdown = make(chan struct{})
    maxConnections = 5
    connSemaphore = make(chan struct{}, maxConnections)
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter Mode (server/client): ")
    mode, _ := reader.ReadString('\n')
    mode = strings.TrimSpace(mode)

    switch mode {
    case "server":
        runServer()
    case "client":
        runClient()
    default:
        fmt.Println("Invalid Mode")
    }
}

// Server function
func runServer() {
	listener, err := quic.ListenAddr("localhost:4242", generateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	fmt.Println("Server started on localhost:4242")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
            select {
                case <- shutdown:
                    return
                default:
                    log.Printf("Failed to accept connection: %v", err)
                    continue
                }
		}

        select {
        // accept connection if within limit
        case connSemaphore <- struct{}{}:
            fmt.Println("New connection accepted")
            go handleConnection(conn)
        // reject connection if limit reached
        default:
            fmt.Println("Maximum Connections Reached")

            // send client rejection message
            stream, err := conn.OpenStreamSync(context.Background())
			if err == nil {
				_, err = stream.Write([]byte("Connection rejected: maximum connection limit reached"))
				if err != nil {
					log.Printf("Failed to send rejection message: %v", err)
				}
				stream.Close()
			}
            conn.CloseWithError(0, "Maximum Connection Reached")
        }

	}
}

func handleConnection(conn quic.Connection) {
    // Release slot when connection is done
    defer func() {
        <- connSemaphore
    }()

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("Failed to accept stream: %v", err)
            return
		}

		go func(s quic.Stream) {
			defer s.Close()

			for {
				buf := make([]byte, 1024)
				n, err := s.Read(buf)
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading from stream: %v", err)
                        return
					}
				}

				message := string(buf[:n])
				fmt.Printf("Received: %s\n", message)
			}
		}(stream)
	}
}

// Client function
func runClient() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

    // configures max time connection stays alive
	quicConf := &quic.Config{
        MaxIdleTimeout: 60 * time.Second,
    }

	ctx := context.Background()
	conn, err := quic.DialAddr(ctx, "localhost:4242", tlsConf, quicConf)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	fmt.Println("Connected to server")

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		log.Fatalf("Failed to open stream: %v", err)
	}

    go func() {
		for {
			buf := make([]byte, 1024)
			n, err := stream.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from stream: %v", err)
				}
				return
			}
			message := string(buf[:n])
			fmt.Printf("Client received: %s\n", message)

			if message == "Connection rejected: maximum connection limit reached" {
				fmt.Println("Client connection rejected: maximum connection limit reached")
				stream.Close()
				conn.CloseWithError(0, "Client exiting due to max connections reached")
			}
		}
	}()


	for {
		fmt.Print("Enter message: ")
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		message = strings.TrimSpace(message)

		_, err := stream.Write([]byte(message))
        fmt.Println("Sent Message")
		if err != nil {
			log.Fatalf("Failed to send data: %v", err)
		}

        if message == "exit" {
            conn.CloseWithError(0x42, "Client Closed Connection")
            os.Exit(0)
        }
	}
}

// generateTLSConfig generates a TLS configuration for secure communication
func generateTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Failed to load TLS certificates: %v", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-echo-example"},
	}
}

