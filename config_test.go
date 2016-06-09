package goslim

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	as := assert.New(t)
	c := &config{
		Mode:       "DEV",
		ListenAddr: "localhost",
		ListenPort: "8080",
		PubDir:     "public",
	}
	c.loadConfig("testfiles/config_test.yaml")

	as.Equal("TEST", c.Mode)
	as.Equal("localhost", c.ListenAddr)
	as.Equal("8080", c.ListenPort)
	as.Equal("static", c.PubDir)
	as.Equal("MySQL", c.Db.SQLDriver)
	as.Equal("tcp", c.Db.Protocol)
	as.Equal("localhost", c.Db.Hostname)
	as.Equal("root", c.Db.Username)
	as.Equal("123456", c.Db.Password)
	as.Equal("test", c.Db.Dbname)
	as.Equal("utf8", c.Db.Charset)

	// mode: TEST
	// listenaddr: localhost
	// listenport: 8080
	// pubdir: static
	// db:
	//     sqldriver: "MySQL"
	//     protocal: "tcp"
	//     hostname: "localhost"
	//     username: "root"
	//     password: "123456"
	//     dbname: "testdb"
	//     charset: "utf-8"
}
