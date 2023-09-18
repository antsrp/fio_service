package domain

type PersonCommon struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Patronym string `json:"patronym"`
}

type Person struct {
	Id int `json:"id,omitempty"`
	PersonCommon
	Age         uint   `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

type PersonMessage struct {
	PersonCommon
	EntityError string `json:"error,omitempty"`
}
