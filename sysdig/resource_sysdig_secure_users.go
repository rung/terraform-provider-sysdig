package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigSecureUsers() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigUsersCreate,
		Update: resourceSysdigUsersUpdate,
		Read:   resourceSysdigUsersRead,
		Delete: resourceSysdigUsersDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"system_role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ROLE_USER",
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSysdigUsersCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	users := usersFromResourceData(d)

	users, err := client.CreateUsers(users)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(users.ID))
	d.Set("version", users.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	id, _ := strconv.Atoi(d.Id())
	u, err := client.GetUsersById(id)

	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("version", u.Version)
	d.Set("system_role", u.SystemRole)
	d.Set("email", u.Email)
	d.Set("first_name", u.FirstName)
	d.Set("last_name", u.LastName)

	return nil
}

func resourceSysdigUsersUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	u := usersFromResourceData(d)

	u.Version = d.Get("version").(int)
	u.ID, _ = strconv.Atoi(d.Id())

	_, err := client.UpdateUsers(u)

	return err
}

func resourceSysdigUsersDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteUsers(id)
}

func usersFromResourceData(d *schema.ResourceData) (u secure.Users) {
	u = secure.Users{
		SystemRole: d.Get("system_role").(string),
		Email:      d.Get("email").(string),
		FirstName:  d.Get("first_name").(string),
		LastName:   d.Get("last_name").(string),
	}
	return u
}
