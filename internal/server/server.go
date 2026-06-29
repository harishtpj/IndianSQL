package server

import (
	"fmt"
	"net"

	"github.com/go-mysql-org/go-mysql/server"
	"github.com/harishtpj/indiansql/internal/engine"
)

type Server struct {
	DBFile    string `yaml:"database"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	srv       *server.Server
	sqlEngine *engine.SQLEngine
}

func (s Server) getAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s *Server) Serve() error {
	addr := s.getAddr()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.srv = server.NewDefaultServer()
	s.sqlEngine, err = engine.NewSQLEngine(s.DBFile)
	if err != nil {
		return err
	}
	fmt.Println("[Server] Listening at", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		fmt.Printf("[Server] Client %s connected!\n", conn.LocalAddr())
		go s.handle(conn)
	}
}

func (s *Server) handle(c net.Conn) {
	conn, err := s.srv.NewConn(c, s.Username, s.Password, NewHandler(s.sqlEngine))
	if err != nil {
		fmt.Printf("[Server] Error: %v\n", err)
		c.Close()
		return
	}

	for {
		if err := conn.HandleCommand(); err != nil {
			fmt.Printf("[Server] Error: %v\n", err)
			return
		}
	}
}
