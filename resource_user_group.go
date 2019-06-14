package main

import (
	"fmt"
	"strings"

	ug "github.com/dtcookie/dynatrace/apis/onprem/user_groups"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupCreate,
		Read:   resourceUserGroupRead,
		Update: resourceUserGroupUpdate,
		Delete: resourceUserGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"is_cluster_admin_group": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"viewer": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"manage_settings": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"agent_install": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"log_viewer": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"view_sensitive_request_data": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"configure_request_capture_data": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceUserGroupCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	api := ug.NewAPI(&rest.Config{Insecure: true, NoProxy: false},
		rest.NewCredentials(config.ClusterUrl, config.ApiToken))

	userGroupConfig := getUserGroupConfig(d)

	_, err := api.Create(userGroupConfig)

	if err != nil {
		return err
	}

	d.SetId(strings.ToLower(d.Get("name").(string)))
	return resourceUserGroupRead(d, m)
}

func resourceUserGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	api := ug.NewAPI(&rest.Config{Insecure: true, NoProxy: false},
		rest.NewCredentials(config.ClusterUrl, config.ApiToken))

	_, err := api.Get(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	return nil
}

func resourceUserGroupUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	api := ug.NewAPI(&rest.Config{Insecure: true, NoProxy: false},
		rest.NewCredentials(config.ClusterUrl, config.ApiToken))

	userGroupConfig := getUserGroupConfig(d)
	userGroupConfig.ID = d.Id()

	_, err := api.Update(userGroupConfig)
	if err != nil {
		return err
	}

	return resourceUserGroupRead(d, m)
}

func resourceUserGroupDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	api := ug.NewAPI(&rest.Config{Insecure: true, NoProxy: false},
		rest.NewCredentials(config.ClusterUrl, config.ApiToken))

	err := api.Delete(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func getUserGroupConfig(d *schema.ResourceData) *ug.GroupConfig {
	accessRights := make(map[string][]string)

	if d.Get("viewer") != nil {
		viewerInterface := d.Get("viewer").([]interface{})
		accessRights = populateMapHelper("VIEWER", accessRights, viewerInterface)
	}

	if d.Get("manage_settings") != nil {
		manageSettingsInterface := d.Get("manage_settings").([]interface{})
		accessRights = populateMapHelper("MANAGE_SETTINGS", accessRights, manageSettingsInterface)
	}

	if d.Get("agent_install") != nil {
		agentInstallInterface := d.Get("agent_install").([]interface{})
		accessRights = populateMapHelper("AGENT_INSTALL", accessRights, agentInstallInterface)
	}

	if d.Get("log_viewer") != nil {
		logViewerInterface := d.Get("log_viewer").([]interface{})
		accessRights = populateMapHelper("LOG_VIEWER", accessRights, logViewerInterface)
	}

	if d.Get("view_sensitive_request_data") != nil {
		VSSRInterface := d.Get("view_sensitive_request_data").([]interface{})
		accessRights = populateMapHelper("VIEW_SENSITIVE_REQUEST_DATA", accessRights, VSSRInterface)
	}

	if d.Get("configure_request_capture_data") != nil {
		CRCInterface := d.Get("configure_request_capture_data").([]interface{})
		accessRights = populateMapHelper("CONFIGURE_REQUEST_CAPTURE_DATA", accessRights, CRCInterface)
	}

	groupConfig := ug.GroupConfig{
		Name:                d.Get("name").(string),
		IsClusterAdminGroup: d.Get("is_cluster_admin_group").(bool),
		LDAPGroupNames:      []string{},
		AccessRight:         accessRights,
	}

	return &groupConfig
}

func populateMapHelper(key string, m map[string][]string, config []interface{}) map[string][]string {

	for _, env := range config {
		m[key] = append(m[key], fmt.Sprint(env))
	}

	return m
}
