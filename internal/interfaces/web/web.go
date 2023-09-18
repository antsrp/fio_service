package web

import "net/http"

type Connector interface {
	Do(*http.Request, interface{}) (int, error)
}
