package main

import (
	"flag"
	"strings"
	"time"

	"github.com/Xe/olegdb-ee/goleg"
	"github.com/tidwall/redcon"
)

var (
	port    = flag.String("port", "6660", "port that olegdb should listen on")
	sanity  = flag.Bool("sanity", false, "let the madness enter your being?")
	dataDir = flag.String("data", "./store", "where data should be stored")
)

func init() {
	flag.Parse()
}

func main() {
	db, err := goleg.Open(*dataDir, "primary", goleg.F_APPENDONLY|goleg.F_LZ4|goleg.F_SPLAYTREE)
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

		if result != nil {
			conn.WriteString(string(result))
		} else {
			conn.WriteError("ERR Could not unjar. Maybe key doesn't exist?")
		}

	case "scoop":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		key := string(cmd.Args[1])

		result := s.db.Scoop(key)

		if result == 0 {
			conn.WriteString(string(result))
		} else {
			conn.WriteError("ERR Could not scoop. Maybe key doesn't exist?")
		}

	case "mebbe":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		key := string(cmd.Args[1])

		success, keys := s.db.PrefixMatch(key)

		if success && keys != nil && len(keys) != 0 {
			conn.WriteArray(len(keys))

			for _, key := range keys {
				conn.WriteString(key)
			}
		} else {
			conn.WriteError("ERR prefix match returned no results")
		}

	case "dump":
		if len(cmd.Args) != 1 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		worked, keys := s.db.DumpKeys()
		if worked {
			conn.WriteArray(len(keys))

			for _, key := range keys {
				conn.WriteString(key)
			}
		} else {
			conn.WriteError("ERR Some error happened. Please try again.")
		}

	case "spoil":
		if len(cmd.Args) != 3 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		key := string(cmd.Args[1])
		dur, err := time.ParseDuration(string(cmd.Args[2]))
		if err != nil {
			conn.WriteError("ERR " + err.Error())
			return
		}

		then := time.Now().Add(dur)

		ret := s.db.Spoil(key, then)

		if ret == 0 {
			conn.WriteString("OK")
		} else {
			conn.WriteError("ERR Could not spoil. Maybe key doesn't exist?")
		}

	case "sniff":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		key := string(cmd.Args[1])

		expTime, exists := s.db.Expiration(key)
		if exists {
			showTime := expTime.Format(time.RFC3339)
			conn.WriteString(showTime)
		} else {
			conn.WriteError("ERR Could not sniff. Maybe key doesn't exist?")
		}

	case "squish":
		worked := s.db.Squish()
		if worked {
			conn.WriteString("OK")
		} else {
			conn.WriteError("ERR Shit's fucked fam")
		}

	case "canhas":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		key := string(cmd.Args[1])

		exists := s.db.Exists(key)

		if exists {
			conn.WriteInt(1)
		} else {
			conn.WriteInt(0)
		}

	case "uptime":
		uptime := s.db.Uptime() * (1e9)

		uptimeD := time.Duration(uptime)

		conn.WriteString(uptimeD.String())
	}
}
