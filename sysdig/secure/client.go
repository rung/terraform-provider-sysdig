package secure

import (
	"io"
	"net/http"
)

type SysdigSecureClient interface {
	CreatePolicy(Policy) (Policy, error)
	DeletePolicy(int) error
	UpdatePolicy(Policy) (Policy, error)
	GetPolicyById(int) (Policy, error)

	GetUserRulesFile() (UserRulesFile, error)
	UpdateUserRulesFile(UserRulesFile) (UserRulesFile, error)

	CreateNotificationChannel(NotificationChannel) (NotificationChannel, error)
	GetNotificationChannelById(int) (NotificationChannel, error)
	DeleteNotificationChannel(int) error
	UpdateNotificationChannel(NotificationChannel) (NotificationChannel, error)

	CreatePoliciesPriority(PoliciesPriority) (PoliciesPriority, error)
	UpdatePoliciesPriority(PoliciesPriority) (PoliciesPriority, error)
	GetPoliciesPriority() (PoliciesPriority, error)

	CreateUsers(Users) (Users, error)
	GetUsersById(int) (Users, error)
	DeleteUsers(int) error
	UpdateUsers(Users) (Users, error)
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

	return client.httpClient.Do(request)
}
