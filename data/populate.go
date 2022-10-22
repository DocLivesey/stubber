package data

import (
	// "fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	// main "github.com/DocLivesey/stubber"
)

type Stub struct {
	jar   string
	path  string
	state bool
	pid   string
}

func (i Stub) Jar() string  { return i.jar }
func (i Stub) Path() string { return i.path }
func (i Stub) State() string {
	if i.state {
		return "On"
	}
	return "Off"
}
func (i Stub) Pid() string         { return i.pid }
func (i Stub) FilterValue() string { return i.jar }

func firstWord(value string) string {
	// Loop over all indexes in the string.
	count := 0
	start := false
	for i := range value {
		// If we encounter a space, reduce the count.
		if (value[i] == ' ') && !start {
			count += 1
		} else if value[i] != ' ' {
			start = true
		} else if value[i] == ' ' && start {
			// count -= 1
			// if count == 0 {
			return value[count:i]
			// }
		}
	}
	// Return the entire string.
	return value
}

func Populate() []Stub { //[]main.Stub {
	out, err := exec.Command("ps", "-e", "-o", "pid,command").CombinedOutput()
	if err != nil {
		log.Fatalf("Error on ps call: %s", err)
	}

	lines := strings.Split(string(out), "\n")

	var stubs []Stub //[]main.Stub

	dir := "/home/kuro/dev/tmp" //, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get workdir:  %s", err)
	}

	var dirRun func(string)

	dirRun = func(path string) {
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatalf("FileStat error\n\t%s", err)
		}
		if fileInfo.IsDir() {
			dirs, err := os.ReadDir(path)
			if err != nil {
				log.Fatalf("Error on reading directory: %s", err)
			}

			for _, d := range dirs {
				p := path + "/" + d.Name()
				dirRun(p)
				// println(d.Name())
			}
		} else {
			if strings.Contains(fileInfo.Name(), ".jar") {
				// fmt.Println("File name is", fileInfo.Name())
				// s := main.Stub{}
				// s.SetJar(fileInfo.Name())
				// s.SetPath(path)
				// s.SetState(false)
				// s.SetPid("-")
				s := Stub{fileInfo.Name(), path, false, "-"}
				for _, l := range lines {
					if strings.Contains(l, fileInfo.Name()) {
						// s.SetPid(firstWord(l))
						// s.SetState(true)
						s.pid = firstWord(l)
						s.state = true
						break
					}
				}
				stubs = append(stubs, s)
			}
		}
	}
	dirRun(dir)
	return stubs
}
