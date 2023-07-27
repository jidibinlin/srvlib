package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Please set queriable log host ENV DLOG_DIALS=user:password@host:port,user2:password2@host2:port2
func main() {
	timeSeq := flag.String("t", time.Now().Format("0102"), "日志日期,如0601")
	keyword := flag.String("k", "", "关键字")
	flag.Parse()

	if *keyword == "" {
		flag.Usage()
		return
	}

	var fNameSuffix string
	switch true {
	case len(*timeSeq) == 1:
		fNameSuffix = time.Now().Format("01-0") + *timeSeq
	case len(*timeSeq) == 2:
		fNameSuffix = time.Now().Format("01-") + *timeSeq
	case len(*timeSeq) == 3:
		fNameSuffix = fmt.Sprintf("0%s-%s", (*timeSeq)[0:1], (*timeSeq)[1:])
	case len(*timeSeq) == 4:
		fNameSuffix = fmt.Sprintf("%s-%s", (*timeSeq)[0:2], (*timeSeq)[2:4])
	default:
		fNameSuffix = time.Now().Format("01-02")
	}
	fNameSuffix += ".log"

	tmpFp, err := os.CreateTemp("", "dlog")
	if err != nil {
		panic(err)
	}

	fmt.Println("created cache file:" + tmpFp.Name())

	dials := os.Getenv("DLOG_DIALS")
	fmt.Println("search log host:", dials)
	var (
		content         string
		appendContentMu sync.Mutex
		wg              sync.WaitGroup
	)
	for _, dial := range strings.Split(dials, ",") {
		wg.Add(1)
		go func(dial string) {
			defer wg.Done()

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

			defer sshCli.Close()

			// 建立新会话
			session, err := sshCli.NewSession()
			if err != nil {
				log.Fatalf("new session error: %s", err.Error())
			}

			defer session.Close()

			var b bytes.Buffer
			session.Stdout = &b
			cmd := fmt.Sprintf(`. ~/.bash_profile; cat ${TLOGDIR}*.%s|grep "%s"`, fNameSuffix, *keyword)
			fmt.Println("run:" + cmd + " from " + components[1])
			if err := session.Run(cmd); err != nil {
				fmt.Println("Failed to run: " + cmd + " from " + components[1] + ", error:" + err.Error())
				return
			}

			ip := strings.Split(components[1], ":")[0]
			lines := strings.Split(b.String(), "\n")

			appendContentMu.Lock()
			for _, l := range lines {
				if len(l) == 0 {
					continue
				}
				content = content + fmt.Sprintln("["+ip+"]"+l)
			}
			appendContentMu.Unlock()

			// 建立新会话
			session, err = sshCli.NewSession()
			if err != nil {
				log.Fatalf("new session error: %s", err.Error())
			}

			session.Stdout = &bytes.Buffer{}
			cmd = fmt.Sprintf(`. ~/.bash_profile; cat ${TLOGDIR}*/*.%s|grep "%s"`, fNameSuffix, *keyword)
			fmt.Println("run:" + cmd + " from " + components[1])
			if err := session.Run(cmd); err != nil {
				fmt.Println("Failed to run: " + cmd + " from " + components[1] + ", error:" + err.Error())
				return
			}

			ip = strings.Split(components[1], ":")[0]
			lines = strings.Split(b.String(), "\n")

			appendContentMu.Lock()
			defer appendContentMu.Unlock()
			for _, l := range lines {
				if len(l) == 0 {
					continue
				}
				content = content + fmt.Sprintln("["+ip+"]"+l)
			}
		}(dial)
	}

	wg.Wait()

	_, err = tmpFp.Write([]byte(content))
	if err != nil {
		panic(err)
	}

	fmt.Println("cached file:" + tmpFp.Name())

	fmt.Println("cat", strings.ReplaceAll(tmpFp.Name(), "\\", "/")+"|sort -k4")

	catCmd := exec.Command("cat", strings.ReplaceAll(tmpFp.Name(), "\\", "/"))
	sortCmd := exec.Command("sort", "-k4")
	catCmd.Stderr = &bytes.Buffer{}
	sortCmd.Stderr = &bytes.Buffer{}
	sortCmd.Stdout = &bytes.Buffer{}

	sortCmd.Stdin, err = catCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = sortCmd.Start()
	if err != nil {
		//panic(err.Error() + ":" + sortCmd.Stderr.(*bytes.Buffer).String())
	}
	err = catCmd.Run()
	if err != nil {
		//panic(err.Error() + ":" + catCmd.Stderr.(*bytes.Buffer).String())
	}
	err = sortCmd.Wait()
	if err != nil {
		//panic(err.Error() + ":" + sortCmd.Stderr.(*bytes.Buffer).String())
	}

	fmt.Println(sortCmd.Stdout.(*bytes.Buffer).String())
}
