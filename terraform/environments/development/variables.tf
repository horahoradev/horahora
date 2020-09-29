variable "region" {
  description = "AWS region for the infrastructure"
  default     = "us-west-1"
}

variable "environment" {
  description = "The environment value (dev or prod)"
  default     = "dev"
}

variable "videodb_password" {}

variable "schedulerdb_password" {}

variable "userdb_password" {}



