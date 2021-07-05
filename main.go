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
	viper.SetDefault("commands", []map[string][]string{{"name": {"echo"}, "arg": {"Hello", " ", "world!"}, "async": {"false"}, "user": {"root"}}})
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

	type Conf struct {
		Host     string
		Port     int
		Path     string
		User     string
		Commands []map[string][]string
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
			cmd := exec.Command(command["name"][0], command["arg"]...)
			cmd.Dir = dir
			if conf.User != "root" && command["user"][0] != "root" {
				var usr *user.User
				if command["user"][0] != "" {
					usr, err = user.Lookup(command["user"][0])
				} else if conf.User != "" {
					usr, err = user.Lookup(conf.User)
				}

				if err != nil {
					w.WriteHeader(400)
					fmt.Fprintf(w, fmt.Sprint(err))
					return
				}

				uid, _ := strconv.ParseUint(usr.Uid, 10, 32)
				gid, _ := strconv.ParseUint(usr.Gid, 10, 32)

				cmd.SysProcAttr = &syscall.SysProcAttr{}
				cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
			}
			output, err := cmd.CombinedOutput()
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, fmt.Sprint(err)+": "+string(output))
				return
			}
			alloutput += string(output)
		}
		/*cmd := exec.Command("git", "-C", dir, "pull", "https://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@github.com/faradey/indemorio-www.git", "master")
		cmd = exec.Command("composer", "install")
		cmd = exec.Command("php", dir+"/bin/magento", "deploy:mode:set", "developer")
		cmd = exec.Command("php", dir+"/bin/magento", "s:up")
		cmd = exec.Command("php", dir+"/bin/magento", "deploy:mode:set", "production")*/

		fmt.Fprintf(w, alloutput)
	})

	log.Fatal(http.ListenAndServe(viper.GetString("host")+":"+viper.GetString("port"), nil))
}
