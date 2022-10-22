package data

import (
	"log"
	"os/exec"
)

const java string = "java"

func ChangeState(stub Stub) {
	if stub.state {
		if err := exec.Command("kill", stub.pid).Run(); err != nil {
			log.Fatalf("Fail to kill stub\n%s", err)
		}
	} else {
		cmd := exec.Command("nohup", java, "-Xmx1G", "-jar", stub.path, "&> /dev/null &")
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatalf("Fail to start stub\n%s", err)
		}
		//For some inexplicable reason if system waits for answer from Stdout(?), it waits forever
		// cmd.Wait()
	}
}
