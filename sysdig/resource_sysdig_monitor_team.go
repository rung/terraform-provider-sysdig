package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigMonitorTeam() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigMonitorTeamCreate,
		Update: resourceSysdigMonitorTeamUpdate,
		Read:   resourceSysdigMonitorTeamRead,
		Delete: resourceSysdigMonitorTeamDelete,

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
			"use_sysdig_capture": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"use_custom_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"use_aws_metrics": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"user_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},

						"role": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ROLE_TEAM_STANDARD",
						},
					},
				},
			},
			"entry_point": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Explore",
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

func resourceSysdigMonitorTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigMonitorClient

	team := monitorTeamFromResourceData(d)

	team, err := client.CreateTeam(team)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(team.ID))
	d.Set("version", team.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigMonitorTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigMonitorClient

	id, _ := strconv.Atoi(d.Id())
	t, err := client.GetTeamById(id)

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
	d.Set("use_sysdig_capture", t.CanUseSysdigCapture)
	d.Set("use_custom_events", t.CanUseSysdigCapture)
	d.Set("use_aws_metrics", t.CanUseSysdigCapture)
	d.Set("entry_point", t.EntryPoint.Module)
	d.Set("default_team", t.DefaultTeam)
	d.Set("user_roles", t.UserRoles)

	return nil
}

func resourceSysdigMonitorTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigMonitorClient

	t := monitorTeamFromResourceData(d)

	t.Version = d.Get("version").(int)
	t.ID, _ = strconv.Atoi(d.Id())

	_, err := client.UpdateTeam(t)

	return err
}

func resourceSysdigMonitorTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigMonitorClient

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteTeam(id)
}

func monitorTeamFromResourceData(d *schema.ResourceData) monitor.Team {
	t := monitor.Team{
		Theme:               d.Get("theme").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		ScopeBy:             d.Get("scope_by").(string),
		Filter:              d.Get("filter").(string),
		CanUseSysdigCapture: d.Get("use_sysdig_capture").(bool),
		CanUseCustomEvents:  d.Get("use_custom_events").(bool),
		CanUseAwsMetrics:    d.Get("use_aws_metrics").(bool),
		EntryPoint: monitor.EntryPoint{
			Module: d.Get("entry_point").(string),
		},
		DefaultTeam: d.Get("default_team").(bool),
		Products:    []string{"SDC"},
	}

	userRoles := []monitor.UserRoles{}
	for _, userRole := range d.Get("user_roles").(*schema.Set).List() {
		ur := userRole.(map[string]interface{})
		userRoles = append(userRoles, monitor.UserRoles{
			Email: ur["email"].(string),
			Role:  ur["role"].(string),
		})
	}
	t.UserRoles = userRoles

	return t
}
