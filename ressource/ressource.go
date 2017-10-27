package ressource

import (
	ui "github.com/gizak/termui"
	"time"
	formatTime "github.com/jehiah/go-strftime"
)

type ConfigurationFile struct {
	Amqp        Amqp                    `json:"amqp"`
	PathConsole string                  `json:"console_path"`
	Sentry      Sentry                  `json:"sentry"`
	Redis       Redis                   `json:"redis"`
	HttpStatus  map[string]HttpEndpoint `json:"http-status"`
}

type Sentry struct {
	Host             string `json:"host"`
	ApiKey           string `json:"api_key"`
	ProjectSlug      string `json:"project_slug"`
	OrganisationSlug string `json:"organization_slug"`
	MaxItem          int    `json:"max_item"`
}

type Amqp struct {
	Host     string           `json:"host"`
	Port     string           `json:"port"`
	WebPort  string           `json:"web_port"`
	User     string           `json:"user"`
	Password string           `json:"password"`
	Queue    map[string]Queue `json:"queue"`
}

type Queue struct {
	AlertWarningMessage int `json:"alert_warning_message"`
	AlertCriticMessage  int `json:"alert_critical_message"`
}

type Redis struct {
	Host string   `json:"host"`
	Info []string `json:"info"`
}

type HttpEndpoint struct {
	Url               string `json:"url"`
	BasicAuthUser     string `json:"basic_auth_user"`
	BasicAuthPassword string `json:"basic_auth_password"`
	HttpCodeExpected  int    `json:"http_code_expected"`
}

const (
	CriticalLevel = 1
	WarningLevel  = 2
	CriticalColor = ui.ColorRed
	WarningColor  = ui.ColorYellow
)

func FormatDate(date time.Time) string {
	return formatTime.Format("%Y-%m-%d %H:%M:%S", date)
}

func ErrorToSliceString(errs []error) []string {
	errorsStr := []string{}
	for _, err := range errs {
		errorsStr = append(errorsStr, err.Error())
	}

	return errorsStr
}