package benefication

import "github.com/antsrp/fio_service/internal/mapper"

type Common struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
}

type DataCacheValue struct {
	Age     uint
	Gender  string
	Country []Country
}

func (d DataCacheValue) MarshalBinary() ([]byte, error) {
	return mapper.ToJSON(d, mapper.NewIndent("", ""))
}

func (d *DataCacheValue) UnmarshalBinary(data []byte) error {
	r, err := mapper.FromJSON[DataCacheValue](data)
	if err != nil {
		return err
	}
	d.Age = r.Age
	d.Country = r.Country
	d.Gender = r.Gender

	return nil
}
