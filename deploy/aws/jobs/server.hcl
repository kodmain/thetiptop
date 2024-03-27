job "server" {
  datacenters = ["eu-west-3"]
  type = "service"

  group "reverse-proxy" {
    count = 1

    network {
      port "http" { static = 80 }
      port "unavailable" { to = 80 }
      port "https" { static = 443 }
      port "traefik" { static = 8080 }
    }

    task "unavailable" {
      driver = "docker"
      resources {
        cpu    = 100
        memory = 10
      }

      service {
        name = "unavailable"
        port = "unavailable"
        provider = "nomad"
        tags = [
            "traefik.enable=true",
            "traefik.http.routers.unavailable.rule=HostRegexp(`{host:.+}`)",
            "traefik.http.routers.unavailable.entrypoints=https",
            "traefik.http.routers.unavailable.priority=1",
        ]
      }
      
      config {
        image = "hashicorp/http-echo:latest"
        ports = ["unavailable"]
        args = [
            "-text=<h1>503 Service Unavailable</h1>",
            "-listen=:80",
            "-status-code=503"
        ]
      }
    }

    task "traefik" {
      driver = "docker"
      resources {
        cpu    = 250
        memory = 128
      }

      service {
        name = "traefik"
        port = "http"
        provider = "nomad" # or "consul"
        
        tags = [
          # Common
          "traefik.enable=true",
          "traefik.http.routers.traefik.tls.certresolver=le",
          "traefik.http.routers.traefik.tls.domains[0].main=kodmain.run",
          "traefik.http.routers.traefik.tls.domains[0].sans=*.kodmain.run",

          # Traefik
          "traefik.http.routers.traefik.rule=Host(`traefik.kodmain.run`)",
          "traefik.http.routers.traefik.entrypoints=https",
          "traefik.http.routers.traefik.service=traefik",
          "traefik.http.services.traefik.loadbalancer.server.port=8080",

          # Nomad
          "traefik.http.routers.nomad.rule=Host(`nomad.kodmain.run`)",
          "traefik.http.routers.nomad.entrypoints=https",
          "traefik.http.routers.nomad.service=nomad",
          "traefik.http.services.nomad.loadbalancer.server.port=4646",
        ]
      }
      
      config {
        image = "traefik:latest"
        ports = ["http", "https", "traefik"]
        volumes = [
          "/home/ec2-user/ssl:/home/ec2-user"
        ]
        args = [
          # API
          "--api.insecure=true",
          "--api.dashboard=true",
          
          # Nomad
          "--providers.nomad=true",
          "--providers.nomad.endpoint.address=http://172.17.0.1:4646",
          "--providers.nomad.endpoint.token=NOMADTOKEN",
          "--providers.nomad.exposedByDefault=false",

          # Entrypoints
          "--entrypoints.traefik.address=:8080",
          "--entryPoints.http.address=:80",
          "--entrypoints.http.http.redirections.entrypoint.to=https",
          "--entrypoints.http.http.redirections.entrypoint.scheme=https",
          "--entryPoints.https.address=:443",

          # TLS
          "--entrypoints.https.http.tls=true",

          # ACME
          "--certificatesresolvers.le.acme.email=contact@kodmain.io",
          "--certificatesresolvers.le.acme.storage=/home/ec2-user/acme.json",
          "--certificatesresolvers.le.acme.dnschallenge=true",
          "--certificatesresolvers.le.acme.dnschallenge.provider=route53",
          "--log.level=DEBUG",
        ]
      }
    }
  }
}