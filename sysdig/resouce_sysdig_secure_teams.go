package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigSecureTeams() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigTeamsCreate,
		Update: resourceSysdigTeamsUpdate,
		Read:   resourceSysdigTeamsRead,
		Delete: resourceSysdigTeamsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"theme": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "#73A1F7",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope_by": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "container",
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"advanced_users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default_team": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSysdigTeamsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	teams := teamsFromResourceData(d)

	teams, err := client.CreateTeams(teams)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(teams.ID))
	d.Set("version", teams.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigTeamsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	id, _ := strconv.Atoi(d.Id())
	t, err := client.GetTeamsById(id)

	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("version", t.Version)
	d.Set("theme", t.Theme)
	d.Set("name", t.Name)
	d.Set("description", t.Description)
	d.Set("scope_by", t.ScopeBy)
	d.Set("filter", t.Filter)
	d.Set("advanced_users", t.AdvancedUsers)
	d.Set("default_team", t.DefaultTeam)

	return nil
}

func resourceSysdigTeamsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	t := teamsFromResourceData(d)

	t.Version = d.Get("version").(int)
	t.ID, _ = strconv.Atoi(d.Id())

	_, err := client.UpdateTeams(t)

	return err
}

func resourceSysdigTeamsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteTeams(id)
}

func teamsFromResourceData(d *schema.ResourceData) (u secure.Teams) {
	t := secure.Teams{
		Theme:       d.Get("theme").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ScopeBy:     d.Get("scope_by").(string),
		Filter:      d.Get("filter").(string),
		DefaultTeam: d.Get("default_team").(bool),
	}

	userRoles := []secure.UserRoles{}
	for _, user := range d.Get("advanced_users").([]interface{}) {
		u := user.(string)
		userRoles = append(userRoles, secure.UserRoles{
			UserId: u,
			Role:   "ROLE_TEAM_EDIT",
		})
	}
	t.AdvancedUsers = userRoles

	return t
}
