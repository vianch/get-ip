//go:build darwin

package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var (
	serviceMapOnce sync.Once
	serviceMap     map[string]string
)

var errNoCurrentWiFi = errors.New("no current WiFi network")

// Collect enumerates IPv4-bearing, up, non-loopback interfaces and enriches
// each with macOS-specific metadata (service name, gateway, DHCP/Manual).
func Collect() ([]Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	services := serviceNameMap()

	var out []Interface
	for _, ifi := range ifaces {
		if ifi.Flags&net.FlagUp == 0 {
			continue
		}
		if ifi.Flags&net.FlagLoopback != 0 {
			continue
		}
		ipv4 := firstIPv4(ifi)
		if ipv4 == "" {
			continue
		}

		row := Interface{
			Device:      ifi.Name,
			ServiceName: services[ifi.Name],
			IPv4:        ipv4,
			MAC:         ifi.HardwareAddr.String(),
		}

		if gw, err := gatewayFor(ifi.Name); err == nil {
			row.Gateway = gw
		} else {
			row.Gateway = "—"
		}

		row.Mode = modeFor(ifi.Name)

		out = append(out, row)
	}
	return out, nil
}

func firstIPv4(ifi net.Interface) string {
	addrs, err := ifi.Addrs()
	if err != nil {
		return ""
	}
	for _, a := range addrs {
		var ip net.IP
		switch v := a.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil {
			continue
		}
		v4 := ip.To4()
		if v4 == nil {
			continue
		}
		return v4.String()
	}
	return ""
}

func serviceNameMap() map[string]string {
	serviceMapOnce.Do(func() {
		out, err := exec.Command("networksetup", "-listnetworkserviceorder").Output()
		if err != nil {
			serviceMap = map[string]string{}
			return
		}
		serviceMap = parseServiceOrder(string(out))
	})
	return serviceMap
}

var hardwarePortRe = regexp.MustCompile(`\(Hardware Port:\s*([^,]+),\s*Device:\s*([^)]+)\)`)

func parseServiceOrder(s string) map[string]string {
	m := map[string]string{}
	for _, match := range hardwarePortRe.FindAllStringSubmatch(s, -1) {
		port := strings.TrimSpace(match[1])
		dev := strings.TrimSpace(match[2])
		if dev == "" {
			continue
		}
		m[dev] = port
	}
	return m
}

// modeFor: ipconfig getpacket emits an empty body for statically-configured
// interfaces — that absence is the signal, not an error condition.
func modeFor(dev string) Mode {
	cmd := exec.Command("ipconfig", "getpacket", dev)
	out, err := cmd.Output()
	if err != nil {
		return ModeManual
	}
	if len(bytes.TrimSpace(out)) == 0 {
		return ModeManual
	}
	if bytes.Contains(out, []byte("yiaddr")) {
		return ModeDHCP
	}
	return ModeManual
}

var gatewayRe = regexp.MustCompile(`(?m)^\s*gateway:\s*(\S+)`)

func gatewayFor(dev string) (string, error) {
	if gw, err := runRouteGet("-ifscope", dev, "default"); err == nil && gw != "" {
		return gw, nil
	}
	if gw, err := runRouteGet("default"); err == nil && gw != "" {
		return gw, nil
	}
	return "", errors.New("no gateway found")
}

func runRouteGet(args ...string) (string, error) {
	full := append([]string{"-n", "get"}, args...)
	out, err := exec.Command("route", full...).Output()
	if err != nil {
		return "", err
	}
	m := gatewayRe.FindStringSubmatch(string(out))
	if len(m) < 2 {
		return "", errors.New("no gateway line")
	}
	return m[1], nil
}

type spairportRoot struct {
	SPAirPortDataType []spairportEntry `json:"SPAirPortDataType"`
}

type spairportEntry struct {
	Interfaces []spairportIface `json:"spairport_airport_interfaces"`
}

type spairportIface struct {
	CurrentNetwork *spairportNetwork `json:"spairport_current_network_information"`
}

type spairportNetwork struct {
	Name        string `json:"_name"`
	SignalNoise string `json:"spairport_signal_noise"`
	Channel     string `json:"spairport_network_channel"`
}

// WiFi queries system_profiler for the currently-associated network only.
func WiFi() (WiFiInfo, error) {
	out, err := exec.Command("system_profiler", "SPAirPortDataType", "-json").Output()
	if err != nil {
		return WiFiInfo{}, err
	}
	return parseWiFiJSON(out)
}

func parseWiFiJSON(b []byte) (WiFiInfo, error) {
	var root spairportRoot
	if err := json.Unmarshal(b, &root); err != nil {
		return WiFiInfo{}, err
	}
	for _, e := range root.SPAirPortDataType {
		for _, ifi := range e.Interfaces {
			if ifi.CurrentNetwork == nil {
				continue
			}
			n := ifi.CurrentNetwork
			info := WiFiInfo{
				SSID:    n.Name,
				Channel: n.Channel,
				Signal:  n.SignalNoise,
			}
			// Location Services gating: SSID either missing or literal "<private>".
			if info.SSID == "" || info.SSID == "<private>" {
				info.Hidden = true
				info.SSID = "(SSID hidden — grant Location Services to Terminal)"
			}
			return info, nil
		}
	}
	return WiFiInfo{}, errNoCurrentWiFi
}
