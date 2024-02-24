job "monitoring" {
  datacenters = ["eu-west-3"]
  type = "service"

  // Groupe pour Prometheus
  group "prometheus" {
    count = 1

    network {
      port "http" { static = 9090 }
    }

    task "prometheus" {
      driver = "docker"
      resources {
        cpu    = 100
        memory = 128
      }

      artifact {
        source = "https://github.com/kodmain/thetiptop/deploy/jobs/monitoring/prometheus/configuration.yml"
        destination = "local/datasource.yaml"
      }

      config {
        image = "prom/prometheus:latest"
        ports = ["http"]
        volumes = [
          "local/prometheus.yml:/etc/prometheus/prometheus.yml",
        ]
      }

      service {
        name = "prometheus"
        port = "http"
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
  }

  group "node-exporter" {
    count = 1

    network {
      port "http" {
        to = 9100
      }
    }

    task "node-exporter" {
      driver = "docker"

      config {
        image = "prom/node-exporter:latest"
        network_mode = "host"
        ports = ["http"]

        // Activer le mode privilégié
        privileged = false

        // Monter les volumes nécessaires pour accéder aux métriques système
        volumes = [
          "/proc:/host/proc:ro",
          "/sys:/host/sys:ro",
          "/:/rootfs:ro"
        ]

        // Passer des paramètres supplémentaires à Node Exporter
        args = [
          "--path.procfs", "/host/proc",
          "--path.sysfs", "/host/sys",
          "--collector.filesystem.ignored-mount-points", "^/(sys|proc|dev|host|etc)($|/)"
        ]
      }

      resources {
        cpu    = 100
        memory = 128
      }

      service {
        name = "node-exporter"
        port = "http"
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

  }

  // Groupe pour Grafana
  group "grafana" {
    count = 1

    network {
      port "http" { static = 3000 }
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
        source = "https://github.com/kodmain/thetiptop/deploy/jobs/monitoring/grafana/datasource.yml"
        destination = "local/datasource.yaml"
      }
      artifact {
        source = "https://github.com/kodmain/thetiptop/deploy/jobs/monitoring/grafana/aws.json"
        destination = "local/aws.json"
      }
      artifact {
        source = "https://github.com/kodmain/thetiptop/deploy/jobs/monitoring/grafana/node.json"
        destination = "local/node.json"
      }
      artifact {
        source = "https://github.com/kodmain/thetiptop/deploy/jobs/monitoring/grafana/nomad.json"
        destination = "local/nomad.json"
      }

      config {
        image = "grafana/grafana:latest"
        ports = ["http"]
        volumes = [
          "local/datasource.yaml:/etc/grafana/provisioning/datasources/default.yaml",
          "local/aws.json:/etc/grafana/provisioning/dashboards/aws.json",
          "local/node.json:/etc/grafana/provisioning/dashboards/node.json",
          "local/nomad.json:/etc/grafana/provisioning/dashboards/nomad.json"
        ]
      }

      service {
        name = "grafana"
        port = "http"  
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
