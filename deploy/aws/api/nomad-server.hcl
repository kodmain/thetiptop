# Spécifie le datacenter de cet agent Nomad
datacenter = "eu-west-3"

# Access Control Lists (ACL)
acl {
  enabled = true
}

# Configuration du serveur
server {
    enabled = true
    bootstrap_expect = 1 # Indique que le cluster est complet avec un seul serveur
}

# Active le client
client {
    enabled = true
    cpu_total_compute = 5000 # Limite le nombre de MHz que Nomad peut utiliser
    options = {
      "docker.volumes.enabled" = "true"
    }
}

# Configuration des adresses
addresses {
    http = "0.0.0.0" # Écoute sur toutes les adresses IP pour l'API HTTP
    rpc  = "0.0.0.0" # Écoute sur toutes les adresses IP pour les communications RPC internes
    serf = "0.0.0.0" # Écoute sur toutes les adresses IP pour la communication Serf (utilisé pour le clustering)
}

# Configuration des ports
ports {
    http = 4646 # Port pour l'API HTTP
    rpc  = 4647 # Port pour les communications RPC
    serf = 4648 # Port pour la communication Serf
}

telemetry {
  collection_interval = "1s"
  disable_hostname = true
  prometheus_metrics = true
  publish_allocation_metrics = true
  publish_node_metrics = true
}

# Journalisation
log_level = "INFO"

# Répertoire des données
data_dir = "/home/ec2-user"