# terraform-backend-hsdp
An extendable HTTP backend implementation for terraform

## Features

* State encryption with AES-256-GCM
* Extensible store: currently supports S3, more to come

## Overview

The primary goal of this project is to offer storage of [Terraform state](https://www.terraform.io/docs/state/index.html) on the HSDP platform with little to no setup required. 
Currently, we use CF credentials to authenticate access to the backend. 
These credentials are expected to be in your pipeline already. 
The best practice is to use a CF functional account for authentication.
Future iterations may introduce service key credentials similar to the HSDP Docker Registry.

The core is derived from [bhoriuchi/terraform-backend-http](https://github.com/bhoriuchi/terraform-backend-http)

## Usage

### 1. Add a `backend.tf` to your terraform definition containing

```hcl
terraform {
  backend "http" {
    address = "https://tfstate.eu1.phsdp.com/my-state"
    lock_address = "https://tfstate.eu1.phsdp.com/my-state"
    unlock_address = "https://eu1.phsdp.com/my-state"
  }
}
```

The path of the URL will serve as the `key` to you state, so the value should be exactly the same in `address`, `lock_address` and `unlock_address`

### 2. Initialize Terraform

```shell
terraform init \
  -backend-config="username=YOUR-CF-LOGIN" \
  -backend-config="password=YOUR-CF-PASSWORD"
```

In case you want to use a variable key for your storage you should also specify the values of `address`, `lock_address` and `unlock_address` in the terraform init command. In that
case you can remove these values from `backend.tf` as shown in step 1.

### 3. Plan and apply

```shell
terraform plan
...

terraform apply
...
```

# License
License is MIT
