package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/roshbhatia/sshtea/host"
)

func LoadHosts() []list.Item {
	configPath := os.ExpandEnv("$HOME/.ssh/config")
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading SSH config: %v\n", err)
		return nil
	}

	lines := strings.Split(string(content), "\n")
	var hosts []list.Item
	var currentHost host.Host
	lineNum := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Host ") {
			if currentHost.Name != "" {
				hosts = append(hosts, currentHost)
			}
			currentHost = host.Host{Name: strings.TrimPrefix(line, "Host "), ConfigLine: lineNum}
		} else if strings.HasPrefix(line, "HostName ") {
			currentHost.Hostname = strings.TrimPrefix(line, "HostName ")
		} else if strings.HasPrefix(line, "User ") {
			currentHost.User = strings.TrimPrefix(line, "User ")
		} else if strings.HasPrefix(line, "Port ") {
			currentHost.Port = strings.TrimPrefix(line, "Port ")
		}
		lineNum++
	}

	if currentHost.Name != "" {
		hosts = append(hosts, currentHost)
	}

	return hosts
}

func SaveHosts(hosts []list.Item) error {
	configPath := os.ExpandEnv("$HOME/.ssh/config")
	var content strings.Builder

	for _, item := range hosts {
		h := item.(host.Host)
		content.WriteString(fmt.Sprintf("Host %s\n", h.Name))
		if h.Hostname != "" {
			content.WriteString(fmt.Sprintf("    HostName %s\n", h.Hostname))
		}
		if h.User != "" {
			content.WriteString(fmt.Sprintf("    User %s\n", h.User))
		}
		if h.Port != "" {
			content.WriteString(fmt.Sprintf("    Port %s\n", h.Port))
		}
		content.WriteString("\n")
	}

	return ioutil.WriteFile(configPath, []byte(content.String()), 0600)
}
