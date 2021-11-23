#!/usr/bin/env bash

go build .
mv terraform-provider-typesense ~/.terraform.d/plugins/terraform-provider-typesense_v0.0.0-snapshot
