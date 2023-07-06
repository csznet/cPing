package conf

type ExRes struct {
	Status bool   `json:"status"`
	Result string `json:"result"`
}

type Conf struct {
	Client string `json:"client"`
	Server string `json:"server"`
	Token  string `json:"token"`
}

type ExReq struct {
	To    string `json:"to"`
	Token string `json:"token"`
	Stamp string `json:"stamp"`
}

type Client struct {
	List []Conf `json:"list"`
}
