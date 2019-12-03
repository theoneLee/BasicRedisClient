package main

import (
	"basic_redis_client/core"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	port := "6379"
	host := "localhost"

	conn := client(host, port) // new tcp client to connect redis
	defer conn.Close()

	for {
		fmt.Printf("%s:%s>", host, port)
		bio := bufio.NewReader(os.Stdin)
		input, _, err := bio.ReadLine()
		if err != nil {
			fmt.Println(err)
		}
		s := strings.Split(string(input), " ")
		req := core.Request(s...)
		conn.Write([]byte(req))
		reply := core.Reply{}
		reply.Conn = conn
		reply.Reply()

		if reply.Err != nil {
			fmt.Println("err:", reply.Err)
		}
		var res []byte
		if reply.IsMulti {
			for index, _ := range reply.MultiReply {
				res = reply.MultiReply[index]
				fmt.Println("result:", string(res), "\nerr:", err)
			}
		} else {
			res = reply.SingleReply
			fmt.Println("result:", string(res), "\nerr:", err)
		}
	}
}

func client(host, port string) *net.TCPConn {
	porti, err := strconv.Atoi(port)
	if err != nil {
		panic("port is error")
	}
	tcpAddr := &net.TCPAddr{IP: net.ParseIP(host), Port: porti}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(err)
	}
	return conn
}
