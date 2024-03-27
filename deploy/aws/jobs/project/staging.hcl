variable "image" {
  description = "Application to deploy"
  type        = string
  default     = "ghcr.io/kodmain/thetiptop:latest"
}

job "staging" {
  name = "thetiptop: staging"
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
      name = "staging"
      port = "http"
      provider = "nomad"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.staging.rule=Host(`staging.api.kodmain.run`)",
        "traefik.http.routers.staging.tls.certresolver=le",
        "traefik.http.routers.staging.entrypoints=https",
        "traefik.http.routers.staging.service=staging",
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