package commands

import (
	"fmt"
	"github.com/faradey/deployer/responser"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

type Commander struct {
	AsyncCommands *AsyncCommandsStruct
	Output        *responser.ResponseStruct
	Channel       chan bool
}

type DefaultCommandStruct struct {
	Uid   uint64
	Gid   uint64
	Shell string
	Try   int
	Cd    string
}

type CommandStruct struct {
	Command string
	Async   bool
	DefaultCommandStruct
	tried  int
	result bool
}

type AsyncCommandsStruct struct {
	Commands []*CommandStruct
}

func (t *Commander) Runner(command *CommandStruct, asyncGroup bool) {
	if asyncGroup {
		t.AsyncCommands.Commands = append(t.AsyncCommands.Commands, command)
	} else if command.Async {
		go t.run(command)
	} else {
		t.run(command)
	}
}

func (t *Commander) RunAsyncCommands() bool {
	result := true
	countCommands := len(t.AsyncCommands.Commands)
	if countCommands > 0 {
		ch := make(chan bool, 100)
		for _, command := range t.AsyncCommands.Commands {
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

func (t *Commander) CreateAsyncGroup() {
	t.AsyncCommands = new(AsyncCommandsStruct)
}

func (t *Commander) run(command *CommandStruct) bool {
	if command.tried == command.Try || command.result {
		t.Channel <- true
		return command.result
	}
	command.result = true
	command.tried++
	cmd := exec.Command("/bin/"+command.Shell, "-c", command.Command)
	dir, err := filepath.Abs(command.Cd)
	if err != nil {
		t.Output.AddMessage("$ " + command.Command + "\nOption CD is specified incorrectly")
		command.result = false
	}
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(command.Uid), Gid: uint32(command.Gid)}
	output, err := cmd.CombinedOutput()
	usr, err := user.LookupId(strconv.FormatUint(command.Uid, 10))
	if err != nil {
		t.Output.AddMessage(cmd.Dir + "$ " + command.Command + "\n" + string(output))
		t.Output.AddMessage(fmt.Sprint(err))
		command.result = false
	}
	hostname, err := os.Hostname()
	if err != nil {
		t.Output.AddMessage(usr.Username + "@:" + cmd.Dir + "$ " + command.Command + "\n" + string(output))
		t.Output.AddMessage(fmt.Sprint(err))
		command.result = false
	}
	t.Output.AddMessage(usr.Username + "@" + hostname + ":" + cmd.Dir + "$ " + command.Command + "\n" + string(output))
	if err != nil {
		t.Output.AddMessage(fmt.Sprint(err))
		command.result = false
	}

	return t.run(command)
}

func (t *Commander) runAsync(command *CommandStruct, ch chan bool) {
	ch <- t.run(command)
}

func GetAttributesStruct() *CommandStruct {
	return new(CommandStruct)
}
