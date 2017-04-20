package lib

import (
	"extres"

	"gopkg.in/gomail.v2"
)

// AppConfig is the shared struct of configuration values
//=======================================================================================
var AppConfig extres.ExternalResources // db, email, etc

// ReadConfig will read the configuration file "config.json" if
// it exists in the current directory
//=======================================================================================
func ReadConfig() {
	fname := "config.json"
	err := extres.ReadConfig(fname, &AppConfig)
	Errcheck(err)
}

// SMTPDialAndSend sends a single SMTP email message, m, using the SMTP information in
// the config.json file
//=======================================================================================
func SMTPDialAndSend(m *gomail.Message) error {
	d := gomail.NewDialer(AppConfig.SMTPHost, AppConfig.SMTPPort, AppConfig.SMTPLogin, AppConfig.SMTPPass)
	return d.DialAndSend(m)
}
