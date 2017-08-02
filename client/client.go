package main

import "os"
import "net"
import "log"
import "io/ioutil"
import "crypto/sha256"
import "bytes"
import "fmt"
import "bufio"

func runCmd(ip, cmd string) string {
	pwHash, err := ioutil.ReadFile("passwdHash")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}
	conn.Write(append(pwHash, []byte(cmd)...))
	var buf bytes.Buffer
	buf.ReadFrom(conn)
	return buf.String()
}
func main() {
	if os.Args[1] == "pwgen" {
		x := sha256.Sum256([]byte(os.Args[2]))
		ioutil.WriteFile("passwdHash",
			x[:],
			0600)
	} else if os.Args[1] == "console" {
		for {
			read := bufio.NewReader(os.Stdin)
			cmd, _ := read.ReadString('\n')
			fmt.Println(runCmd(os.Args[2], cmd))
		}
	} else if os.Args[1] == "cmd" {
		fmt.Println(runCmd(os.Args[2], os.Args[3]))
	}
}
