package internal

import "encoding/base64"

type Config struct {
	mosquittoPid    int
	pskFilePath     string
	basicAuthHeader string
	port            string
	crt             string
	key             string
	aclFile         string
}

func NewConfig(mosquittoPid int, pskFilePath string,
	basicAuthLogin string,
	basicAuthPass string,
	port string,
	crt string,
	key string,
	aclFile string) Config {
	if basicAuthLogin != "" && basicAuthPass != "" {
		return Config{
			mosquittoPid:    mosquittoPid,
			pskFilePath:     pskFilePath,
			basicAuthHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte(basicAuthLogin+":"+basicAuthPass)),
			port:            ":" + port,
			crt:             crt,
			key:             key,
			aclFile:         aclFile,
		}
	} else {
		return Config{
			mosquittoPid:    mosquittoPid,
			pskFilePath:     pskFilePath,
			basicAuthHeader: "",
			port:            ":" + port,
			crt:             crt,
			key:             key,
			aclFile:         aclFile,
		}
	}
}

func (c *Config) isTLS() bool {
	return c.crt != "" && c.key != ""
}
