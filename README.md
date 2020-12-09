# terraform-backend-hsdp
An extendable HTTP backend implementation for terraform

## Features

* State encryption with AES-256-GCM
* Custom state metadata extraction
* Extensible store: supports S3, more to come

## Overview

The core is based on the excellent original implementation from [bhoriuchi/terraform-backend-http](https://github.com/bhoriuchi/terraform-backend-http)
Goal of this project is to support hosting Terraform state on the HSDP platform with little to no setup required.  Currently we use CF credentials to authenticate access to the backend. These credentials are expected to be in your pipeline already anyway. The best practice for this is to use a CF functional account and not your personal CF credentials.
At some point we may introduce service specific credentials similar to the HSDP Docker Registry provides.

# License
License is MIT
