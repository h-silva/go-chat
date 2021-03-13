package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	Nickname string
	msg      chan<- string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	log.Println("Welcome to GO-CHAT")

	listener, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		panic(err.Error())
	}

	go broadcaster()

	for {
		cnn, err := listener.Accept()

		if err != nil {
			log.Println(err.Error())
			continue
		}

		go handleConn(cnn)
	}
}

func broadcaster() {
	clients := make(map[client]bool)

	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.msg <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.msg)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)

	go clientWritter(conn, ch)

	who := conn.RemoteAddr().String()

	cli := client{Nickname: who, msg: ch}

	entering <- cli

	input := bufio.NewScanner(conn)

	for input.Scan() {
		messages <- cli.Nickname + " said: " + input.Text()
	}

	leaving <- cli
	conn.Close()
}

func clientWritter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
