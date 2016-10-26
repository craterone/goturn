package libs

import (
	"os"
	"encoding/json"
	"log"
)

const (
	APP_NAME    = "实时猫 TURN/STUN 服务器"
	APP_VERSION = "0.0.1"
)

var (
	config_path_array = []string{"config.json", "/etc/goturn/config.json"}
)

var (
	Config            *Configuration
)

type Configuration struct {
	Realm string

	LogLevel string `json:"logLevel"`

	LogToFile bool `json:"logToFile"`

	LogFilePath string `json:"logFilePath"`

	ErrLogFilePath string `json:"errLogFilePath"`

}


func LoadConfigurationModule() {
	var err error
	switch {
	case len(*config) > 0:
		log.Printf("Read config from  -> %s", *config)
		Config, err = readConfigByPath(*config)
	case len(*config) == 0:
		log.Println("Read config from default path")
		Config, err = readConfigDefault()
	}

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
		return
	}

	PrintModuleLoaded("Configuration")
}

func readConfigDefault() (*Configuration, error) {

	for _, config_path := range config_path_array {
		config, err := readConfigByPath(config_path)

		if err != nil {
			continue
		}

		return config, nil
	}

	return nil, ERROR_NO_SUITABLE_CONFIG
}


func readConfigByPath(config_path string) (*Configuration, error) {

	file, err := os.Open(config_path)

	defer file.Close()

	if err != nil {
		log.Printf("Can not read config from %s", config_path)
		return nil, err
	}

	decoder := json.NewDecoder(file)

	var config Configuration
	json_err := decoder.Decode(&config)

	if json_err != nil {
		log.Printf("Json parses failed from %s , error -> %v", config_path, json_err)
		return nil, json_err
	}

	log.Printf("Read config from %s successfully\n", config_path)

	return &config, nil
}