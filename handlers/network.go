package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
)

type Network interface {
	GetPublicIPV6() (string, error)
	GetPublicIPV4() (string, error)
	GetPrivateIPV4() (string, error)
}

type ipResponse struct {
	IP string `json:"ip"`
}

func NewNetwork(ipv6FetchAddress string, ipv4FetchAddress string, ipv4InterfaceName string) Network {
	return &NetworkHandler{
		ipv6FetchAddress:  ipv6FetchAddress,
		ipv4FetchAddress:  ipv4FetchAddress,
		ipv4InterfaceName: ipv4InterfaceName,
	}
}

type NetworkHandler struct {
	ipv6FetchAddress  string
	ipv4FetchAddress  string
	ipv4InterfaceName string
}

func (n *NetworkHandler) GetPublicIPV4() (string, error) {
	return getPublicIP(n.ipv4FetchAddress)
}

func (n *NetworkHandler) GetPublicIPV6() (string, error) {
	return getPublicIP(n.ipv6FetchAddress)
}

func (n *NetworkHandler) GetPrivateIPV4() (string, error) {
	var addrs []net.Addr
	var err error

	if n.ipv4InterfaceName != "" {
		iface, err := net.InterfaceByName(n.ipv4InterfaceName)
		if err != nil {
			return "", err
		}

		addrs, err = iface.Addrs()
		if err != nil {
			return "", err
		}
	} else {
		addrs, err = net.InterfaceAddrs()
		if err != nil {
			return "", err
		}
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no private ipv4 address found")
}

func getPublicIP(apiAddr string) (string, error) {
	httpResponse, err := http.Get(apiAddr)
	if err != nil {
		return "", err
	}

	responseBodyBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return "", err
	}

	res := &ipResponse{}

	err = json.Unmarshal(responseBodyBytes, res)
	if err != nil {
		return "", err
	}

	return res.IP, nil
}
