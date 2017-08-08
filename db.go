/*config read to verify normal user*/
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"crypto/sha1"
)

func dbh(dsn string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return db, err
	}
	return db, nil
}

func Query(db *sql.DB, q string) (*sql.Rows, error) {
	if Verbose {
		log.Printf("Query: %s\n", q)
	}
	return db.Query(q)
}

func QueryRow(db *sql.DB, q string) *sql.Row {
	if Verbose {
		log.Printf("Query: %s", q)
	}
	return db.QueryRow(q)
}

func ExecQuery(db *sql.DB, q string) (sql.Result, error) {
	if Verbose {
		log.Printf("ExecQuery: %s\n", q)
	}
	return db.Exec(q)
}

func insertlog(db *sql.DB, t *query) bool {
	insertSql := `
	insert into query_log(bindport, client, client_port, server, server_port, sql_type, 
	sql_string, create_time) values (%d, '%s', %d, '%s', %d, '%s', '%s', now())
	`
	_, err := ExecQuery(db, fmt.Sprintf(insertSql, t.bindPort, t.client, t.cport, t.server, t.sport, t.sqlType, t.sqlString))
	if err != nil {
		return false
	}
	return true
}

func userSecret(db *sql.DB, user string) (secret string, err error) {
	selectSQL := `select secret from otp_secret where name = '%s'`
	err = QueryRow(db, fmt.Sprintf(selectSQL, user)).Scan(&secret)
	if err != nil {
		return "", err
	}
	return secret, nil
}

// calculate mysql password
func calcPassword(scramble, password []byte) []byte {
	if len(password) == 0 {
		return nil
	}

	// stageHash = SHA1(password)
	crypt := sha1.New()
	crypt.Write(password)
	stage1 := crypt.Sum(nil)

	// scrambleHash = SHA1(scramble + SHA1(stage1Hash))
	// inner Hash
	crypt.Reset()
	crypt.Write(stage1)
	hash := crypt.Sum(nil)

	// outer Hash
	crypt.Reset()
	crypt.Write(scramble)
	crypt.Write(hash)
	scramble = crypt.Sum(nil)

	// token = scrambleHash XOR stageHash
	for i, _ := range scramble {
		scramble[i] ^= stage1[i]
	}
	return scramble
}
