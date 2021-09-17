package main

import (
	"errors"
	"fmt"
	"github.com/faradey/deployer/commands"
	"github.com/faradey/deployer/parser"
	"github.com/faradey/deployer/responser"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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

	http.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
		output := new(responser.ResponseStruct)
		commander := new(commands.Commander)
		commander.Output = output
		commander.Shell = shell
		var usr *user.User
		usr, _ = user.Current()
		var uid uint64
		uid, err = strconv.ParseUint(usr.Uid, 10, 32)
		if err != nil {
			output.SetMessage("An error occured while getting the user id by name " + usr.Username)
			output.SendError(w)
		}
		var gid uint64
		userGroups, err := usr.GroupIds()
		if err != nil {
			output.SetMessage("For the user named " + usr.Username + ", it was not possible to get the list of groups")
			output.SendError(w)
		}
		userGroup, err := user.LookupGroupId(userGroups[0])
		if err != nil {
			output.SetMessage("For the user named " + usr.Username + ", it was not possible to get the group")
			output.SendError(w)
		}
		gid, err = strconv.ParseUint(userGroup.Gid, 10, 32)
		if err != nil {
			output.SetMessage("For the user named " + usr.Username + ", it was not possible to get the group in uint64")
			output.SendError(w)
		}

		runTry := 1
		asyncGroup := false

		for _, row := range config {
			switch strings.ToLower(row.Name) {
			case "user":
				userName := strings.TrimSpace(row.Data)
				if userName != "" {
					usr, err = user.Lookup(userName)
					if err != nil {
						output.SetMessage("The user with the name " + userName + " is not in the system")
						output.SendError(w)
					}
					uid, err = strconv.ParseUint(usr.Uid, 10, 32)
					if err != nil {
						output.SetMessage("An error occured while getting the user id by name " + userName)
						output.SendError(w)
					}
				}
			case "user_group":
				userGroupName := strings.TrimSpace(row.Data)
				if userGroupName != "" {
					userGroup, err := user.LookupGroup(userGroupName)
					if err != nil {
						output.SetMessage("For the user named " + usr.Username + ", it was not possible to get the group")
						output.SendError(w)
					}
					gid, err = strconv.ParseUint(userGroup.Gid, 10, 32)
					if err != nil {
						output.SetMessage("For the user named " + usr.Username + ", it was not possible to get the group in uint64")
						output.SendError(w)
					}
				}
			case "try":
				runTry, err = strconv.Atoi(row.Data)
				if err != nil {
					output.SetMessage("Option TRY is specified incorrectly")
					output.SendError(w)
				}
			case "async_group_start":
				asyncGroup = true
			case "async_group_end":
				asyncGroup = false
				commander.RunAsync()
			case "run":
				if strings.TrimSpace(row.Data) != "" {
					commander.Runner(strings.TrimSpace(row.Data), uid, gid, runTry, asyncGroup)
				}
			}
		}

		output.Finish(w)
	})
	log.Println("Start Listener Host: " + host + " and Port: " + port)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func runCommand(w http.ResponseWriter, r *http.Request, conf Conf, command map[string][]string, dir, alloutput string) (string, error) {
	for i := 0; i < tryCount; i++ {
		cmd := exec.Command(command["name"][0], command["arg"]...)
		cmd.Dir = dir
		userOs := ""
		if val, ok := command["user"]; ok && len(val[0]) > 0 {
			userOs = command["user"][0]
		}

		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, alloutput+"\n"+fmt.Sprint(err))
			return "", err
		}
		userName := usr.Username
		uid, err := strconv.ParseUint(usr.Uid, 10, 32)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, alloutput+"\n"+fmt.Sprint(err))
			return "", err
		}
		gid := usr.Gid
		userGroup := ""
		if val, ok := command["user-group"]; ok && len(val[0]) > 0 {
			userGroup = command["user-group"][0]
		}
		var grp *user.Group
		if userGroup != "" {
			grp, err = user.LookupGroup(userGroup)
		} else if conf.UserGroup != "" {
			grp, err = user.LookupGroup(conf.UserGroup)
		} else {
			grps, err := currentUser.GroupIds()
			if err == nil {
				grp, err = user.LookupGroupId(grps[0])
			}
		}

		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, alloutput+"\n"+fmt.Sprint(err))
			return "", err
		}

		if grp != nil {
			gid = grp.Gid
			userGroup = grp.Name
		}
		Gid, _ := strconv.Atoi(gid)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(Gid)}
		output, err := cmd.CombinedOutput()
		userNGD := "\n" + "The command is executed by user " + userName + ":" + userGroup + "\n"
		if err != nil {
			if tryCount == i+1 {
				w.WriteHeader(400)
				fmt.Fprintf(w, alloutput+userNGD+fmt.Sprint(cmd)+"\n"+fmt.Sprint(err)+": "+string(output))
				return "", err
			} else {
				continue
			}
		}
		return userNGD + fmt.Sprint(cmd) + "\n" + string(output) + "\n", nil
	}
	w.WriteHeader(400)
	fmt.Fprintf(w, "something went wrong")
	return "", errors.New("something went wrong")
}
