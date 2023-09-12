package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

type Network interface {
	GetIPV6() (string, error)
	GetIPV4() (string, error)
}

type ipResponse struct {
	IP string `json:"ip"`
}

func NewNetwork(ipv6FetchAddress string, ipv4FetchAddress string) Network {
	return &NetworkHandler{
		ipv6FetchAddress: ipv6FetchAddress,
		ipv4FetchAddress: ipv4FetchAddress,
	}
}

type NetworkHandler struct {
	ipv6FetchAddress string
	ipv4FetchAddress string
}

func (n *NetworkHandler) GetIPV4() (string, error) {
	return getIP(n.ipv4FetchAddress)
}

func (n *NetworkHandler) GetIPV6() (string, error) {
	return getIP(n.ipv6FetchAddress)
}

func getIP(apiAddr string) (string, error) {
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
