#!/bin/sh
# Obtient l'adresse DNS interne de l'hôte et l'exporte comme variable d'environnement
export INTERNAL_DNS=$(curl http://169.254.169.254/latest/meta-data/local-hostname)