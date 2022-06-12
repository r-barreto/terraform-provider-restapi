package internal

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccRestApiCreateWithDefinedID(t *testing.T) {
	id := 5678
	name := "test"

	configString := `
		resource "restapi_call" "test" {
			endpoint  = "http://localhost:8080"
		    custom_id = %d
			create {
				path    = "/api/objects/{id}"
				body    = jsonencode({
					"name": "%s" 
				})
			}
			read {
				path = "/api/objects/{id}"
			}
			delete {
				path = "/api/objects/{id}"
			}
		}
		
		output "restapi_output" {
			value = jsondecode(restapi_call.test.create_output).name
		}
	`
	config := fmt.Sprintf(configString, id, name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("restapi_output", name),
				),
			},
		},
	})
}

func TestAccRestApiCreateWithIDFromResponse(t *testing.T) {
	id := 5678
	name := "test"

	configString := `
		resource "restapi_call" "test" {
			endpoint  = "http://localhost:8080"
			create {
				path    = "/api/objects/{id}"
				id_path = "$.id"
				body    = jsonencode({
					"id": "%d",
					"name": "%s" 
				})
			}
			read {
				path = "/api/objects/{id}"
			}
			delete {
				path = "/api/objects/{id}"
			}
		}
		
		output "restapi_output" {
			value = jsondecode(restapi_call.test.create_output).name
		}
	`
	config := fmt.Sprintf(configString, id, name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("restapi_output", name),
				),
			},
		},
	})
}

func TestAccRestApiCreateWithJsonPath(t *testing.T) {
	id := 5678
	name := "test"

	configString := `
		resource "restapi_call" "test" {
			endpoint  = "http://localhost:8080"
			custom_id = %d
			create {
				path      = "/api/objects/{id}"
				json_path = "$.response[0]"
				body      = jsonencode({
					"response": [{
						"name": "%s"
					}]
				})
			}
			read {
				path = "/api/objects/{id}"
			}
			delete {
				path = "/api/objects/{id}"
			}
		}
		
		output "restapi_output" {
			value = jsondecode(restapi_call.test.create_output).name
		}
	`
	config := fmt.Sprintf(configString, id, name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("restapi_output", name),
				),
			},
		},
	})
}

func TestAccRestApiDocExample(t *testing.T) {
	config := `
		provider "restapi" {}

		locals {
			id = 1234
		}

		resource "restapi_call" "test" {
		  endpoint  = "http://localhost:8080"
		  
		  create {
			path      = "/api/objects/${local.id}"
			id_path   = "$.response[0].id"
			json_path = "$.response[0]"
			body      = jsonencode({
			  "response": [{
				"id": local.id,
				"name": "api-test"
			  }]
			})
		  }
		  
		  read {
			path = "/api/objects/{id}"
		  }
		  
		  delete {
			path = "/api/objects/{id}"
		  }
		}
		
		output "restapi_output_name" {
		  value = jsondecode(restapi_call.test.create_output).name
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("restapi_output_name", "api-test"),
				),
			},
		},
	})
}
