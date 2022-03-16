package main

import (
	"log"
	"net"
	"runtime"
)

var helloWorldText = "HTTP/1.1 200 OK\r\nContext-Type: text/plain\r\nServer: tcp-http-implementation\r\n\r\nhello world\r\n"

func main() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println("http served on :80")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept err: %v\n", err)
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	var (
		writeCh   = make(chan struct{}, 1)
		releaseCh = make(chan struct{}, 1)
	)

	go readConn(conn, writeCh)

	go writeConn(conn, writeCh, releaseCh)

	go releaseConn(conn, releaseCh)
}

func readConn(conn net.Conn, writeCh chan<- struct{}) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		runtime.Goexit()
	}
	writeCh <- struct{}{}
}

func writeConn(conn net.Conn, writeCh <-chan struct{}, releaseCh chan<- struct{}) {
	<-writeCh
	_, err := conn.Write([]byte(helloWorldText))
	if err != nil {
		runtime.Goexit()
	}
	releaseCh <- struct{}{}
}

func releaseConn(conn net.Conn, releaseCh <-chan struct{}) {
	select {
	case <-releaseCh:
		_ = conn.Close()
		break
	}
}
