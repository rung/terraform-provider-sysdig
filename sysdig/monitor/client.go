package monitor

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type SysdigMonitorClient interface {
	CreateAlert(Alert) (Alert, error)
	DeleteAlert(int) error
	UpdateAlert(Alert) (Alert, error)
	GetAlertById(int) (Alert, error)

	CreateTeam(Team) (Team, error)
	GetTeamById(int) (Team, error)
	DeleteTeam(int) error
	UpdateTeam(Team) (Team, error)
}

func NewSysdigMonitorClient(apiToken string, url string) SysdigMonitorClient {
	return &sysdigMonitorClient{
		SysdigMonitorAPIToken: apiToken,
		URL:                   url,
		httpClient:            http.DefaultClient,
	}
}

type sysdigMonitorClient struct {
	SysdigMonitorAPIToken string
	URL                   string
	httpClient            *http.Client
}

func (c *sysdigMonitorClient) doSysdigMonitorRequest(method string, url string, payload io.Reader) (*http.Response, error) {
	request, _ := http.NewRequest(method, url, payload)
	request.Header.Set("Authorization", "Bearer "+c.SysdigMonitorAPIToken)
	request.Header.Set("Content-Type", "application/json")

	file, _ := os.OpenFile("/tmp/terraform.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	defer file.Close()
	fmt.Fprintf(file, url+":")
	fmt.Fprintln(file, payload)

	return c.httpClient.Do(request)
}
