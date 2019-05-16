package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

func loadPrivateKey(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(key), nil
}

func runRemoteCmd(ip string, port int, privatekey ssh.AuthMethod, user string, cmd string, timeout int) (string, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			privatekey,
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Always accept key.
			return nil
		},
		Timeout: time.Duration(timeout) * time.Second,
	}

	client, err := ssh.Dial("tcp", ip+":"+strconv.Itoa(port), config)
	if err != nil {
		return "", err
	}

	defer client.Close()
	session, err := client.NewSession()

	if err != nil {
		return "", err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("vt100", 80, 40, modes); err != nil {
		return "", err
	}

	stdoutp, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}

	err = session.Run(cmd)
	if err != nil {
		return "", err
	}

	stdoutbuffer, err := ioutil.ReadAll(stdoutp)
	if err != nil {
		return "", err
	}
	return string(stdoutbuffer), nil
}

func main() {
	privateKey, err := loadPrivateKey("/.ssh/de-1")
	if err != nil {
		panic(err)
	}

	output, err := runRemoteCmd("de", 22, privateKey, "root", "cat /proc/cmdline", 10)
	if err != nil {
		panic(err)
	}
	fmt.Println("output:" + output)
}
