variable "url" {
  description = "URL to preview the front"
  type        = string
  default     = "https://api.kodmain.run"
}

variable "image" {
  description = "Application to deploy"
  type        = string
  default     = "ghcr.io/kodmain/thetiptop"
}

variable "environment" {
  description = "Environment to use (sandbox, staging, production)"
  type        = string
  default     = "sandbox"
}

job "${var.environment}" {
  assert(in(var.environment, ["sandbox", "staging", "production"]), "Environnement invalide.")

  datacenters = ["eu-west-3"]
  type = "service"

  group "project" {
    count = 1

    network {
      port "http" { to = 80 }
    }

    service {
      name = "project"
      port = "http"
      provider = "nomad"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.project-${var.environment}.rule=Host(`${var.url}`)",
        "traefik.http.routers.project-${var.environment}.entrypoints=https",
        "traefik.http.routers.project-${var.environment}.service=project",
      ]
    }

    task "project-${var.environment}" {
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
