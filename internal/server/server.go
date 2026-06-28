package server

import (
	"fmt"
	"net"

	"github.com/go-mysql-org/go-mysql/server"
	"github.com/harishtpj/indiansql/internal/handler"
)

type Server struct {
	Addr string
}

func (s *Server) Serve(dbFile string) error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	fmt.Println("[Server] Listening at", s.Addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		fmt.Printf("[Server] Client %s connected!\n", conn.LocalAddr())
		go s.handle(conn, dbFile)
	}
}

func (s *Server) handle(c net.Conn, dbFile string) {
	db, err := handler.NewREPLContext(dbFile)
	if err != nil {
		fmt.Printf("[Server/ReplCtx] Error: %v\n", err)
	}

	srv := server.NewDefaultServer()
	conn, err := srv.NewConn(c, "root", "", &Handler{db})
	if err != nil {
		fmt.Printf("[Server] Error: %v\n", err)
	}

	for {
		if err := conn.HandleCommand(); err != nil {
			fmt.Printf("[Server] Error: %v\n", err)
		}
	}
}
