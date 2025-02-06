package models

import (
	"net/url"
)

type Server struct {
	config  Config
	URL     *url.URL
	IsAlive bool
}
