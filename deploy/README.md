# Documentation pour le déploiment

## Description
Les images générer par le workflow sont sur une base "scratch" prêt à l'usage
$ `docker run --pid=host -p 8000:80 thetiptop:latest`

## Attention
Le `--pid=host` est très importants il permet de pouvoir exécuter fiber dans docker en utilisant profitant des bénéfice du multithreading 