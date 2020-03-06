package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccMonitorTeam(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		Providers: map[string]terraform.ResourceProvider{
			"sysdig": sysdig.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: monitorTeamWithName(rText()),
			},
			{
				Config: monitorTeamWithOneUser(rText()),
			},
			{
				Config: monitorTeamWithTwoUser(rText()),
			},
			{
				Config: monitorTeamMinimumConfiguration(rText()),
			},
		},
	})
}

func monitorTeamWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample" {
  name               = "sample-%s"
  description        = "%s"
  scope_by           = "container"
  filter             = "container.image.repo = \"sysdig/agent\""
}`, name, name)
}

func monitorTeamWithOneUser(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+monitor-team@sysdig.com"
}

resource "sysdig_monitor_team" "sample" {
  name               = "sample-%s"
  description        = "%s"
  scope_by           = "container"
  filter             = "container.image.repo = \"sysdig/agent\""
  use_sysdig_capture = false

  user_roles {
    email = sysdig_user.sample.email
    role  = "ROLE_TEAM_EDIT"
  }
}`, name, name)
}

func monitorTeamWithTwoUser(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample1" {
  email      = "terraform-test+monitor-team-1@sysdig.com"
}

resource "sysdig_user" "sample2" {
  email      = "terraform-test+monitor-team-2@sysdig.com"
}

resource "sysdig_monitor_team" "sample" {
  name               = "sample-%s"
  description        = "%s"
  scope_by           = "container"
  filter             = "container.image.repo = \"sysdig/agent\""
  use_sysdig_capture = false

  user_roles {
    email = sysdig_user.sample1.email
    role  = "ROLE_TEAM_EDIT"
  }

  user_roles {
    email = sysdig_user.sample2.email
    role  = "ROLE_TEAM_MANAGER"
  }
}`, name, name)
}

func monitorTeamMinimumConfiguration(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample" {
  name      = "sample-%s"
}`, name)
}
