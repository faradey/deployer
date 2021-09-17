package commands

import (
	"github.com/faradey/deployer/responser"
	"os/exec"
)

type Commander struct {
	AsyncCommands []CommandStruct
	Output        *responser.ResponseStruct
	Shell         string
}

type CommandStruct struct {
	Command string
	Uid     uint64
	Gid     uint64
	Try     int
	Shell   string
}

func (t *Commander) Runner(commandStr string, uid, gid uint64, runTry int, asyncGroup bool) {
	command := CommandStruct{Command: commandStr, Uid: uid, Gid: gid, Try: runTry, Shell: t.Shell}
	if asyncGroup {
		t.AsyncCommands = append(t.AsyncCommands, command)
	} else {
		run(command)
	}
}

func (t *Commander) RunAsync() {
	if len(t.AsyncCommands) > 0 {
		ch := make(chan string)
		for _, command := range t.AsyncCommands {
			go run(command)
		}
	}
}

func run(command CommandStruct) {
	cmd := exec.Command("/bin/"+command.Shell, "-c", command.Command)
	err := cmd.Run()
}
