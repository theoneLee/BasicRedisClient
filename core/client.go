package core

import (
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

}
