package secure_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func TestCreateTeams(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	team, err := sysdigSecureClient.CreateTeams(aTeams())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteTeams(team.ID)
	assert.NotEqual(t, 0, team.ID)
}

func TestUpdateTeams(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
	team, err := sysdigSecureClient.CreateTeams(aTeams())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteUsers(team.ID)
	assert.Equal(t, "test-team", team.Name)

	team.Name = "Changed Name"
	newTeam, err := sysdigSecureClient.UpdateTeams(team)

	assert.Nil(t, err)
	assert.Equal(t, "Changed Name", newTeam.Name)
}

func TestUpdateTeamsFailsWhenDoesNotExist(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	team, err := sysdigSecureClient.CreateTeams(aTeams())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteUsers(team.ID)

	nonExistentID := 0
	team.ID = nonExistentID
	team.Name = "Changed Name"
	_, err = sysdigSecureClient.UpdateTeams(team)

	assert.NotNil(t, err)
}

func TestGetTeamsById(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
	team, err := sysdigSecureClient.CreateTeams(aTeams())
	assert.Nil(t, err)

	newTeam, err := sysdigSecureClient.GetTeamsById(team.ID)

	assert.Nil(t, err)
	assert.Equal(t, team.ID, newTeam.ID)
	assert.Equal(t, team.Version, newTeam.Version)
	assert.Equal(t, team.Name, newTeam.Name)

	defer sysdigSecureClient.DeleteUsers(team.ID)
}

func TestGetTeamsByIdFailsWhenDoesNotExist(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	nonExistentID := 0
	_, err := sysdigSecureClient.GetTeamsById(nonExistentID)

	assert.NotNil(t, err)
}

func aTeams() secure.Teams {
	return secure.Teams{
		Name:    "test-team",
		Theme:   "#73A1F7",
		ScopeBy: "container",
	}
}
