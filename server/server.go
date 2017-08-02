package main

import "net"
import "os/exec"
import "crypto/sha256"
import "io/ioutil"
import "os"
import "bytes"

var pwHash [sha256.Size]byte

func handleconn(conn net.Conn) {
	b := make([]byte, 232)
	conn.Read(b)
	if sha256.Sum256(b[:32]) == pwHash {
		end := 32
		for end < len(b) && b[end] != byte(0) {
			end++
		}
		cmd := exec.Command("cmd", "/C", string(b[32:end]))
		var output bytes.Buffer
		cmd.Stdout = &output
		err := cmd.Run()
		if err == nil {
			conn.Write(output.Bytes())
		} else {
			conn.Write([]byte("Error")])
		}
	}
	conn.Close()
}
func main() {
	if len(os.Args) < 2 || os.Args[1] != "pwgen" {
		pwHashSlice, err := ioutil.ReadFile("passwdHash")
		if err != nil {
			return
		}
		copy(pwHash[:], pwHashSlice)
		ln, err2 := net.Listen("tcp", ":3924")
		if err2 != nil {
			return
		}
		for {
			conn, _ := ln.Accept()
			go handleconn(conn)
		}
	} else {
		x := sha256.Sum256([]byte(os.Args[2]))
		x = sha256.Sum256(x[:])
		ioutil.WriteFile("passwdHash",
			x[:],
			0600)
	}
}
