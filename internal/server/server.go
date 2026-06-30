package server

import (
	"fmt"
	"net"

	"github.com/go-mysql-org/go-mysql/server"
	"github.com/harishtpj/indiansql/internal/engine"
)

type Server struct {
	DBFile   string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	srv      *server.Server
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
	fmt.Println("[Server] Listening at", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		fmt.Printf("[Server] Client %s connected!\n", conn.RemoteAddr().String())
		go s.handle(conn)
	}
}

func (s *Server) handle(c net.Conn) {
	sqlEngine, err := engine.NewSQLEngine(s.DBFile)
	if err != nil {
		fmt.Printf("[Server/Engine] Error: %v\n", err)
		c.Close()
		return
	}

	conn, err := s.srv.NewConn(c, s.Username, s.Password, NewHandler(sqlEngine))
	if err != nil {
		fmt.Printf("[Server] Error: %v\n", err)
		sqlEngine.Close()
		c.Close()
		return
	}

	for {
		if err := conn.HandleCommand(); err != nil {
			fmt.Printf("[Server] Error: %v\n", err)
			sqlEngine.Close()
			return
		}
	}
}
