variable "tags" {
  type        = map(string)
  description = "The default tags to use for AWS resources"
  default = {
    App = "vdlg-login"
  }
}

variable "region" {
  type        = string
  description = "The default region to use for AWS"
  default     = "us-east-1"
}

variable "vpc_name" {
  type        = string
  description = "The name of the VPC"
  default     = "vpc-vdlg"
}