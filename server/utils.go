package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func getFilesFromClient(fileName string, conn net.Conn) error {
	defer conn.Close()
	// try n create a file on the server
	file, err := os.Create("cdbr." + fileName)
	if err != nil {
		return fmt.Errorf("[!] Impossible to create the file: %v", err)
	}
	defer file.Close()

	n, err := io.Copy(file, conn)
	if err != nil {
		return fmt.Errorf("[!] Something went wrong while copying the file: %v", err)
	}

	fmt.Printf("[*] Received %d bytes and saved to %s\n", n, "cdbr."+fileName)

	conn.Close()
	return nil
}

func sendFilesToClient(fileName string, conn net.Conn) error {
	// imagine sending a non existing file...
	file, err := os.Open(strings.TrimSpace(fileName))
	if err != nil {
		conn.Write([]byte("[!] The file does not exist..."))
		return fmt.Errorf("[!] Could not open the file")
	}
	defer file.Close()

	_, err = conn.Write([]byte(fileName))
	if err != nil {
		return fmt.Errorf("[!] Could not write the file name")
	}

	n, err := io.Copy(conn, file)
	if err != nil {
		return fmt.Errorf("[!] Error while copying file contents: %v", err)
	}

	fmt.Printf("[*] %d bytes sent", n) //should be a debug maybe???

	return nil
}

func trustClient(conn net.Conn) error {
	var ip string = conn.RemoteAddr().String()
	fmt.Printf("[?] Someone is trying to connect, their ip is: %s, do you want to accept the connection? (y/n)\n", ip)
	var resp string

	for {
		fmt.Scanln(&resp)
		switch resp {
		case "y":
			fmt.Println("[*] Proceeding with the connection...")
			return nil
		case "n":
			fmt.Println("[*] Closing the connection...")
			return fmt.Errorf("connection aborted by server")
		default:
			fmt.Println("[*] Please enter 'y' or 'n'")
		}
	}
}
