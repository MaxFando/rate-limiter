package network

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
}

type IpNetwork struct {
	Ip   string `json:"ip" db:"prefix"`
	Mask string `json:"mask" db:"mask"`
}
