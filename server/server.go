package server

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/Toolnado/alligator/cache/interfaces"
	"github.com/Toolnado/alligator/commands"
)

type Options struct {
	Addr  string
	Cache interfaces.Cacher
}

type Server struct {
	Opts   Options
	leader bool
}

func New(ops Options, lead bool) *Server {
	return &Server{
		Opts:   ops,
		leader: lead,
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Opts.Addr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept connection error:", err)
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("read conn error:", err)
			return
		}

		bufMSG := buf[:n]
		go s.handleCommand(conn, bufMSG)
	}
}

func (s *Server) handleCommand(conn net.Conn, raw []byte) {
	cmd, err := commands.ParseCommand(raw)
	if err != nil {
		log.Println("parse command error:", err)
		return
	}
	switch cmd.Name {
	case commands.SET_COMMAND:
		log.Println("handling the set command: ", cmd)
		err = s.Opts.Cache.Set(cmd.Key, cmd.Value, cmd.TTL)
		if err != nil {
			log.Println("set command error:", err)
		}
	case commands.GET_COMMAND:
		log.Println("handling the get command: ", cmd)
		value, err := s.Opts.Cache.Get(cmd.Key)
		if err != nil {
			log.Println("get command error:", err)
			return
		}
		_, err = conn.Write(value)
		if err != nil {
			log.Println("write to conn error:", err)
		}
	case commands.DELETE_COMMAND:
		log.Println("handling the delete command: ", cmd)
		err = s.Opts.Cache.Delete(cmd.Key)
		if err != nil {
			log.Println("delete command error:", err)
			return
		}
	case commands.HAS_COMMAND:
		log.Println("handling the has command: ", cmd)
		exist := s.Opts.Cache.Has(cmd.Key)
		_, err = conn.Write([]byte(strconv.FormatBool(exist)))
		if err != nil {
			log.Println("write to conn error:", err)
		}
	}
}
