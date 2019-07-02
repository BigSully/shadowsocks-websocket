package shadowsocks

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

func ParseArgs() (config map[string]interface{}) {
	// default
	config = map[string]interface{}{
		"server":        "0.0.0.0",
		"server_port":   80,
		"local_port":    1080,
		"local_address": "127.0.0.1",
		"password":      "barfoo!",
		"method":        "aes-256-cfb",
		"timeout":       600,
	}

	// choose to parse the option -c manually because of the quirks of the flag library
	for index, arg := range os.Args {
		if arg != "-c" || index >= len(os.Args)-1 {
			continue
		}
		configFile := os.Args[index+1]
		// load the config.json file if it exists
		if _, err := os.Stat(configFile); err == nil {
			data, err := ioutil.ReadFile(configFile)
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal([]byte(string(data)), &config); err != nil {
				panic(err)
			}
			return
		}

	}

	options := map[string]string{
		"s": "server",
		"p": "server_port",
		"l": "local_port",
		"b": "local_address",
		"k": "password",
		"m": "method",
		"t": "timeout"}

	// update them with environment variables
	for _, name := range options {
		if value, isExist := os.LookupEnv(strings.ToUpper(name)); isExist {
			config[name] = value
		}
	}

	// update them with arguments    parse arguments manually
	for index, arg := range os.Args {
		// must start with a dash char -  and have next arg
		if !strings.HasPrefix(arg, "-") {
			continue
		}
		if index >= len(os.Args)-1 || strings.HasPrefix(os.Args[index+1], "-") {
			config[arg[1:]] = nil
			continue
		}
		// update map if the opt exists in the options
		next := os.Args[index+1]
		if name, ok := options[arg[1:]]; ok {
			config[name] = next
		} else {
			config[arg[1:]] = next
		}
	}

	return
}
