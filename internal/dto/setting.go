// setting.go
package dto

type SettingDto struct {
	AllowRegistration *bool                  `json:"allowRegistration"`
	FrontendSettings  map[string]interface{} `json:"frontendSettings"`
	SmtpEnabled       *bool                  `json:"smtpEnabled"`
	SmtpDriver        *string                `json:"smtpDriver"`
	SmtpHost          *string                `json:"smtpHost"`
	SmtpPort          *int                   `json:"smtpPort"`
	SmtpUser          *string                `json:"smtpUser"`
	SmtpPass          *string                `json:"smtpPass"`
	SmtpFrom          *string                `json:"smtpFrom"`
	SmtpEncryption    *string                `json:"smtpEncryption"`
	SmtpTLS           *bool                  `json:"smtpTLS"`
	StorageEnabled    *bool                  `json:"storageEnabled"`
	StorageConfig     map[string]interface{} `json:"storageConfig"`
}
