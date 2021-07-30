# terraform-backend-hsdp
An extendable HTTP backend implementation for terraform

# Features

* Encrypt state at rest with AES-256-GCM
* Extensible store: currently supports S3, more to come
* HSDP UAA integration: use LDAP / functional account credentials for auth
* Allow list support: restrict use of an instance backend to specific accounts

# Overview

The primary goal of this project is to offer storage of [Terraform state](https://www.terraform.io/docs/state/index.html) on the HSDP platform with little to no setup required. 
Currently, we use CF credentials to authenticate access to the backend. 
These credentials are expected to be in your pipeline already. 
The best practice is to use a CF functional account for authentication.
Future iterations may introduce service key credentials similar to the HSDP Docker Registry.

The core is derived from [bhoriuchi/terraform-backend-http](https://github.com/bhoriuchi/terraform-backend-http)

# Install
When self-hosting, you should deploy both the S3 bucket and the application deployment
in a separate space in order to limit who has access. Terraform state will contain operator
level secrets so only operators within your organization should have access.

## Provision an S3 bucket
Create an S3 bucket:
```shell
cf cs hsdp-s3 s3_bucket my-tfstate-bucket
```

## Deploy the service

Use the following `manifest.yml` as an example

```yaml
---
applications:
- name: tfstate
  env:
    TFSTATE_KEY: SecretKeyHereThisIsUsedForEncryption
    TFSTATE_REGIONS: us-east,eu-west
  docker:
    image: philipslabs/terraform-backend-hsdp:v0.1.0
  services:
  - my-tfstate-bucket
  routes:
  - route: my-tfstate.eu1.phsdp.com
  processes:
  - type: web
    instances: 1
    memory: 64M
    disk_quota: 1024M
    health-check-type: port
```

Save this to a `manifest.yml` and make the necessary changes i.e. the appname and routes. Then deploy:

```shell
cf push -f manifest.yml
```
After a few seconds you should have a running backend

## Configuration
| Environment | Description | Required | Default |
|-------------|-------------|----------|---------|
| TFSTATE\_KEY | The encryption key for storage at rest | `Yes` | |
| TFSTATE\_ALLOW\_LIST | Comma separated list of allows users | `No` |`""` (every valid LDAP user can access) |
| TFSTATE\_REGIONS | The HSDP regions to validate LDAP accounts in | `No` | `"us-east,eu-west"` |  

# Usage

### 1. Add a `backend.tf` to your terraform definition containing

```hcl
terraform {
  backend "http" {
    address        = "https://my-tfstate.eu1.phsdp.com/my-state"
    lock_address   = "https://my-tfstate.eu1.phsdp.com/my-state"
    unlock_address = "https://my-tfstate.eu1.phsdp.com/my-state"
  }
}
```

The path of the URL will serve as the `key` to your state, so the value should be exactly the same in `address`, `lock_address` and `unlock_address`

### 2. Initialize Terraform

```shell
terraform init \
  -backend-config="username=YOUR-CF-LOGIN" \
  -backend-config="password=YOUR-CF-PASSWORD"
```

In case you want to use a variable key for your storage you should also specify the values of `address`, `lock_address` and `unlock_address` in the terraform init command. In that
case you can remove these values from `backend.tf` as shown in step 1:

```shell
terraform init \
  -backend-config="username=${username}" \
  -backend-config="password=${password}" \
  -backend-config="address=https://my-tfstate.eu1.phsdp.com/${key}" \
  -backend-config="lock_address=https://my-tfstate.eu1.phsdp.com/${key}" \
  -backend-config="unlock_address=https://my-tfstate.eu1.phsdp.com/${key}"

```

### 3. Plan and apply

```shell
terraform plan
...

terraform apply
...
```

# License
License is MIT
