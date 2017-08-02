package main

import "net"
import "os"
import "os/exec"
import "crypto/sha256"
import "crypto/rsa"
import "crypto/x509"
import "crypto/rand"
import "io/ioutil"
import "bytes"
import "fmt"

var pwHash [sha256.Size]byte
var pub *rsa.PublicKey
var priv *rsa.PrivateKey

func downloadData(conn net.Conn) []byte {
	var buf, data []byte
	for {
		buf = make([]byte, 1)
		conn.Read(buf)
		amtToGet := uint8(buf[0])
		if amtToGet == 0 {
			break
		}
		buf = make([]byte, amtToGet)
		conn.Read(buf)
		data = append(data, buf...)
	}
	return data
}
func handleconn(conn net.Conn) {
	b := downloadData(conn)
	b, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, b, []byte("label"))
	if err != nil || len(b) < 32 {
		return
	}
	if sha256.Sum256(b[:32]) == pwHash {
		/*end := 32
		for end < len(b) && b[end] != byte(0) {
			end++
		}*/
		cmd := exec.Command("cmd", "/C", string(b[32: /*end*/]))
		var output bytes.Buffer
		cmd.Stdout = &output
		err := cmd.Run()
		if err == nil {
			conn.Write(output.Bytes())
		} else {
			conn.Write([]byte("Error"))
		}
	}
	conn.Close()
}
func main() {
	if len(os.Args) > 1 && os.Args[1] == "pwgen" {
		x := sha256.Sum256([]byte(os.Args[2]))
		x = sha256.Sum256(x[:])
		ioutil.WriteFile("passwdHash",
			x[:],
			0600)
		fmt.Println("Success")
	} else if len(os.Args) > 1 && os.Args[1] == "privgen" {
		priv, _ = rsa.GenerateKey(rand.Reader, 2048)
		ioutil.WriteFile("rsaPriv",
			x509.MarshalPKCS1PrivateKey(priv),
			0600)
		fmt.Println("Success")
	} else if len(os.Args) > 1 && os.Args[1] == "pubgen" {
		privFile, err := ioutil.ReadFile("rsaPriv")
		if err != nil {
			return
		}
		priv, err = x509.ParsePKCS1PrivateKey(privFile)
		if err != nil {
			return
		}

		x, err := x509.MarshalPKIXPublicKey(priv.Public())
		ioutil.WriteFile("rsaPub",
			x,
			0600)
		fmt.Println("Success")
	} else {
		pwHashSlice, err := ioutil.ReadFile("passwdHash")
		if err != nil {
			return
		}
		copy(pwHash[:], pwHashSlice)

		privFile, err := ioutil.ReadFile("rsaPriv")
		if err != nil {
			return
		}
		priv, err = x509.ParsePKCS1PrivateKey(privFile)
		if err != nil {
			return
		}

		p := priv.Public()
		pub = p.(*rsa.PublicKey)

		ln, err := net.Listen("tcp", ":3924")
		if err != nil {
			return
		}
		for {
			conn, _ := ln.Accept()
			go handleconn(conn)
		}
	}
}
