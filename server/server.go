package server

import (
	"github.com/faradey/deployer/commands"
	"github.com/faradey/deployer/parser"
	"github.com/faradey/deployer/responser"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var config []parser.ConfigRow

type MainConfigStruct struct {
	Host    string
	Port    string
	UrlPath string
	Dir     string
	Shell   string
}

func StartServer() {
	defConf := commands.DefaultCommandStruct{}
	mainConfig := GetMainConfig()
	defConf.Shell = mainConfig.Shell
	var usr *user.User
	usr, _ = user.Current()
	var err error
	defConf.Uid, err = strconv.ParseUint(usr.Uid, 10, 32)
	if err != nil {
		panic("An error occured while getting the user id by name " + usr.Username)
	}
	userGroups, err := usr.GroupIds()
	if err != nil {
		panic("For the user named " + usr.Username + ", it was not possible to get the list of groups")
	}
	userGroup, err := user.LookupGroupId(userGroups[0])
	if err != nil {
		panic("For the user named " + usr.Username + ", it was not possible to get the group")
	}

	defConf.Gid, err = strconv.ParseUint(userGroup.Gid, 10, 32)
	if err != nil {
		panic("For the user named " + usr.Username + ", it was not possible to get the group in uint64")
	}
	defConf.Try = 1
	defConf.Cd = ""

	http.HandleFunc(mainConfig.UrlPath, func(w http.ResponseWriter, r *http.Request) {
		output := new(responser.ResponseStruct)
		commander := new(commands.Commander)
		commander.Output = output
		commander.Channel = make(chan bool, len(config))
		asyncGroup := false
		countRunCommands := 0

		commandDefConf := defConf

		for _, row := range config {
			rowData := strings.TrimSpace(row.Data)
			optionName := strings.ToLower(row.Name)
			switch optionName {
			case "user":
				userData := strings.SplitN(rowData, " ", 2)
				if len(userData) > 0 {
					userName := userData[0]
					if userName != "" {
						usr, err = user.Lookup(userName)
						if err != nil {
							output.AddMessage("The user with the name " + userName + " is not in the system")
							output.SendError(w)
							return
						}
						commandDefConf.Uid, err = strconv.ParseUint(usr.Uid, 10, 32)
						if err != nil {
							output.AddMessage("An error occured while getting the user id by name " + userName)
							output.SendError(w)
							return
						}

						if len(userData) > 1 {
							userGroupName := strings.TrimSpace(userData[1])
							if userGroupName != "" {
								userGroup, err := user.LookupGroup(userGroupName)
								if err != nil {
									output.AddMessage("For the user named " + usr.Username + ", it was not possible to get the group")
									output.SendError(w)
									return
								}
								commandDefConf.Gid, err = strconv.ParseUint(userGroup.Gid, 10, 32)
								if err != nil {
									output.AddMessage("For the user named " + usr.Username + ", it was not possible to get the group in uint64")
									output.SendError(w)
									return
								}
							}
						}
					}
				}
			case "try":
				commandDefConf.Try, err = strconv.Atoi(rowData)
				if err != nil {
					output.AddMessage("Option TRY is specified incorrectly")
					output.SendError(w)
					return
				}
			case "cd":
				commandDefConf.Cd = rowData
			case "async_group_start":
				asyncGroup = true
				commander.CreateAsyncGroup()
			case "async_group_end":
				asyncGroup = false
				commander.RunAsyncCommands()
			case "run", "async_run":
				if rowData != "" {
					countRunCommands++
					commandAttributes := commands.GetAttributesStruct()
					commandAttributes.Command = rowData
					commandAttributes.DefaultCommandStruct = commandDefConf
					localAsyncGroup := asyncGroup
					if optionName == "async_run" {
						commandAttributes.Async = true
						localAsyncGroup = false
					}

					commander.Runner(commandAttributes, localAsyncGroup)
				}
			}
		}

		if countRunCommands > 0 {
			loop := true
			i := 0
			for loop {
				select {
				case <-commander.Channel:
					i++
					if i >= countRunCommands {
						loop = false
					}
				}
			}
		}

		output.Finish(w)
	})
	log.Println("Start Listener Host: " + mainConfig.Host + " and Port: " + mainConfig.Port)
	log.Fatal(http.ListenAndServe(mainConfig.Host+":"+mainConfig.Port, nil))
}

func getDir(arg0 string) string {
	dir, err := filepath.Abs(filepath.Dir(arg0))
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func GetMainConfig() MainConfigStruct {
	mainConfig := MainConfigStruct{}
	mainConfig.Dir = getDir(os.Args[0])
	mainConfig.Host = ""
	mainConfig.Port = ""
	mainConfig.UrlPath = ""
	mainConfig.Shell = "bash"

	config = parser.GetConfig(mainConfig.Dir)
	for _, row := range config {
		switch strings.ToLower(row.Name) {
		case "host":
			mainConfig.Host = strings.TrimSpace(row.Data)
		case "port":
			mainConfig.Port = strings.TrimSpace(row.Data)
		case "path":
			mainConfig.UrlPath = strings.TrimSpace(row.Data)
		case "shell":
			if strings.TrimSpace(row.Data) != "" {
				mainConfig.Shell = strings.TrimSpace(row.Data)
			}
		}
	}

	return mainConfig
}
