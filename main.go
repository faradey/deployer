package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

type Conf struct {
	Host      string
	Port      int
	Path      string
	User      string
	UserGroup string
	Commands  []map[string][]string
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	/* Default options */
	viper.SetDefault("host", "")
	viper.SetDefault("port", "8083")
	viper.SetDefault("path", "/deploy/123456789/abcdefg")
	viper.SetDefault("user", "root")
	viper.SetDefault("user-group", "root")
	viper.SetDefault("commands", []map[string][]string{{"name": {"echo"}, "arg": {"Hello", " ", "world!"}, "async": {"false"}, "user": {"root"}, "user-group": {"root"}}})
	/* END Default options */

	viper.SetConfigName("deployer-config")
	viper.AddConfigPath(dir + "/.")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = viper.SafeWriteConfig()
			if err != nil {
				log.Fatal(err)
			} else {
				log.Print("Create a new config file")
			}
		} else {
			log.Fatal(err)
		}
	}

	http.HandleFunc(viper.GetString("path"), func(w http.ResponseWriter, r *http.Request) {
		alloutput := ""
		var conf Conf
		err := viper.Unmarshal(&conf)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, fmt.Sprint(err))
			return
		}
		for _, command := range conf.Commands {
			if val, ok := command["async"]; ok && val[0] == "true" {
				go runCommand(w, r, conf, command, dir, alloutput)
			} else {
				outputTemp, err := runCommand(w, r, conf, command, dir, alloutput)
				if err != nil {
					w.WriteHeader(400)
					fmt.Fprintf(w, fmt.Sprint(err))
					return
				}
				alloutput += outputTemp
			}
		}

		fmt.Fprintf(w, alloutput)
	})

	log.Fatal(http.ListenAndServe(viper.GetString("host")+":"+viper.GetString("port"), nil))
}

func runCommand(w http.ResponseWriter, r *http.Request, conf Conf, command map[string][]string, dir, alloutput string) (string, error) {
	cmd := exec.Command(command["name"][0], command["arg"]...)
	cmd.Dir = dir
	userOs := ""
	if val, ok := command["user"]; ok && len(val[0]) > 0 {
		userOs = command["user"][0]
	}
	var err error
	if conf.User != "root" && userOs != "root" {
		var usr *user.User
		if userOs != "" {
			usr, err = user.Lookup(userOs)
		} else if conf.User != "" {
			usr, err = user.Lookup(conf.User)
		}

		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, alloutput+"\n"+fmt.Sprint(err))
			return "", err
		}
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
		if conf.UserGroup != "root" && userGroup != "root" {
			var grp *user.Group
			if userGroup != "" {
				grp, err = user.LookupGroup(userGroup)

			} else if conf.UserGroup != "" {
				grp, err = user.LookupGroup(conf.UserGroup)
			}

			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, alloutput+"\n"+fmt.Sprint(err))
				return "", err
			}
			gid = grp.Gid
		}
		Gid, _ := strconv.Atoi(gid)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(Gid)}
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, alloutput+"\n"+fmt.Sprint(err)+": "+string(output))
		return "", err
	}
	return fmt.Sprint(cmd) + "\n" + string(output) + "\n", nil
}
