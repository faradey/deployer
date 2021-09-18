package commands

import (
	"fmt"
	"github.com/faradey/deployer/responser"
	"os/exec"
	"path/filepath"
	"syscall"
)

type Commander struct {
	AsyncCommands []*CommandStruct
	Output        *responser.ResponseStruct
}

type CommandStruct struct {
	Command string
	Uid     uint64
	Gid     uint64
	Try     int
	Shell   string
	Cd      string
	Async   bool
	tried   int
	result  bool
}

func (t *Commander) Runner(command *CommandStruct, asyncGroup bool) {
	if asyncGroup {
		t.AsyncCommands = append(t.AsyncCommands, command)
	} else if command.Async {
		go t.run(command)
	} else {
		t.run(command)
	}
}

func (t *Commander) RunAsyncCommands() bool {
	result := true
	countCommands := len(t.AsyncCommands)
	if countCommands > 0 {
		ch := make(chan bool)
		for _, command := range t.AsyncCommands {
			go t.runAsync(command, ch)
		}

		loop := true
		i := 0
		for loop {
			select {
			case msg := <-ch:
				i++
				if !msg || i == countCommands {
					result = msg
					loop = false
				}
			}
		}
	}

	return result
}

func (t *Commander) run(command *CommandStruct) bool {
	if command.tried == command.Try {
		return command.result
	}
	command.tried++
	cmd := exec.Command("/bin/"+command.Shell, "-c", command.Command)
	dir, err := filepath.Abs(command.Cd)
	if err != nil {
		t.Output.AddMessage("Option CD is specified incorrectly")
		return false
	}
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(command.Uid), Gid: uint32(command.Gid)}
	output, err := cmd.CombinedOutput()
	t.Output.AddMessage(cmd.Dir + "$ " + command.Command + "\n" + string(output))
	if err != nil {
		t.Output.AddMessage(fmt.Sprint(err))
		command.result = false
		return t.run(command)
	}

	return true
}

func (t *Commander) runAsync(command *CommandStruct, ch chan bool) {
	ch <- t.run(command)
}

func GetAttributesStruct() *CommandStruct {
	return new(CommandStruct)
}
