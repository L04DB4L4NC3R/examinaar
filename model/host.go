package model

type HostType struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Port1    string `json:"port1"`
	Port2    string `json:"port2"`
	Image1   string `json:"image1"`
	Image2   string `json:"image2"`
}
