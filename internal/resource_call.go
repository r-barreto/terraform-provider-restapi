package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CommonResource struct {
	Path       string
	HttpMethod string
	Body       string
	IdPath     string
	JsonPath   string
	Headers    map[string]interface{}
}

func commonResource(defaultHttpMethod string) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Description: "The path to call. It supports replacement by the tag {id}. For example: /{id}",
				Optional:    true,
				Default:     "/{id}",
			},
			"http_method": {
				Type:        schema.TypeString,
				Description: "The http method to call the api.",
				Optional:    true,
				Default:     defaultHttpMethod,
			},
			"body": {
				Type:        schema.TypeString,
				Description: "The request body as a string. It supports replacement by the tag {id}.",
				Optional:    true,
				ForceNew:    true,
				Default:     "",
			},
			"id_path": {
				Type:        schema.TypeString,
				Description: "The json path from where the id can retrieved in a success call. For example: \"$.response.id\"",
				Optional:    true,
				Default:     nil,
			},
			"json_path": {
				Type:        schema.TypeString,
				Description: "The json path you want to return. For example: $.objects[0]",
				Optional:    true,
				Default:     nil,
			},
			"headers": {
				Type:        schema.TypeMap,
				Description: "A map of headers to pass to the rest api call",
				Optional:    true,
				Elem:        schema.TypeString,
				Default:     nil,
			},
		},
	}
}

func call() *schema.Resource {
	return &schema.Resource{
		CreateContext: createCall,
		UpdateContext: updateCall,
		DeleteContext: deleteCall,
		ReadContext:   readCall,
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Description: "A endpoint where the terraform provider will point to, this must include the http(s) schema and port number.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"custom_id": {
				Description: "An ID to manage the resource. If not provided, the id_path must be informed in the create call.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"create": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A element to define how to create the resource",
				MaxItems:    1,
				Elem:        commonResource("POST"),
			},
			"read": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A element to define how to read the resource",
				MaxItems:    1,
				Elem:        commonResource("GET"),
			},
			"update": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A element to define how to update the resource.",
				MaxItems:    1,
				Elem:        commonResource("PUT"),
			},
			"delete": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A element to define how to delete the resource.",
				MaxItems:    1,
				Elem:        commonResource("DELETE"),
			},
			"create_output": {
				Type:        schema.TypeString,
				Description: "The output generated when the resource was created. Use jsondecode to decode this output if it is a JSON object.",
				Computed:    true,
			},
			"raw_output": {
				Type:        schema.TypeString,
				Description: "The output from the last operation.",
				Computed:    true,
			},
		},
	}
}

func createCall(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	id := data.Get("custom_id").(string)
	endpoint := data.Get("endpoint").(string)
	create := data.Get("create").([]interface{})

	diags := convertAndCall(create, id, endpoint, data)

	if err := data.Set("create_output", data.Get("raw_output")); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readCall(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	id := data.Id()
	endpoint := data.Get("endpoint").(string)
	read := data.Get("read").([]interface{})

	return convertAndCall(read, id, endpoint, data)
}

func updateCall(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	id := data.Id()
	endpoint := data.Get("endpoint").(string)
	update := data.Get("update").([]interface{})

	return convertAndCall(update, id, endpoint, data)
}

func deleteCall(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	id := data.Id()
	endpoint := data.Get("endpoint").(string)
	deleteResource := data.Get("delete").([]interface{})

	return convertAndCall(deleteResource, id, endpoint, data)
}

func convertAndCall(resource []interface{}, id string, endpoint string, data *schema.ResourceData) diag.Diagnostics {
	if len(resource) > 0 {
		resourceMap := resource[0].(map[string]interface{})
		commonResource := CommonResource{
			Body:       resourceMap["body"].(string),
			IdPath:     resourceMap["id_path"].(string),
			Path:       resourceMap["path"].(string),
			HttpMethod: resourceMap["http_method"].(string),
			JsonPath:   resourceMap["json_path"].(string),
			Headers:    resourceMap["headers"].(map[string]interface{}),
		}
		return runCall(commonResource, id, endpoint, data)
	} else {
		return runCall(CommonResource{}, id, endpoint, data)
	}
}

func runCall(commonResource CommonResource, id string, endpoint string, data *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if id == "" && commonResource.IdPath == "" {
		return diag.FromErr(fmt.Errorf("both id and create.id_path are empty. Please specify one"))
	}

	response, err := Call(endpoint, commonResource.Path, commonResource.HttpMethod, commonResource.Body, id, commonResource.Headers)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = data.Set("raw_output", response); err != nil {
		return diag.FromErr(err)
	}

	responseMap, _ := toMap(*response)

	if commonResource.IdPath != "" {
		id, err := jsonpath.Get(commonResource.IdPath, responseMap)
		idString := fmt.Sprintf("%v", id)
		if err != nil || idString == "" {
			return diag.FromErr(fmt.Errorf("error querying the id path. ID path: %s, error: %s", commonResource.IdPath, err))
		}

		data.SetId(idString)
	} else if id != "" {
		data.SetId(id)
	} else {
		return diag.FromErr(fmt.Errorf("both id and idPath are empty for the %s operation", commonResource.HttpMethod))
	}

	if commonResource.JsonPath != "" {
		responseJson, err := jsonpath.Get(commonResource.JsonPath, responseMap)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error querying the json path. JSON path: %s, error: %s", commonResource.JsonPath, err))
		}
		responseString, err := toString(responseJson)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error converting JSON response to string. JSON path: %s, error: %s", commonResource.JsonPath, err))
		}

		if err = data.Set("raw_output", responseString); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func toMap(response string) (map[string]interface{}, error) {
	if response == "" {
		return nil, fmt.Errorf("cannot convert to map. Received empty string")
	}
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(response), &jsonMap); err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func toString(response any) (string, error) {
	if responseString, err := json.Marshal(response); err != nil {
		return "", err
	} else {
		return string(responseString), nil
	}
}
