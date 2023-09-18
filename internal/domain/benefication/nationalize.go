package benefication

type Country struct {
	ID          string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type Nationalize struct {
	Common
	Country []Country `json:"country"`
}
