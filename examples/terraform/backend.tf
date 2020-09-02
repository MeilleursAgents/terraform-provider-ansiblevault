terraform {
  required_version = ">= 0.13"

  required_providers {
    ansiblevault = {
      source  = "MeilleursAgents/ansiblevault"
      version = "~> 2.0"
    }
  }
}
