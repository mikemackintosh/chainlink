package main

var cmds = map[string][]string{
	"get_dns":      []string{"networksetup", "-getdnsservers", "%inf"},
	"set_dns":      []string{"networksetup", "-setdnsservers", "%inf"},
	"clear_dns":    []string{"networksetup", "-setdnsservers", "%inf", "\"Empty\""},
	"list_devices": []string{"networksetup", "-listallnetworkservices"},
	"restart_dns":  []string{"killall", "-HUP", "mDNSResponder"},
}
