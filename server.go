package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

var games = make(map[string]game)

type gamer struct {
	score byte
	pos   uint16
	conn  net.Conn
}
type ball struct {
	xpos uint16
	ypos uint16
}
type game struct {
	balls  []ball
	gamers [2]*gamer
}

func logErr(e error) {
	os.Stderr.WriteString(e.Error())
}

func handleConnection(conn net.Conn) {
	fmt.Println("Client connected:", conn.RemoteAddr())
	conn.Write([]byte("fuck indonesia" + "\n"))
	for {
		s, err := readString(conn)
		if err != nil {
			logErr(err)
		}
		if len(s) > 1 {
			g, ok := games[s]
			if !ok {
				games[s] = game{}
			}
			fmt.Println(s)
			conn.Write([]byte("received gameid\n"))
			numGamers := len(games[s].gamers)
			if numGamers >= 2 {
				return
			}
			var gameAdd *gamer = &gamer{}
			gameAdd.conn = conn
			g.gamers[0] = gameAdd
			break
		}
	}
	for {
		var buffer []byte = make([]byte, 1024)
		numread, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				fmt.Println("Closed connection to", conn.RemoteAddr())
				return
			}
			logErr(err)
			return
		}
		if numread > 0 {
			fmt.Print(string(buffer[:numread]))
			conn.Write([]byte("ACK" + "\n"))
		}
	}
}

func startGame(game) {

}

func readString(conn net.Conn) (string, error) {
	var buffer []byte = make([]byte, 1024)
	numread, err := conn.Read(buffer)
	if err != nil {
		if err == io.EOF {
			conn.Close()
			fmt.Println("Closed connection to", conn.RemoteAddr())
			return "", err
		}
		logErr(err)
		return "", err
	}
	if numread > 0 {
		return string(buffer[:numread]), nil
	}
	return "", errors.New("nothing receieved")
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn)
	}
}
