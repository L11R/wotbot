package wargaming

import "time"

type Config struct {
	ApplicationID string        `long:"application-id" env:"APPLICATION_ID" description:"Wargaming API application_id" required:"yes"`
	HTTPTimeout   time.Duration `long:"http-timeout" env:"HTTP_TIMEOUT" description:"HTTP Wargaming API call timeout" default:"10s"`
}
