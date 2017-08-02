package main

import "os"
import "net"
import "log"
import "io/ioutil"
import "crypto/sha256"
import "crypto/rsa"
import "crypto/rand"
import "crypto/x509"
import "bytes"
import "fmt"
import "bufio"

var pub *rsa.PublicKey

func encode(msg []byte) []byte {
	enc, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, msg, []byte("label"))
	retVal := make([]byte, 0)
	for {
		if len(enc) < 255 {
			retVal = append(retVal, byte(uint8(len(enc)+1))) //all this casting crap necessary?
			retVal = append(retVal, enc...)
			retVal = append(retVal, byte(0))
			break
		} else {
			retVal = append(retVal, byte(uint8(255))) //all this casting crap necessary?
			retVal = append(retVal, enc[:255]...)
			enc = enc[255:]
		}
	}
	return retVal
}
func runCmd(ip, cmd string) string {
	pwHash, err := ioutil.ReadFile("passwdHash")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}
	conn.Write(encode(append(pwHash, []byte(cmd)...)))
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
		pubFile, err := ioutil.ReadFile("rsaPub")
		if err != nil {
			return
		}
		x, err := x509.ParsePKIXPublicKey(pubFile)
		if err != nil {
			return
		}
		pub = x.(*rsa.PublicKey)

		for {
			read := bufio.NewReader(os.Stdin)
			cmd, _ := read.ReadString('\n')
			fmt.Println(runCmd(os.Args[2], cmd))
		}
	} else if os.Args[1] == "cmd" {
		fmt.Println(runCmd(os.Args[2], os.Args[3]))
	}
}
