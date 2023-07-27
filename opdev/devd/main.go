package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	cmd := flag.String("c", "", "cmd:update,build,restart,deploy(run:update->build->restart altogether)")
	flag.Parse()

	if *cmd == "" {
		flag.Usage()
		return
	}

	dials := os.Getenv("DLOG_DIALS")
	for _, dial := range strings.Split(dials, ",") {
		func() {
			components := strings.Split(dial, "@")

			if len(components) != 2 {
				return
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

			// 建立新会话
			session, err := sshCli.NewSession()
			if err != nil {
				log.Fatalf("new session error: %s", err.Error())
			}

			defer sshCli.Close()

			switch *cmd {
			case "update":
				var b bytes.Buffer
				session.Stdout = &b
				remoteCmd := fmt.Sprintf(". ~/.bash_profile;cd ${HOME}/trunk;sh ${HOME}/trunk/sh/update.sh")
				fmt.Println("run:" + remoteCmd + " from " + components[1])
				if err := session.Run(remoteCmd); err != nil {
					fmt.Println("Failed to run: " + remoteCmd + " from " + components[1] + ", error:" + err.Error())
					return
				}
			case "build":
				var b bytes.Buffer
				session.Stdout = &b
				remoteCmd := fmt.Sprintf(". ~/.bash_profile;cd ${HOME}/trunk;sh ./sh/gamesrv_b.sh")
				fmt.Println("run:" + remoteCmd + " from " + components[1])
				if err := session.Run(remoteCmd); err != nil {
					fmt.Println("Failed to run: " + remoteCmd + " from " + components[1] + ", error:" + err.Error())
					return
				}
				session.Close()

				// 建立新会话
				session, err = sshCli.NewSession()
				if err != nil {
					log.Fatalf("new session error: %s", err.Error())
				}
				session.Stdout = &bytes.Buffer{}
				remoteCmd = fmt.Sprintf(". ~/.bash_profile;cd ${HOME}/trunk;sh ./sh/fightsrv_b.sh")
				fmt.Println("run:" + remoteCmd + " from " + components[1])
				if err := session.Run(remoteCmd); err != nil {
					fmt.Println("Failed to run: " + remoteCmd + " from " + components[1] + ", error:" + err.Error())
					return
				}
			case "restart":
				var b bytes.Buffer
				session.Stdout = &b
				remoteCmd := fmt.Sprintf(". ~/.bash_profile;cd ${HOME}/trunk;sh ./sh/srv.sh restart")
				fmt.Println("run:" + remoteCmd + " from " + components[1])
				if err := session.Run(remoteCmd); err != nil {
					fmt.Println("Failed to run: " + remoteCmd + " from " + components[1] + ", error:" + err.Error())
					return
				}
			case "deploy":
				var b bytes.Buffer
				session.Stdout = &b
				remoteCmd := fmt.Sprintf(". ~/.bash_profile;cd ${HOME}/trunk;" +
					"sh ./sh/update.sh;" +
					"sh ./sh/gamesrv_b.sh;" +
					"sh ./sh/fightsrv_b.sh;" +
					"sh ./sh/srv.sh restart")
				fmt.Println("run:" + remoteCmd + " from " + components[1])
				if err := session.Run(remoteCmd); err != nil {
					fmt.Println("Failed to run: " + remoteCmd + " from " + components[1] + ", error:" + err.Error())
					return
				}
			}
		}()
	}
}
