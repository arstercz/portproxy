package main

import (
	"io"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

type Proxy struct {
	bind, backend *net.TCPAddr
	sessionsCount int32
	pool          *recycler
}

func New(bind, backend string, size uint32) *Proxy {
	a1, err := net.ResolveTCPAddr("tcp", bind)
	if err != nil {
		log.Fatalln("resolve bind error:", err)
	}

	a2, err := net.ResolveTCPAddr("tcp", backend)
	if err != nil {
		log.Fatalln("resolve backend error:", err)
	}

	return &Proxy{
		bind:          a1,
		backend:       a2,
		sessionsCount: 0,
		pool:          NewRecycler(size),
	}
}

func (t *Proxy) pipe(dst, src *Conn, c chan int64, tag string) {
	defer func() {
		dst.CloseWrite()
		dst.CloseRead()
	}()
	if strings.EqualFold(tag, "send") {
		proxyLog(src, dst)
		c <- 0
	} else {
		n, err := io.Copy(dst, src)
		if err != nil {
			log.Print(err)
		}
		c <- n
	}
}

func (t *Proxy) transport(conn net.Conn) {
	start := time.Now()
	conn2, err := net.DialTCP("tcp", nil, t.backend)
	if err != nil {
		log.Print(err)
		return
	}
	connectTime := time.Now().Sub(start)
	log.Printf("proxy: %s ==> %s", conn2.LocalAddr().String(),
		conn2.RemoteAddr().String())
	start = time.Now()
	readChan := make(chan int64)
	writeChan := make(chan int64)
	var readBytes, writeBytes int64

	atomic.AddInt32(&t.sessionsCount, 1)
	var bindConn, backendConn *Conn
	bindConn = NewConn(conn, t.pool)
	backendConn = NewConn(conn2, t.pool)

	go t.pipe(backendConn, bindConn, writeChan, "send")
	go t.pipe(bindConn, backendConn, readChan, "receive")

	readBytes = <-readChan
	writeBytes = <-writeChan
	transferTime := time.Now().Sub(start)
	log.Printf("r: %d w:%d ct:%.3f t:%.3f [#%d]", readBytes, writeBytes,
		connectTime.Seconds(), transferTime.Seconds(), t.sessionsCount)
	atomic.AddInt32(&t.sessionsCount, -1)
}

func (t *Proxy) Start() {
	ln, err := net.ListenTCP("tcp", t.bind)
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		log.Printf("client: %s ==> %s", conn.RemoteAddr().String(),
			conn.LocalAddr().String())
		go t.transport(conn)
	}
}
