package network

import (
	"errors"
	"github.com/MaxFando/rate-limiter/internal/domain/network/ip"
	"github.com/MaxFando/rate-limiter/internal/domain/network/mask"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IP       ip.IP  `json:"ip"`
}

func NewRequest(
	login string,
	password string,
	ipNew string,
) (Request, error) {
	if login == "" || password == "" {
		return Request{}, errors.New("login or password is empty")
	}

	_ip, err := ip.New(ipNew)
	if err != nil {
		return Request{}, err
	}

	return Request{
		Login:    login,
		Password: password,
		IP:       _ip,
	}, nil
}

type IPNetwork struct {
	IP   ip.IP     `json:"ip" db:"prefix"`
	Mask mask.Mask `json:"mask" db:"mask"`
}

func NewIPNetwork(ipValue string, maskValue string) (IPNetwork, error) {
	_ip, err := ip.New(ipValue)
	if err != nil {
		return IPNetwork{}, err
	}

	_mask, err := mask.New(maskValue)
	if err != nil {
		return IPNetwork{}, err
	}

	return IPNetwork{
		IP:   _ip,
		Mask: _mask,
	}, nil
}
