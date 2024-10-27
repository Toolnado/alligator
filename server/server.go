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
		if err := s.handleSetCommand(cmd); err != nil {
			// TODO:respond
			return
		}
	case commands.GET_COMMAND:
		if err := s.handleGetCommand(conn, cmd); err != nil {
			// TODO:respond
			return
		}
	case commands.DELETE_COMMAND:
		if err := s.handleDeleteCommand(cmd); err != nil {
			// TODO:respond
			return
		}
	case commands.HAS_COMMAND:
		if err := s.handleHasCommand(conn, cmd); err != nil {
			// TODO:respond
			return
		}
	}
}

func (s *Server) handleSetCommand(cmd commands.CMD) error {
	if err := s.Opts.Cache.Set(cmd.Key, cmd.Value, cmd.TTL); err != nil {
		return fmt.Errorf("server.handleSetCommand error: %s", err)
	}
	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, cmd commands.CMD) error {
	if value, err := s.Opts.Cache.Get(cmd.Key); err != nil {
		return fmt.Errorf("server.handleGetCommand error: %s", err)
	} else {
		if _, err = conn.Write(value); err != nil {
			return fmt.Errorf("server.handleGetCommand error: %s", err)
		}
	}
	return nil
}

func (s *Server) handleDeleteCommand(cmd commands.CMD) error {
	if err := s.Opts.Cache.Delete(cmd.Key); err != nil {
		return fmt.Errorf("server.handleDeleteCommand error: %s", err)
	}
	return nil
}

func (s *Server) handleHasCommand(conn net.Conn, cmd commands.CMD) error {
	exist := s.Opts.Cache.Has(cmd.Key)
	if _, err := conn.Write([]byte(strconv.FormatBool(exist))); err != nil {
		return fmt.Errorf("server.handleHasCommand error: %s", err)
	}
	return nil
}
