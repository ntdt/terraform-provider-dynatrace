# Dynatrace Terraform Provider

This is a simple terraform provider for Dynatrace Managed (on premise version of Dynatrace). It allows you to create users and user-groups, and assign users to those groups.

## Example usage:

```
locals {
  env_id = "blahblah"
}

resource "dynatrace_user" "terraform_user" {
  user_id    = "terraform-test-user123"
  first_name = "billy"
  last_name  = "campoli"
  email      = "fake-email@fake.com"

  user_groups = [
    "${dynatrace_user_group.test_group.id}",
  ]
}

resource "dynatrace_user_group" "test_group" {
  name   = "terraform-group1"
  viewer = ["${local.env_id}"]
}

resource "dynatrace_user_group" "test_group2" {
  name            = "terraform-group2"
  manage_settings = ["${local.env_id}"]
}
```

## Build provider:

`go build -o terraform-provider-dynatrace`

## Provider Setup:

Set the following environment variables:

1) `DT_API_TOKEN`: cluster api token to talk to managed cluster.
2) `API_BASE_URL`: URL to access dynatrace managed cluster.