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

function defineTask(string name, string image) {
  return {
    driver = "docker"
    resources {
      cpu    = 1000 
      memory = 128  
    }
    config {
      image = image
      ports = ["http"]
    }
  }
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
        "traefik.http.routers.api.rule=Host(`${var.url}`)",
        "traefik.http.routers.api.entrypoints=https",
        "traefik.http.routers.api.service=api",
      ]
    }

    task "api" = defineTask("api", var.image)
  }
}
