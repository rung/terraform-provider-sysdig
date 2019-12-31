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
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memberships": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ROLE_TEAM_READ",
						},
					},
				},
			},
			"default_team": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	d.Set("filter", t.Filter)
	d.Set("memberships", t.Memberships)
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
		Filter:      d.Get("filter").(string),
		DefaultTeam: d.Get("default_team").(bool),
	}

	memberships := []secure.Memberships{}
	for _, membership := range d.Get("memberships").(*schema.Set).List() {
		ms := membership.(map[string]interface{})
		var m secure.Memberships
		m.UserId = ms["user_id"].(string)
		m.Role = ms["role"].(string)
		memberships = append(memberships, m)
	}

	t.Memberships = memberships

	return t
}
