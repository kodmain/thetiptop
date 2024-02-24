job "project" {
  datacenters = ["eu-west-3"]
  type = "service"

  group "whoami" {
    count = 1

    network {
      port "http" { to = 80 }
    }

    service {
      name = "whoami"
      port = "http"
      provider = "nomad"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.whoami.rule=Host(`whoami.kodmain.run`)",
        "traefik.http.routers.whoami.entrypoints=https",
        "traefik.http.routers.whoami.service=whoami",
      ]
    }

    task "whoami" {
      driver = "docker"
      resources {
        cpu    = 1000
        memory = 10
      }
      config {
        image = "traefik/whoami:latest"
        ports = ["http"]
      }
    }
  }
}
