package secure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) GetUsersById(id int) (u Users, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodGet, client.GetUsersUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	u = UsersFromJSON(body)

	return
}

func (client *sysdigSecureClient) CreateUsers(uRequest Users) (u Users, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPost, client.GetUserUrl(), uRequest.ToJSON())

	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errors.New(response.Status)
		return
	}

	u = UsersFromJSON(body)
	return
}

func (client *sysdigSecureClient) UpdateUsers(uRequest Users) (u Users, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPut, client.GetUsersUrl(u.ID), uRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	u = UsersFromJSON(body)
	return

}

func (client *sysdigSecureClient) DeleteUsers(id int) error {
	response, err := client.doSysdigSecureRequest(http.MethodDelete, client.GetUsersUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return errors.New(response.Status)
	}
	return nil
}

func (client *sysdigSecureClient) GetUserUrl() string {
	return fmt.Sprintf("%s/api/users", client.URL)
}

func (client *sysdigSecureClient) GetUsersUrl(id int) string {
	return fmt.Sprintf("%s/api/users/%d", client.URL, id)
}
