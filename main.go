package main

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

type CommandResponse struct {
	Status  bool
	Message string
}

var config []parser.ConfigRow

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	/*f, err := os.OpenFile(dir+"/deployer_errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	defer f.Close()

	log.SetOutput(f)*/

	host := ""
	port := ""
	urlPath := ""
	shell := "bash"

	config = parser.GetConfig(dir)
	for _, row := range config {
		switch strings.ToLower(row.Name) {
		case "host":
			host = strings.TrimSpace(row.Data)
		case "port":
			port = strings.TrimSpace(row.Data)
		case "path":
			urlPath = strings.TrimSpace(row.Data)
		case "shell":
			if strings.TrimSpace(row.Data) != "" {
				shell = strings.TrimSpace(row.Data)
			}
		}
	}

	var usr *user.User
	usr, _ = user.Current()
	var uid uint64
	uid, err = strconv.ParseUint(usr.Uid, 10, 32)
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
	var gid uint64
	gid, err = strconv.ParseUint(userGroup.Gid, 10, 32)
	if err != nil {
		panic("For the user named " + usr.Username + ", it was not possible to get the group in uint64")
	}

	http.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
		output := new(responser.ResponseStruct)
		commander := new(commands.Commander)
		commander.Output = output
		asyncGroup := false

		for _, row := range config {
			commandAttributes := commands.GetAttributesStruct()
			commandAttributes.Shell = shell
			commandAttributes.Try = 1
			commandAttributes.Cd = ""
			commandAttributes.Uid = uid
			commandAttributes.Gid = gid

			rowData := strings.TrimSpace(row.Data)
			switch strings.ToLower(row.Name) {
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
						commandAttributes.Uid, err = strconv.ParseUint(usr.Uid, 10, 32)
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
								commandAttributes.Gid, err = strconv.ParseUint(userGroup.Gid, 10, 32)
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
				commandAttributes.Try, err = strconv.Atoi(rowData)
				if err != nil {
					output.AddMessage("Option TRY is specified incorrectly")
					output.SendError(w)
					return
				}
			case "cd":
				commandAttributes.Cd = rowData
			case "async_group_start":
				asyncGroup = true
			case "async_group_end":
				asyncGroup = false
				commander.RunAsyncCommands()
			case "run":
				if rowData != "" {
					commandAttributes.Command = rowData
					commander.Runner(commandAttributes, asyncGroup)
				}
			}
		}

		output.Finish(w)
	})
	log.Println("Start Listener Host: " + host + " and Port: " + port)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}
