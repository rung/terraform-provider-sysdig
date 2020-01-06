package secure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) GetTeamsById(id int) (t Teams, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodGet, client.GetTeamsUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	t = TeamsFromJSON(body)

	return
}

func (client *sysdigSecureClient) CreateTeams(tRequest Teams) (t Teams, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPost, client.GetTeamUrl(), tRequest.ToJSON())

	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errors.New(response.Status)
		return
	}

	t = TeamsFromJSON(body)
	return
}

func (client *sysdigSecureClient) UpdateTeams(tRequest Teams) (t Teams, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPut, client.GetTeamsUrl(tRequest.ID), tRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	t = TeamsFromJSON(body)
	return
}

func (client *sysdigSecureClient) DeleteTeams(id int) error {
	response, err := client.doSysdigSecureRequest(http.MethodDelete, client.GetTeamsUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}
	return nil
}

func (client *sysdigSecureClient) GetTeamUrl() string {
	return fmt.Sprintf("%s/api/teams", client.URL)
}

func (client *sysdigSecureClient) GetTeamsUrl(id int) string {
	return fmt.Sprintf("%s/api/teams/%d", client.URL, id)
}
