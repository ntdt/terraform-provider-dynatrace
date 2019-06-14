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
