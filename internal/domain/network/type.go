package network

import (
	"errors"
	"github.com/MaxFando/rate-limiter/internal/domain/network/ip"
	"github.com/MaxFando/rate-limiter/internal/domain/network/mask"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Ip       ip.Ip  `json:"ip"`
}

func NewRequest(
	Login string,
	Password string,
	Ip string,
) (Request, error) {
	if Login == "" || Password == "" {
		return Request{}, errors.New("login or password is empty")
	}

	_ip, err := ip.New(Ip)
	if err != nil {
		return Request{}, err
	}

	return Request{
		Login:    Login,
		Password: Password,
		Ip:       _ip,
	}, nil
}

type IpNetwork struct {
	Ip   ip.Ip     `json:"ip" db:"prefix"`
	Mask mask.Mask `json:"mask" db:"mask"`
}

func NewIpNetwork(ipValue string, maskValue string) (IpNetwork, error) {
	_ip, err := ip.New(ipValue)
	if err != nil {
		return IpNetwork{}, err
	}

	_mask, err := mask.New(maskValue)
	if err != nil {
		return IpNetwork{}, err
	}

	return IpNetwork{
		Ip:   _ip,
		Mask: _mask,
	}, nil
}
