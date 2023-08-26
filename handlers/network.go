package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

type Network interface {
	GetIPV6() (string, error)
}

type ipResponse struct {
	IP string `json:"ip"`
}

func NewNetwork(ipv6FetchAddress string) Network {
	return &NetworkHandler{
		ipv6FetchAddress: ipv6FetchAddress,
	}
}

type NetworkHandler struct {
	ipv6FetchAddress string
}

func (n *NetworkHandler) GetIPV6() (string, error) {
	httpResponse, err := http.Get(n.ipv6FetchAddress)
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
