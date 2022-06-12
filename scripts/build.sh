#!/bin/bash

NAME="restapi"
VERSION=$(cat version)

echo "building terraform-provider-${NAME}_v${VERSION}"

go build -o "terraform-provider-${NAME}_v${VERSION}"