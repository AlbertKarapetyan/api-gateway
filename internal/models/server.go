package models

import (
	"net/url"
)

type Server struct {
	URL     *url.URL
	IsAlive bool
}
