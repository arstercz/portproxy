package main

import (
	"log"
)

//read more client-server protocol from http://dev.mysql.com/doc/internals/en/text-protocol.html
const (
	comQuit byte = iota + 1
	comInitDB
	comQuery
	comFieldList
	comCreateDB
	comDropDB
	comRefresh
	comShutdown
	comStatistics
	comProcessInfo
	comConnect
	comProcessKill
	comDebug
	comPing
	comTime
	comDelayedInsert
	comChangeUser
	comBinlogDump
	comTableDump
	comConnectOut
	comRegiserSlave
	comStmtPrepare
	comStmtExecute
	comStmtSendLongData
	comStmtClose
	comStmtReset
	comSetOption
	comStmtFetch
)

func proxyLog(src, dst *Conn) {
	buffer := make([]byte, Bsize)
	clientIp := src.conn.RemoteAddr().String()
	serverIp := dst.conn.RemoteAddr().String()
	for {
		n, err := src.Read(buffer)
		if err != nil {
			return
		}
		if n >= 5 {
			switch buffer[4] {
			case comQuit:
				log.Printf("From %s To %s; Quit: %s\n", clientIp, serverIp, "user quit")
			case comInitDB:
				log.Printf("From %s To %s; schema: use %s\n", clientIp, serverIp, string(buffer[5:n]))
			case comQuery:
				log.Printf("From %s To %s; Query: %s\n", clientIp, serverIp, string(buffer[5:n]))
			case comFieldList:
				log.Printf("From %s To %s; Table columns list: %s\n", clientIp, serverIp, string(buffer[5:n]))
			case comConnect:
				log.Printf("Internal: internal command in the server\n")
			case comRefresh:
				log.Printf("From %s To %s; Refresh: command: %s\n", clientIp, serverIp, string(buffer[5:n]))
			case comStmtPrepare:
				log.Printf("From %s To %s; Prepare Query: %s\n", clientIp, serverIp, string(buffer[5:n]))
			case comStmtExecute:
				log.Printf("From %s To %s; Prepare Args: %s\n", clientIp, serverIp, string(buffer[5:n]))
			case comProcessKill:
				log.Printf("From %s To %s; Kill: kill conntion %s\n", clientIp, serverIp, string(buffer[5:n]))
			}
		}

		_, err = dst.Write(buffer[0:n])
		if err != nil {
			return
		}
	}
}
