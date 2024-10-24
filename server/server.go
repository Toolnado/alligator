package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/Toolnado/alligator/cache/interfaces"
	"github.com/Toolnado/alligator/commands"
)

type Options struct {
	Addr  string
	Cache interfaces.Cache
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
		removeBytes := strings.TrimSuffix(string(bufMSG), "\x0d\x0a")
		parts := strings.Split(removeBytes, " ")
		if len(parts) < 2 {
			log.Println("invalid parts")
			return
		}

		go s.handleCommand(conn, parts)
	}
}

func (s *Server) handleCommand(conn net.Conn, parts []string) {
	switch commands.Command(parts[0]) {
	case commands.SET_COMMAND:
		msg, err := commands.ParseSetCommand(parts[1:])
		if err != nil {
			log.Println("parse set command error:", err)
			return
		}
		err = s.Opts.Cache.Set(msg.ItemKey(), msg.ItemValue(), msg.ItemTTL())
		if err != nil {
			log.Println("set command error:", err)
		}
	case commands.GET_COMMAND:
		msg, err := commands.ParseGetCommand(parts[1:])
		if err != nil {
			log.Println("parse get command error:", err)
			return
		}

		value, err := s.Opts.Cache.Get(msg.ItemKey())
		if err != nil {
			log.Println("get command error:", err)
			return
		}
		_, err = conn.Write(value)
		if err != nil {
			log.Println("write to conn error:", err)
		}
	case commands.DELETE_COMMAND:
		msg, err := commands.ParseDeleteCommand(parts[1:])
		if err != nil {
			log.Println("parse delete command error:", err)
			return
		}

		err = s.Opts.Cache.Delete(msg.ItemKey())
		if err != nil {
			log.Println("delete command error:", err)
			return
		}
	case commands.HAS_COMMAND:
		msg, err := commands.ParseHasCommand(parts[1:])
		if err != nil {
			log.Println("parse has command error:", err)
			return
		}

		exist := s.Opts.Cache.Has(msg.ItemKey())
		_, err = conn.Write([]byte(strconv.FormatBool(exist)))
		if err != nil {
			log.Println("write to conn error:", err)
		}
	default:
		log.Println(commands.ErrorInvalidCommand)
	}
}
