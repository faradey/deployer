package commands

import "github.com/faradey/deployer/responser"

type Commander struct {
	AsyncCommands []CommandStruct
	Output        *responser.ResponseStruct
}

type CommandStruct struct {
	Command string
	Uid     uint64
	Gid     uint64
	Try     int
}

func (t *Commander) Runner(commandStr string, uid, gid uint64, runTry int, asyncGroup bool) {
	command := CommandStruct{Command: commandStr, Uid: uid, Gid: gid, Try: runTry}
	if asyncGroup {
		t.AsyncCommands = append(t.AsyncCommands, command)
	} else {
		run(command)
	}
}

func (t *Commander) RunAsync() {
	if len(t.AsyncCommands) > 0 {
		for _, command := range t.AsyncCommands {
			run(command)
		}
	}
}

func run(command CommandStruct) {

}
