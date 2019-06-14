package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DT_API_BASE_URL", ""),
			},
			"api_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DT_API_TOKEN", ""),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dynatrace_user":       resourceUser(),
			"dynatrace_user_group": resourceUserGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := Config{
		ClusterUrl: d.Get("api_url").(string),
		ApiToken:   d.Get("api_token").(string),
	}

	return &config, nil
}
