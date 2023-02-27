package types

type Config struct {
	Endpoint      string `json:"endpoint"`
	Private       bool   `json:"private"`
	Authorization string `json:"authorization"`
}
