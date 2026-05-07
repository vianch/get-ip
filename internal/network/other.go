//go:build !darwin

package network

import "errors"

var errUnsupported = errors.New("unsupported platform")

func Collect() ([]Interface, error) {
	return nil, errUnsupported
}

func WiFi() (WiFiInfo, error) {
	return WiFiInfo{}, errUnsupported
}
