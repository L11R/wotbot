package wargaming

type Config struct {
	ApplicationID string `long:"application_id" env:"APPLICATION_ID" description:"Wargaming API application_id" required:"yes"`
}
