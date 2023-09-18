package service

import (
	"fmt"
	"sync"

	"github.com/antsrp/fio_service/internal/domain/benefication"
	"github.com/antsrp/fio_service/internal/infrastructure/http/web"
	iweb "github.com/antsrp/fio_service/internal/interfaces/web"
)

const (
	requestAgify       = "https://api.agify.io"
	requestGenderize   = "https://api.genderize.io"
	requestNationalize = "https://api.nationalize.io"
)

func beneficate(conn iweb.Connector, name string) (benefication.DataCacheValue, error) {
	options := map[string]interface{}{
		"name": name,
	}
	a, g, n := benefication.Agify{}, benefication.Genderize{}, benefication.Nationalize{}

	var groupErrors []error
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		if _, err := conn.Do(web.MakeGETRequest(requestAgify, options), &a); err != nil {
			groupErrors = append(groupErrors, fmt.Errorf("can't agify name: %w", err))
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := conn.Do(web.MakeGETRequest(requestGenderize, options), &g); err != nil {
			groupErrors = append(groupErrors, fmt.Errorf("can't genderize name: %w", err))
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := conn.Do(web.MakeGETRequest(requestNationalize, options), &n); err != nil {
			groupErrors = append(groupErrors, fmt.Errorf("can't nationalize name: %w", err))
		}
	}()

	wg.Wait()

	if n := len(groupErrors); n > 0 {
		final := groupErrors[0]
		for i := 1; i < n; i++ {
			final = fmt.Errorf("%v,%v", final, groupErrors[i])
		}
		return benefication.DataCacheValue{}, final
	}

	return benefication.DataCacheValue{
		Age:     a.Age,
		Gender:  g.Gender,
		Country: n.Country,
	}, nil
}
