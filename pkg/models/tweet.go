package models

type Keyer interface {
	Key() string
}

type Tweet struct {
	UID    string `json:"UID"`
	Author string `json:"author"`
	Tweet  string `json:"tweet"`
}

func (t Tweet) Key() string { //refer to interface implementation
	return "tweet:" + t.UID
}
