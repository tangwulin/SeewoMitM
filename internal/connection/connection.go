package connection

import (
	"SeewoMitM/internal/log"
	"fmt"
	"github.com/gorilla/websocket"
)

type Connection struct {
	URL        string
	Upstream   *websocket.Conn
	Downstream *websocket.Conn
}

var connectionPool = make([]*Connection, 20)

func GetConnectionPool() *[]*Connection {
	return &connectionPool
}

func GetConnectionPoolSize() int {
	return len(connectionPool)
}

func AddConnection(c *Connection) {
	connectionPool = append(connectionPool, c)
	log.WithFields(log.Fields{"type": "AddConnection"}).Info(fmt.Sprintf("Connection added, URL:%s", c.URL))
}

func RemoveConnection(c *Connection) {
	for i, v := range connectionPool {
		if v == c {
			connectionPool = append(connectionPool[:i], connectionPool[i+1:]...)
		}
	}
	log.WithFields(log.Fields{"type": "RemoveConnection"}).Info(fmt.Sprintf("Connection removed, URL:%s", c.URL))
}
