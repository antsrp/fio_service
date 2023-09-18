package benefication

type Genderize struct {
	Common
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}
