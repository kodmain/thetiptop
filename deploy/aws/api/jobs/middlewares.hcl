variable "grafana_admin_user" {
  description = "Admin username for Grafana"
  type        = string
  default     = "admin" 
}

variable "grafana_admin_password" {
  description = "Admin password for Grafana"
  type        = string
}

job "middlewares" {
  datacenters = ["eu-west-3"]
  type = "service"

  group "private" {
    count = 1

    network {
      port "node-exporter" { static = 9100 }
      port "prometheus" { static = 9090 }
    }

    task "prometheus" {
      driver = "docker"
      resources {
        cpu    = 100
        memory = 127
      }

      artifact {
        source = "https://raw.githubusercontent.com/kodmain/thetiptop/main/deploy/aws/api/jobs/monitoring/prometheus/prometheus.yml"
        destination = "local/prometheus"
      }

      config {
        image = "prom/prometheus:latest"
        ports = ["prometheus"]
        args = [ 
          "--config.file=/etc/prometheus/prometheus.yml",
          "--web.external-url=https://internal.kodmain.run/prometheus",
          "--web.route-prefix=/prometheus"
        ]
        volumes = [
          "local/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml",
        ]
      }

      service {
        name = "prometheus"
        port = "prometheus"
        provider = "nomad"
        tags = [
            "traefik.enable=true",
            "traefik.http.routers.prometheus.rule=Host(`internal.kodmain.run`) && PathPrefix(`/prometheus`)",
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
          "traefik.http.routers.node-exporter-internal.rule=Host(`internal.kodmain.run`) && PathPrefix(`/node-exporter`)",
          "traefik.http.routers.node-exporter-internal.entrypoints=https",
          "traefik.http.routers.node-exporter-internal.middlewares=node-exporter-stripprefix",
          "traefik.http.middlewares.node-exporter-stripprefix.stripprefix.prefixes=/node-exporter",
          "traefik.http.services.node-exporter.loadbalancer.server.port=9100" 
        ]
      }
    }

  }

  group "public" {
    count = 1

    network {
      port "grafana" { static = 3000 }
    }

    task "grafana" {
      driver = "docker"

      env {
        GF_AUTH_ANONYMOUS_ENABLED = "true"
        GF_AUTH_ANONYMOUS_ORG_NAME = "Main Org."
        GF_AUTH_ANONYMOUS_ORG_ROLE = "Viewer"

        GF_SECURITY_ADMIN_USER = "${var.grafana_admin_user}"
        GF_SECURITY_ADMIN_PASSWORD = "${var.grafana_admin_password}"
        GF_LOG_MODE="console"
        GF_PATHS_PROVISIONING="/etc/grafana/provisioning"
        GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH="/var/lib/grafana/dashboards/thetiptop.json"
      }

      resources {
        cpu    = 100
        memory = 128
      }

      artifact {
        source = "https://raw.githubusercontent.com/kodmain/thetiptop/main/deploy/aws/api/jobs/monitoring/grafana/datasources/thetiptop.yml"
        destination = "local/datasources"
      }

      artifact {
        source = "https://raw.githubusercontent.com/kodmain/thetiptop/main/deploy/aws/api/jobs/monitoring/grafana/dashboards/thetiptop.yml"
        destination = "local/dashboards"
      }

      artifact {
        source = "https://raw.githubusercontent.com/kodmain/thetiptop/main/deploy/aws/api/jobs/monitoring/grafana/templates/thetiptop.json"
        destination = "local/template"
      }

      config {
        image = "grafana/grafana:latest"
        ports = ["grafana"]
        volumes = [
          "local/datasources:/etc/grafana/provisioning/datasources",
          "local/dashboards:/etc/grafana/provisioning/dashboards",
          "local/templates:/var/lib/grafana/dashboards"
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