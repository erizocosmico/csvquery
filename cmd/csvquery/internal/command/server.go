package command

import (
	"fmt"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-mysql-server.v0/server"
	"gopkg.in/src-d/go-vitess.v1/mysql"
)

// Server starts a new MySQL server with the CSV files as a backend.
type Server struct {
	baseCmd
	User     string `short:"u" long:"user" default:"root" description:"User name to access the server."`
	Password string `short:"p" long:"password" description:"Password to access the server."`
	Port     int    `short:"P" long:"port" default:"3306" description:"Port in which the server will listen."`
	Host     string `short:"h" long:"host" default:"127.0.0.1" description:"Host name of the server."`
}

// Execute the command.
func (c *Server) Execute([]string) error {
	engine, _, err := c.engine()
	if err != nil {
		return err
	}

	auth := mysql.NewAuthServerStatic()
	auth.Entries[c.User] = []*mysql.AuthServerStaticEntry{
		{Password: c.Password},
	}

	addr := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	s, err := server.NewServer(
		server.Config{
			Protocol: "tcp",
			Address:  addr,
			Auth:     auth,
		},
		engine,
		server.DefaultSessionBuilder,
	)
	if err != nil {
		return fmt.Errorf("csvquery: unable to create server: %s", err)
	}

	logrus.Infof("server started and listening on %s:%d", c.Host, c.Port)
	return s.Start()
}
