package secure_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func TestCreateUsers(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	user, err := sysdigSecureClient.CreateUsers(aUsers())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteUsers(user.ID)
	assert.NotEqual(t, 0, user.ID)
}

func TestUpdateUsers(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
	user, err := sysdigSecureClient.CreateUsers(aUsers())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteUsers(user.ID)
	assert.Equal(t, "root@localhost", user.Email)
	assert.Equal(t, "users", user.FirstName)
	assert.Equal(t, "test", user.LastName)

	user.FirstName = "Changed Name"
	newUser, err := sysdigSecureClient.UpdateUsers(user)

	assert.Nil(t, err)
	assert.Equal(t, "Changed Name", newUser.FirstName)
}

func TestUpdateUsersFailsWhenDoesNotExist(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	user, err := sysdigSecureClient.CreateUsers(aUsers())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteUsers(user.ID)

	nonExistentID := 0
	user.ID = nonExistentID
	user.FirstName = "Changed Name"
	_, err = sysdigSecureClient.UpdateUsers(user)

	assert.NotNil(t, err)
}

func TestGetUsersById(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
	user, err := sysdigSecureClient.CreateUsers(aUsers())
	assert.Nil(t, err)

	newUser, err := sysdigSecureClient.GetUsersById(user.ID)

	assert.Nil(t, err)
	assert.Equal(t, user.ID, newUser.ID)
	assert.Equal(t, user.Version, newUser.Version)
	assert.Equal(t, user.FirstName, newUser.FirstName)
	assert.Equal(t, user.LastName, newUser.LastName)

	defer sysdigSecureClient.DeleteUsers(user.ID)
}

func TestGetUsersByIdFailsWhenDoesNotExist(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	nonExistentID := 0
	_, err := sysdigSecureClient.GetUsersById(nonExistentID)

	assert.NotNil(t, err)
}

func aUsers() secure.Users {
	return secure.Users{
		Email: "root@localhost",
		FirstName: "users",
		LastName: "test",
	}
}

