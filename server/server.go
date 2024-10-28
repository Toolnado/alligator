package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/Toolnado/alligator/cache/interfaces"
	"github.com/Toolnado/alligator/commands"
)

type Options struct {
	Addr       string
	LeaderAddr string
	Leader     bool
}

type Server struct {
	Opts      Options
	Cache     interfaces.Cacher
	mu        sync.Mutex
	followers map[net.Conn]struct{}
}

func New(ops Options, cache interfaces.Cacher) *Server {
	return &Server{
		Opts:  ops,
		Cache: cache,
		mu:    sync.Mutex{},
		followers: func() map[net.Conn]struct{} {
			if ops.Leader {
				return make(map[net.Conn]struct{})
			}
			return nil
		}(),
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Opts.Addr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	if !s.Opts.Leader {
		go func() {
			conn, err := net.Dial("tcp", s.Opts.LeaderAddr)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("connected with leader:", s.Opts.LeaderAddr)
			conn.Write([]byte(commands.JoinCommand))
			s.handleConn(conn)
		}()
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
	log.Println("connection made:", conn.RemoteAddr())
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("server.handleConn error:", err)
		}
	}()

	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("server.handleConn error:", err)
			return
		}

		bufMSG := buf[:n]
		go s.handleCommand(conn, bufMSG)
	}
}

func (s *Server) handleCommand(conn net.Conn, raw []byte) {
	cmd := commands.New(raw)
	msg, err := cmd.Parse()
	if err != nil {
		log.Println(err)
		if _, connWriteErr := conn.Write([]byte(err.Error() + "\n")); connWriteErr != nil {
			log.Println(err, connWriteErr)
		}
		return
	}

	switch msg.Command() {
	case commands.SetCommand:
		err = s.handleSetCommand(msg)
	case commands.GetCommand:
		err = s.handleGetCommand(conn, msg)
	case commands.DeleteCommand:
		err = s.handleDeleteCommand(msg)
	case commands.HasCommand:
		err = s.handleHasCommand(conn, msg)
	case commands.JoinCommand:
		err = s.handleJoinCommand(conn)
	}

	if err != nil {
		log.Println(err)
		if _, connWriteErr := conn.Write([]byte(err.Error() + "\n")); connWriteErr != nil {
			log.Println(err, connWriteErr)
		}
	}
}

func (s *Server) sendUpdateToFollowers(msg commands.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for conn := range s.followers {
		if _, err := conn.Write(msg.Bytes()); err != nil {
			log.Println("server.sendUpdateToFollowers error:", err)
		}
	}
	return nil
}

func (s *Server) handleSetCommand(msg commands.Message) error {
	if err := s.Cache.Set(msg.Key(), msg.Value(), msg.TTL()); err != nil {
		return fmt.Errorf("server.handleSetCommand error: %s", err)
	}
	if s.Opts.Leader {
		if err := s.sendUpdateToFollowers(msg); err != nil {
			return fmt.Errorf("server.handleSetCommand error: %s", err)
		}
	}

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, msg commands.Message) error {
	value, err := s.Cache.Get(msg.Key())
	if err != nil {
		return fmt.Errorf("server.handleGetCommand error: %s", err)
	}
	if _, err = conn.Write(value); err != nil {
		return fmt.Errorf("server.handleGetCommand error: %s", err)
	}
	return nil
}

func (s *Server) handleDeleteCommand(msg commands.Message) error {
	if err := s.Cache.Delete(msg.Key()); err != nil {
		return fmt.Errorf("server.handleDeleteCommand error: %s", err)
	}
	return nil
}

func (s *Server) handleHasCommand(conn net.Conn, msg commands.Message) error {
	exist := s.Cache.Has(msg.Key())
	if _, err := conn.Write([]byte(strconv.FormatBool(exist))); err != nil {
		return fmt.Errorf("server.handleHasCommand error: %s", err)
	}
	return nil
}

func (s *Server) handleJoinCommand(conn net.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.followers[conn] = struct{}{}
	return nil
}
