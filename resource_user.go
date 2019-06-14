package main

import (
	"fmt"

	"github.com/bcampoli/dynatrace/apis/onprem/users"
	"github.com/bcampoli/dynatrace/rest"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user_groups": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	api := users.NewAPI(&rest.Config{Insecure: true, NoProxy: false},
		rest.NewCredentials(config.ClusterUrl, config.ApiToken))

	userConfig := getUserConfig(d)

	_, err := api.Create(userConfig)

	if err != nil {
		return err
	}

	d.SetId(d.Get("user_id").(string))
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	api := users.NewAPI(&rest.Config{Insecure: true, NoProxy: false},
		rest.NewCredentials(config.ClusterUrl, config.ApiToken))

	_, err := api.GetUser(d.Get("user_id").(string))
	if err != nil {
		d.SetId("")
		return err
	}

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	api := users.NewAPI(&rest.Config{Insecure: true, NoProxy: false},
		rest.NewCredentials(config.ClusterUrl, config.ApiToken))

	userConfig := getUserConfig(d)

	_, err := api.Update(userConfig)
	if err != nil {
		return err
	}

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	api := users.NewAPI(&rest.Config{Insecure: true, NoProxy: false},
		rest.NewCredentials(config.ClusterUrl, config.ApiToken))

	_, err := api.Delete(d.Get("user_id").(string))
	if err != nil {
		return err
	}

	return nil
}

func getUserConfig(d *schema.ResourceData) *users.UserConfig {
	groupsInterface := d.Get("user_groups").([]interface{})

	groups := make([]string, len(groupsInterface))

	for i, group := range groupsInterface {
		groups[i] = fmt.Sprint(group)
	}

	userConfig := users.UserConfig{
		ID:        d.Get("user_id").(string),
		EMail:     d.Get("email").(string),
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
		Groups:    groups,
	}

	return &userConfig
}
