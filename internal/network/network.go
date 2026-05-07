package network

type Mode string

const (
	ModeDHCP    Mode = "DHCP"
	ModeManual  Mode = "Manual"
	ModeUnknown Mode = "—"
)

type Interface struct {
	Device      string
	ServiceName string
	IPv4        string
	MAC         string
	Gateway     string
	Mode        Mode
}

type WiFiInfo struct {
	SSID    string
	Channel string
	Signal  string
	Hidden  bool
}

func (i Interface) Label() string {
	if i.ServiceName == "" {
		return i.Device
	}
	return i.Device + " (" + i.ServiceName + ")"
}
