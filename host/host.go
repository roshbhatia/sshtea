package host

import "fmt"

type Host struct {
	Name       string
	Hostname   string
	User       string
	Port       string
	ConfigLine int
}

func (h Host) Title() string       { return h.Name }
func (h Host) Description() string { return fmt.Sprintf("%s@%s:%s", h.User, h.Hostname, h.Port) }
func (h Host) FilterValue() string { return h.Name }
