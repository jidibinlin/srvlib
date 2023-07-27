package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

// Please set backup log host ENV BACKUP_LOG_HOST=user:password@host:port,user2:password2@host2:port2
func main() {
	var ipAddr = strings.ReplaceAll(LocalIp(), ".", "_")
	for {
		dir := os.Getenv("TLOGDIR")
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Fatalln(err)
		}

		var backupLogs []os.DirEntry
		for _, file := range files {
			if !strings.Contains(file.Name(), time.Now().Format("01-02")) {
				backupLogs = append(backupLogs, file)
			}
		}

		sshCli, sftpCli := connect()
		// 建立新会话
		session, err := sshCli.NewSession()
		if err != nil {
			panic("new session error: %s" + err.Error())
		}

		session.Run("mkdir /tmp/logmigrate")
		session.Close()

		for _, file := range backupLogs {
			fp, err := os.Open(path.Join(dir, file.Name()))
			if nil != err {
				panic(err)
			}

			ftpFile, err := sftpCli.Create("/tmp/logmigrate/" + ipAddr + "_" + file.Name())
			if nil != err {
				panic(err)
			}

			fileByte, err := io.ReadAll(fp)
			if nil != err {
				panic(err)
			}

			ftpFile.Write(fileByte)

			fp.Close()
		}

		session, err = sshCli.NewSession()
		if err != nil {
			panic("new session error:" + err.Error())
		}
		if err = session.Run(". ~/.bash_profile; mv /tmp/logmigrate/* ${TLOGDIR}/"); err != nil {
			panic("new session error:" + err.Error())
		}
		session.Close()

		sshCli.Close()
		sftpCli.Close()

		for _, file := range backupLogs {
			if err = os.Remove(dir + file.Name()); err != nil {
				fmt.Println(err)
			}
		}

		time.Sleep(time.Hour)
	}
}

func connect() (*ssh.Client, *sftp.Client) {
	backupHost := os.Getenv("BACKUP_LOG_HOST")

	components := strings.Split(backupHost, "@")

	if len(components) != 2 {
		log.Fatalln("BACKUP_LOG_HOST ENV INCORRECT, FORMAT:user:password@host:port")
	}

	userPasswd := strings.Split(components[0], ":")

	sshCli, err := ssh.Dial("tcp", components[1], &ssh.ClientConfig{
		User: userPasswd[0],
		Auth: []ssh.AuthMethod{ssh.Password(userPasswd[1])},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		log.Fatalln("error occurred:", err)
	}

	sftpCli, err := sftp.NewClient(sshCli)
	if err != nil {
		log.Fatalln("new sftp client error: %w", err)
	}

	return sshCli, sftpCli
}

func LocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic("ip address is " + err.Error())
	}
	var ip = "localhost"
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ip
}
