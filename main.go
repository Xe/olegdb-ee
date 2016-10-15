package main

import (
	"flag"
	"strings"

	"github.com/Xe/olegdb-ee/goleg"
	"github.com/tidwall/redcon"
)

var (
	port   = flag.String("port", "6660", "port that olegdb should listen on")
	sanity = flag.Bool("sanity", false, "let the madness enter your being?")
)

func init() {
	flag.Parse()
}

func main() {
	db, err := goleg.Open("./test.val", "test", 0)
	if err != nil {
		panic(err)
	}

	s := &Server{
		db: db,
	}

	err = redcon.ListenAndServe(":"+*port, s.HandleCommand, nil, nil)
	if err != nil {
		panic(err)
	}
}

type Server struct {
	db goleg.Database
}

func (s *Server) HandleCommand(conn redcon.Conn, cmd redcon.Command) {
	switch strings.ToLower(string(cmd.Args[0])) {
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	case "ping":
		conn.WriteString("PONG")
	case "quit":
		conn.WriteString("OK")
		conn.Close()
	case "jar":
		if len(cmd.Args) != 3 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		key := cmd.Args[1]
		value := cmd.Args[2]
		result := s.db.Jar(string(key), value)

		if result == 0 {
			conn.WriteString("OK")
			return
		} else {
			conn.WriteError("ERR Couldn't save to oleg :(")
		}
	case "unjar":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		key := string(cmd.Args[1])

		result := s.db.Unjar(key)
		if len(result) != 0 {
			conn.WriteString(string(result))
			conn.WriteString("OK")
		}
	}
}
