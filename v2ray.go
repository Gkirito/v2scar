package v2scar

import (
	"context"
	"io"
	"log"
	"os/exec"
	"strings"
)

var V2rayRunFlag = make(chan bool)
var CloseV2rayNotifier = make(chan bool)

func RunV2ray(config string) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "sh", "-c", "v2ray -config="+config)
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %s......", err.Error())
	}
	go cancelV2ray(cancel)
	go asyncLog(stdout)
	_ = cmd.Wait()
}

func cancelV2ray(cancel context.CancelFunc)  {
	<-CloseV2rayNotifier
	log.Println("close v2ray")
	cancel()
}

func asyncLog(reader io.ReadCloser) error {
	cache := ""
	buf := make([]byte, 1024)
	for {
		num, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if num > 0 {
			b := buf[:num]
			s := strings.Split(string(b), "\n")
			line := strings.Join(s[:len(s)-1], "\n")
			if strings.Contains(line, "started") {
				V2rayRunFlag <- true
			}
			log.Printf("%s%s\n", cache, line)
			cache = s[len(s)-1]
		}
	}
}
