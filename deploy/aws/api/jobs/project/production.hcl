variable "image" {
  description = "Application to deploy"
  type        = string
  default     = "ghcr.io/kodmain/thetiptop:latest"
}

job "production" {
  name = "thetiptop: production"
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
      name = "production"
      port = "http"
      provider = "nomad"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.production.rule=Host(`api.kodmain.run`)",
        "traefik.http.routers.production.entrypoints=https",
        "traefik.http.routers.production.service=production",
      ]
    }

    task "api" {
      driver = "docker"
      resources {
        cpu    = 1000 
        memory = 128  
      }
      config {
        force_pull = true
        image = var.image
        ports = ["http"]
      }
    }
  }
}