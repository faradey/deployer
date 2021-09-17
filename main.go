package main

import (
	"errors"
	"fmt"
	"github.com/faradey/deployer/commands"
	"github.com/faradey/deployer/parser"
	"github.com/faradey/deployer/theend"
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

	config = parser.GetConfig(dir)
	for _, row := range config {
		switch strings.ToLower(row.Name) {
		case "host":
			host = strings.TrimSpace(row.Data)
		case "port":
			port = strings.TrimSpace(row.Data)
		case "path":
			urlPath = strings.TrimSpace(row.Data)
		}
	}

	var AllOutput theend.ResponseTheEnd

	http.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
		AllOutput = theend.ResponseTheEnd{}
		commander := new(commands.Commander)
		var usr *user.User
		usr, _ = user.Current()
		var uid uint64
		uid, err = strconv.ParseUint(usr.Uid, 10, 32)
		if err != nil {
			AllOutput.SetMessage("An error occured while getting the user id by name " + usr.Username)
			AllOutput.TheEnd(w)
		}
		var gid uint64
		userGroups, err := usr.GroupIds()
		if err != nil {
			AllOutput.SetMessage("For the user named " + usr.Username + ", it was not possible to get the list of groups")
			AllOutput.TheEnd(w)
		}
		userGroup, err := user.LookupGroupId(userGroups[0])
		if err != nil {
			AllOutput.SetMessage("For the user named " + usr.Username + ", it was not possible to get the group")
			AllOutput.TheEnd(w)
		}
		gid, err = strconv.ParseUint(userGroup.Gid, 10, 32)
		if err != nil {
			AllOutput.SetMessage("For the user named " + usr.Username + ", it was not possible to get the group in uint64")
			AllOutput.TheEnd(w)
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
						AllOutput.SetMessage("The user with the name " + userName + " is not in the system")
						AllOutput.TheEnd(w)
					}
					uid, err = strconv.ParseUint(usr.Uid, 10, 32)
					if err != nil {
						AllOutput.SetMessage("An error occured while getting the user id by name " + userName)
						AllOutput.TheEnd(w)
					}
				}
			case "user_group":
				userGroupName := strings.TrimSpace(row.Data)
				if userGroupName != "" {
					userGroup, err := user.LookupGroup(userGroupName)
					if err != nil {
						AllOutput.SetMessage("For the user named " + usr.Username + ", it was not possible to get the group")
						AllOutput.TheEnd(w)
					}
					gid, err = strconv.ParseUint(userGroup.Gid, 10, 32)
					if err != nil {
						AllOutput.SetMessage("For the user named " + usr.Username + ", it was not possible to get the group in uint64")
						AllOutput.TheEnd(w)
					}
				}
			case "try":
				runTry, err = strconv.Atoi(row.Data)
				if err != nil {
					AllOutput.SetMessage("Option TRY is specified incorrectly")
					AllOutput.TheEnd(w)
				}
			case "async_group_start":
				asyncGroup = true
			case "async_group_end":
				asyncGroup = false
				commander.RunAsync()
			case "run":
				commander.Runner()
			}
		}
		/*

			ch := make(chan CommandResponse)
			for _, command := range conf.Commands {
				if val, ok := command["async"]; ok && val[0] == "true" {
					go runCommand(w, r, conf, command, dir, alloutput)
				} else {
					outputTemp, err := runCommand(w, r, conf, command, dir, alloutput)
					if err != nil {
						return
					}
					alloutput += outputTemp
				}
			}
			i := 0
			for loop := true; loop; {
				select {
				case msg := <-ch:
					fmt.Println(msg)
					i++
					loop = false
				}
			}

			fmt.Fprintf(w, alloutput)*/
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
