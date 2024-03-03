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

job "project" {
  datacenters = ["eu-west-3"]
  type = "service"

  group "project" {
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
      name = "project-${var.environment}"
      port = "http"
      provider = "nomad"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.project-${var.environment}.rule=Host(`${var.url}`)",
        "traefik.http.routers.project-${var.environment}.entrypoints=https",
        "traefik.http.routers.project-${var.environment}.service=project-${var.environment}",
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
