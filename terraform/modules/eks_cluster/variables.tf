variable "name" {
  description = "The name to be associated with this module's resources"
}

variable "environment" {
  description = "The environment value (dev or prod)"
}

variable "subnets" {
  description = "The subnets to be associated with the EKS cluster"
  type        = list(string)
}