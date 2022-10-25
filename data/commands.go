package data

import (
	"log"
	"os/exec"

	"github.com/sadlil/gologger"
)

const java string = "java"

func ChangeState(stub Stub) {
	logger := gologger.GetLogger(gologger.FILE, "log_stubber.log")
	if stub.state {
		if err := exec.Command("kill", stub.pid).Run(); err != nil {
			log.Fatalf("Fail to kill stub\n%s", err)
		}
		logger.Log("Command: kill " + stub.pid)
	} else {
		cmd := exec.Command("nohup", java, "-Xmx1G", "-jar", stub.path, "&> /dev/null &")
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatalf("Fail to start stub\n%s", err)
		}

		logger.Log("Command: " + cmd.String())

		//For some inexplicable reason if system waits for answer from Stdout(?), it waits forever
		// cmd.Wait()
	}
}
