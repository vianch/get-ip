//go:build darwin

package network

import "testing"

const networksetupFixture = `An asterisk (*) denotes that a network service is disabled.
(1) Wi-Fi
(Hardware Port: Wi-Fi, Device: en0)

(2) USB 10/100/1000 LAN
(Hardware Port: USB 10/100/1000 LAN, Device: en7)

(3) Thunderbolt Bridge
(Hardware Port: Thunderbolt Bridge, Device: bridge0)
`

func TestParseServiceOrder(t *testing.T) {
	m := parseServiceOrder(networksetupFixture)
	cases := map[string]string{
		"en0":     "Wi-Fi",
		"en7":     "USB 10/100/1000 LAN",
		"bridge0": "Thunderbolt Bridge",
	}
	for dev, want := range cases {
		got, ok := m[dev]
		if !ok {
			t.Errorf("missing entry for %q", dev)
			continue
		}
		if got != want {
			t.Errorf("dev %q: got %q, want %q", dev, got, want)
		}
	}
	if len(m) != 3 {
		t.Errorf("got %d entries, want 3", len(m))
	}
}

const wifiJSONFixture = `{
  "SPAirPortDataType": [
    {
      "spairport_airport_interfaces": [
        {
          "_name": "en0",
          "spairport_current_network_information": {
            "_name": "MyHomeNetwork",
            "spairport_signal_noise": "-52 dBm / -90 dBm",
            "spairport_network_channel": "36 (5GHz, 80MHz)"
          }
        }
      ]
    }
  ]
}`

func TestParseWiFiJSON(t *testing.T) {
	info, err := parseWiFiJSON([]byte(wifiJSONFixture))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.SSID != "MyHomeNetwork" {
		t.Errorf("SSID: got %q, want %q", info.SSID, "MyHomeNetwork")
	}
	if info.Channel != "36 (5GHz, 80MHz)" {
		t.Errorf("Channel: got %q", info.Channel)
	}
	if info.Signal != "-52 dBm / -90 dBm" {
		t.Errorf("Signal: got %q", info.Signal)
	}
	if info.Hidden {
		t.Errorf("expected Hidden=false")
	}
}

const wifiJSONHiddenFixture = `{
  "SPAirPortDataType": [
    {
      "spairport_airport_interfaces": [
        {
          "spairport_current_network_information": {
            "_name": "<private>",
            "spairport_signal_noise": "-60 dBm / -90 dBm",
            "spairport_network_channel": "11"
          }
        }
      ]
    }
  ]
}`

func TestParseWiFiJSONHidden(t *testing.T) {
	info, err := parseWiFiJSON([]byte(wifiJSONHiddenFixture))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !info.Hidden {
		t.Errorf("expected Hidden=true")
	}
	if info.SSID == "<private>" || info.SSID == "" {
		t.Errorf("expected SSID to be replaced, got %q", info.SSID)
	}
}
