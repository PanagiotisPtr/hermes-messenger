package utils

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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
