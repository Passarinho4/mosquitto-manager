package internal

import "encoding/base64"

type Config struct {
	mosquittoPid int
	pskFilePath  string
	basicAuthHeader string
}

func NewConfig(mosquittoPid int, pskFilePath string, basicAuthLogin string, basicAuthPass string) Config {
	return Config{
		mosquittoPid: mosquittoPid,
		pskFilePath:  pskFilePath,
		basicAuthHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte(basicAuthLogin + ":" + basicAuthPass)),
	}
}
