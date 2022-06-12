package internal

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{

		ConfigureContextFunc: nil,
		ResourcesMap: map[string]*schema.Resource{
			"restapi_call": call(),
		},
	}
}
