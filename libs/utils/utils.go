package utils

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func GetMachineIpAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("Could not find IP address of machine")
}

func GetEnvVariableString(name string, defaultValue string) string {
	listenPortEnvVariable, foundListenPortEnvVariable := os.LookupEnv(name)
	if foundListenPortEnvVariable {
		return listenPortEnvVariable
	}

	return defaultValue
}

func GetEnvVariableInt(name string, defaultValue int) int {
	envVariable, foundEnvVariable := os.LookupEnv(name)
	if foundEnvVariable {
		value, err := strconv.Atoi(envVariable)
		if err != nil {
			return defaultValue
		}
		return value
	}

	return defaultValue
}

func GetEnvVariableBool(name string, defaultValue bool) bool {
	listenPortEnvVariable, foundListenPortEnvVariable := os.LookupEnv(name)
	if foundListenPortEnvVariable {
		lowerCased := strings.ToLower(listenPortEnvVariable)
		if lowerCased == "true" {
			return true
		} else if lowerCased == "false" {
			return false
		} else {
			return defaultValue
		}
	}

	return defaultValue
}

type ConfigLocation struct {
	ConfigPath string `mapstructure:"CONFIG_PATH"`
	ConfigName string `mapstructure:"CONFIG_NAME"`
}

func ProvideConfigLocation() *ConfigLocation {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "."
	}
	name := os.Getenv("CONFIG_NAME")

	return &ConfigLocation{
		ConfigPath: path,
		ConfigName: name,
	}
}

func LoadConfigFromEnvs(config interface{}) error {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "."
	}
	name := os.Getenv("CONFIG_NAME")
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	fmt.Println(path)
	fmt.Println(name)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&config)

	return err
}
