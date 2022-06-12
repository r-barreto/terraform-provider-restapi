# Restapi Terraform Provider


This Terraform provider allows management of a **Rest API** resource.

## Using the provider

Download a binary for your system from the release page and remove the `-os-arch` details, so you're left with `terraform-provider-restapi`.
Use `chmod +x` to make it executable and then either place it at the root of your Terraform folder or in the Terraform plugin folder on your system.


### Example

```terraform
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
```

## Development Guide

### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 1.x
-	[Go](https://golang.org/doc/install) 1.18+
     - correctly setup [GOPATH](http://golang.org/doc/code.html#GOPATH
     - add `$GOPATH/bin` to your `$PATH`
- clone this repository

### Building the provider

To build the provider, run `make build`.

```sh
$ make build
```

### Testing

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of acceptance tests, run `make acceptance`.

```sh
$ make acceptance
```

Alternatively, you can manually start a fake server, run the acceptance tests and then shut down the server.

```sh
$ make start-fakeserver
$ make acceptance
$ make stop-fakeserver
```