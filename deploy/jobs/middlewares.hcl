job "middlewares" {
  datacenters = ["eu-west-3"]
  type = "service"

  group "metrics" {
    count = 1

    network {
      port "node-exporter" { to = 9100 }
      port "prometheus" { static = 9090 }
      port "grafana" { static = 3000 }
    }

    task "prometheus" {
      driver = "docker"
      resources {
        cpu    = 100
        memory = 128
      }

      artifact {
        source = "https://raw.githubusercontent.com/kodmain/thetiptop/main/deploy/jobs/monitoring/prometheus/configuration.yml"
        destination = "local/prometheus"
      }

      config {
        image = "prom/prometheus:latest"
        ports = ["prometheus"]
        volumes = [
          "local/prometheus/configuration.yml:/etc/prometheus/prometheus.yml",
        ]
      }

      service {
        name = "prometheus"
        port = "prometheus"
        provider = "nomad"
        tags = [
            "traefik.enable=true",
            "traefik.http.routers.prometheus.rule=Host(`prometheus.kodmain.run`)",
            "traefik.http.routers.prometheus.entrypoints=https",
            "traefik.http.routers.prometheus.service=prometheus",
            "traefik.http.services.prometheus.loadbalancer.server.port=9090",
        ]
      }
    }

    task "node-exporter" {
      driver = "docker"

      config {
        image = "prom/node-exporter:latest"
        network_mode = "host"
        ports = ["node-exporter"]

        volumes = [
          "/proc:/host/proc:ro",
          "/sys:/host/sys:ro",
          "/:/rootfs:ro"
        ]

        args = [
          "--path.procfs", "/host/proc",
          "--path.sysfs", "/host/sys",
          "--collector.filesystem.ignored-mount-points", "^/(sys|proc|dev|host|etc)($|/)"
        ]
      }

      resources {
        cpu    = 100
        memory = 64
      }

      service {
        name = "node-exporter"
        port = "node-exporter"
        provider = "nomad"

        tags = [
          "traefik.enable=true",
          "traefik.http.routers.node-exporter.rule=Host(`node-exporter.kodmain.run`)",
          "traefik.http.routers.node-exporter.entrypoints=https",
          "traefik.http.routers.node-exporter.service=node-exporter",
          "traefik.http.services.node-exporter.loadbalancer.server.port=9100",
        ]
      }
    }

    task "grafana" {
      driver = "docker"

      env {
        GF_AUTH_ANONYMOUS_ENABLED = "true"
        GF_AUTH_ANONYMOUS_ORG_NAME = "Main Org."
        GF_AUTH_ANONYMOUS_ORG_ROLE = "Viewer"

        GF_SECURITY_ADMIN_USER = "admin"
        GF_SECURITY_ADMIN_PASSWORD = "admin"
      }

      resources {
        cpu    = 100
        memory = 128
      }

      artifact {
        source = "https://raw.githubusercontent.com/kodmain/thetiptop/main/deploy/jobs/monitoring/grafana/datasource.yml"
        destination = "local/datasource"
      }
      artifact {
        source = "https://raw.githubusercontent.com/kodmain/thetiptop/main/deploy/jobs/monitoring/grafana/dashboard/aws.json"
        destination = "local/dashboard"
      }
      artifact {
        source = "https://raw.githubusercontent.com/kodmain/thetiptop/main/deploy/jobs/monitoring/grafana/dashboard/node.json"
        destination = "local/dashboard"
      }
      artifact {
        source = "https://raw.githubusercontent.com/kodmain/thetiptop/main/deploy/jobs/monitoring/grafana/dashboard/nomad.json"
        destination = "local/dashboard"
      }

      config {
        image = "grafana/grafana:latest"
        ports = ["grafana"]
        volumes = [
          "local/datasource:/etc/grafana/provisioning/datasources",
          "local/dashboard:/etc/grafana/provisioning/dashboards",
        ]
      }

      service {
        name = "grafana"
        port = "grafana"  
        provider = "nomad"

        tags = [
            "traefik.enable=true",
            "traefik.http.routers.grafana.rule=Host(`grafana.kodmain.run`)",
            "traefik.http.routers.grafana.entrypoints=https",
            "traefik.http.routers.grafana.service=grafana",
            "traefik.http.services.grafana.loadbalancer.server.port=3000",
        ]
      }
    }
  }
}