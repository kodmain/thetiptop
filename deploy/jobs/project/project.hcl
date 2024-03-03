variable "url" {
  description = "URL to preview the front"
  type        = string
  default     = "https://api.kodmain.run"
}

variable "image" {
  description = "Application to deploy"
  type        = string
  default     = "ghcr.io/kodmain/thetiptop:latest"
}

variable "environment" {
  description = "Environment to use (sandbox, staging, production)"
  type        = string
  default     = "sandbox"

  validation {
    condition     = contains(["sandbox", "staging", "production"], var.environment)
    error_message = "The environment must be one of: sandbox, staging, production."
  }
}

job "environment" {
  name = "${var.environment}"
  datacenters = ["eu-west-3"]
  type = "service"

  group "thetiptop" {
    count = 1

    network {
      port "http" { to = 80 }
    }

    restart {
      attempts = 3
      interval = "30m"
      delay    = "15s"
      mode     = "fail"
    }

    service {
      name = "${var.environment}"
      port = "http"
      provider = "nomad"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.${var.environment}.rule=Host(`${var.url}`)",
        "traefik.http.routers.${var.environment}.entrypoints=https",
        "traefik.http.routers.${var.environment}.service=${var.environment}",
      ]
    }

    task "api" {
      driver = "docker"
      resources {
        cpu    = 1000 
        memory = 128  
      }
      config {
        image = var.image
        ports = ["http"]
      }
    }

  }
}
