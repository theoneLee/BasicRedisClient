package core

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
)

// Reply : redis response
type Reply struct {
	IsMulti     bool
	Err         error
	Conn        *net.TCPConn
	SingleReply []byte
	MultiReply  [][]byte
	Source      []byte
}

// Request : redis command request
func Request(args ...string) string {
	return multiCommandMarshal(args...)
}

func multiCommandMarshal(args ...string) string {
	s := "*"
	s += strconv.Itoa(len(args))
	s += "\r\n"

	for _, value := range args {
		s += "$"
		s += strconv.Itoa(len(value))
		s += "\r\n"
		s += value
		s += "\r\n"

	}

	return s
}

func (r *Reply) Reply() {
	reader := bufio.NewReader(r.Conn)
	b, err := reader.Peek(1) // pre read a first byte to judge what redis response
	if err != nil {
		fmt.Println("conn err")

	}

	if b[0] == byte('*') {
		r.IsMulti = true //it means redis response array ,see https://redis.io/topics/protocol#resp-arrays
		r.MultiReply, r.Err = multiResponse(reader)
	} else { //b[0] will be '+' '-' ':' '$'
		r.IsMulti = false
		r.SingleReply, err = singleResponse(reader)
		if err != nil {
			r.Err = err
			return
		}
	}

}

func singleResponse(reader *bufio.Reader) (res []byte, err error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return
	}
	switch prefix {
	case byte('+'), byte('-'), byte(':'):
		res, _, err = reader.ReadLine()
	case byte('$'):
		// $7\r\naiangwt\r\n
		n, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		l, err := strconv.Atoi(string(n))
		if err != nil {
			return
		}
		p := make([]byte, l+2)
		reader.Read(p)
		res = p[0 : len(p)-2]

	}
	return
}

func multiResponse(reader *bufio.Reader) (res [][]byte, err error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return
	}
	if prefix != byte('*') {
		return res, errors.New("not multi response")
	}
	//*3\r\n$1\r\n3\r\n$1\r\n2\r\n$1\r\n
	l, _, err := reader.ReadLine()
	if err != nil {
		return
	}
	n, err := strconv.Atoi(string(l))
	if err != nil {
		return
	}
	for i := 0; i < n; i++ {
		s, err := singleResponse(reader)
		fmt.Println("i =", i, "result = ", string(s))
		if err != nil {
			return
		}
		res = append(res, s)
	}

	return
}
