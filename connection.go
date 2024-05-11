package main

import (
	"SeewoMitM/internal/log"
	"fmt"
	"github.com/gorilla/websocket"
)

type Connection struct {
	URL            string
	UpstreamConn   *websocket.Conn
	DownstreamConn *websocket.Conn
}

var connectionPool = make([]*Connection, 20)

func GetConnectionPool() *[]*Connection {
	//去除nil元素
	//我也不知道为什么里面会有nil元素，但去掉了就不会有问题了
	var temp []*Connection
	for _, v := range connectionPool {
		if v != nil {
			temp = append(temp, v)
		}
	}
	return &temp
}

func GetConnectionPoolSize() int {
	return len(connectionPool)
}

func AddConnection(c *Connection) {
	//真的会有nil吗？
	if c == nil {
		return
	}
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
