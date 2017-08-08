package main

import (
	"io"
	"log"
	"net"
	"strings"
	"bytes"
	"sync/atomic"
	"time"
	"strconv"
	"fmt"
	"errors"
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

func readInitPacket(server net.Conn, client net.Conn) (c ServerConn, err error) {
	buffer := make([]byte, Bsize)
	n, err := server.Read(buffer)
	data := make([]byte, n)
	copy(data, buffer[0:n])
	//fmt.Printf("1 server\n")
	//for _, v1 := range buffer[0:n] {
	//	fmt.Printf("%x ", v1)
	//}
	//fmt.Printf("\n")
	if err != nil {
		log.Printf("read init packet from server error: %s\n", err.Error())
		server.Close()
		return c, err
	}
	pos := 0
	//protocol
	pos += 4
	c.protocol = uint16(buffer[pos])
	//skip version string
	pos += bytes.IndexByte(buffer[pos:], 0)
	//skip connection id
	pos += 4
	log.Printf("pos: %d\n", pos)
	//salt
	c.salt = buffer[pos + 1 : pos + 1 + 8]
	pos += 9
	//skip filler
	pos++
	//skip capability
	pos += 2
	//skip charset
	pos++
	//skip status flag
	pos += 2
	//skil capability
	pos += 2
	//skip length
	authLen := int(buffer[pos])
	pos++
	//log.Printf("authLen server: %d\n", authLen)
	//skip reserved
	pos += 10
	reLen := authLen - 8
	if reLen < 13 {
		reLen = 13
	}
	c.salt = append(c.salt, buffer[pos : pos  + reLen - 1]...)
	//fmt.Printf("server init:\n")
	//for _, vc := range c.salt {
	//	fmt.Printf("%x ", vc)
	//}
	//fmt.Printf("\n")
	//fmt.Printf("2 server\n")
	//for _, v2 := range data {
	//	fmt.Printf("%x ", v2)
	//}
	//fmt.Printf("\n")
	_, _ = client.Write(data[0:n])
	return c, nil
}

func passwordCheck(salt []byte, user string, auth []byte) error {
	// check otp password
	secret, err := userSecret(Dbh, user)
	if err != nil {
		log.Printf("cannot get secret with user: %s\n", user)
		return errors.New(fmt.Sprintf("cannot get secret with user: %s\n", user))
	}
	totp, err := getOtpPass(secret)
	//log.Printf("get totp with secret: %s, salt: %s\n", secret, string(salt))
	if err != nil {
		log.Printf("cannot get totp with secret: %s\n", secret)
		return errors.New(fmt.Sprintf("cannot get totp with secret: %s\n", secret))
	}
	userPass  := user + strconv.FormatUint(uint64(uint32(totp)), 10)
	checkAuth := calcPassword(salt, []byte(userPass))

	if !bytes.Equal(auth, checkAuth) {
		log.Printf("mismatch password %s for user %s\n", userPass, user)
		log.Printf("check auth: %s, old auth: %s\n", string(checkAuth), auth)
		return errors.New(fmt.Sprintf("mismatch password '%s' for user %s\n", userPass, user))
	}
	return nil
}

func readInitResponse(client net.Conn, server net.Conn, salt []byte) (c ClientConn, err error) {
	buffer := make([]byte, Bsize)
	n, err := client.Read(buffer)
	if err != nil {
		log.Printf("read init response from client error: %s\n", err.Error())
		return c, errors.New(fmt.Sprintf("read init response from client error: %s\n", err.Error()))
	}
	data := make([]byte, 0)
	pos := 0
	//packet header
	pos += 4
	//capability
	pos += 4
	//skip max packet size
	pos += 4
	//charset skip
	pos++
	//skip reserved 23[00]
	pos += 23
	header1 := pos
	//user name
	c.user = string(buffer[pos : pos+bytes.IndexByte(buffer[pos:], 0)])
	pos += len(c.user) + 1

	//auth length and auth
	authLen := int(buffer[pos])
	pos++
	c.auth = buffer[pos : pos+authLen]
	//log.Printf("auth: %s\n", c.auth)
	pos += authLen

	err = passwordCheck(salt, c.user, c.auth)
	if err != nil {
		log.Printf("error: %s", err)
		return c, err
	}

	// make a new response packate
	data = append(data, buffer[0 : header1]...)
	data = append(data, MysqlUser...)
	data = append(data, 0)
	authNew := calcPassword(salt, []byte(MysqlPass))
	authNewLen := len(authNew)
	data = append(data, byte(authNewLen))
	data = append(data, authNew...)
	data = append(data, buffer[pos : n]...)
	dataLen := len(data)
	data[0] = byte(dataLen - 4)

	// write to server conn
	_, _ = server.Write(data[0:dataLen])
	return c, nil
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

	//mysql handshake check.
	servert, err := readInitPacket(conn2, conn)
	if err != nil {
		log.Printf("read server error: %s\n", err)
		conn2.Write([]byte("read server init packet error\n"))
		conn2.Close()
		return
	}
	clientt, err := readInitResponse(conn, conn2, servert.salt)
	if err != nil {
		log.Printf("read client init error: %s\n", err)
		conn.Write([]byte(fmt.Sprintf("read client init response error: %s\n", err.Error())))
		conn2.Write([]byte(fmt.Sprintf("client init response error: %s\n", err.Error())))
		conn2.Close() // close the server conn
		conn.Close()
		return
	}

	log.Printf("server salt: %s\n", string(servert.salt))
	log.Printf("user name: %s\n", clientt.user)

	// proxy
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
