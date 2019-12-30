package secure

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type SysdigSecureClient interface {
	CreatePolicy(Policy) (Policy, error)
	DeletePolicy(int) error
	UpdatePolicy(Policy) (Policy, error)
	GetPolicyById(int) (Policy, error)

	CreateRule(Rule) (Rule, error)
	GetRuleByID(int) (Rule, error)
	UpdateRule(Rule) (Rule, error)
	DeleteRule(int) error

	CreateNotificationChannel(NotificationChannel) (NotificationChannel, error)
	GetNotificationChannelById(int) (NotificationChannel, error)
	DeleteNotificationChannel(int) error
	UpdateNotificationChannel(NotificationChannel) (NotificationChannel, error)

	CreateUsers(Users) (Users, error)
	GetUsersById(int) (Users, error)
	DeleteUsers(int) error
	UpdateUsers(Users) (Users, error)

	CreateTeams(Teams) (Teams, error)
	GetTeamsById(int) (Teams, error)
	DeleteTeams(int) error
	UpdateTeams(Teams) (Teams, error)
}

func NewSysdigSecureClient(sysdigSecureAPIToken string, url string) SysdigSecureClient {
	return &sysdigSecureClient{
		SysdigSecureAPIToken: sysdigSecureAPIToken,
		URL:                  url,
		httpClient:           http.DefaultClient,
	}
}

type sysdigSecureClient struct {
	SysdigSecureAPIToken string
	URL                  string
	httpClient           *http.Client
}

func (client *sysdigSecureClient) doSysdigSecureRequest(method string, url string, payload io.Reader) (*http.Response, error) {
	request, _ := http.NewRequest(method, url, payload)
	request.Header.Set("Authorization", "Bearer "+client.SysdigSecureAPIToken)
	request.Header.Set("Content-Type", "application/json")

	file, _ := os.OpenFile("/tmp/terraform.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	defer file.Close()
	fmt.Fprintln(file, payload)

	return client.httpClient.Do(request)
}
