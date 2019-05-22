package dns

import "os"

var (
	cfAPIKey string = ""
	cfEmail  string = ""
	cfZone   string = ""
)

//InitConfig reads environment variables to initialize config
func InitConfig() {
	cfAPIKey = os.Getenv("CF_API_KEY")
	cfEmail = os.Getenv("CF_EMAIL")
	cfZone = os.Getenv("CF_ZONE")
}
