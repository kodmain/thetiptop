from diagrams import Cluster, Diagram
from diagrams.aws.compute import EC2
from diagrams.aws.database import RDS
from diagrams.aws.network import NLB, Route53, VPC, CloudFront

with Diagram("TheTipTop - Architecture", filename="architecture", show=False, outformat="png"):

    dns = Route53("DNS")
    cdn = CloudFront("CDN")

    with Cluster("VPC"):  # Le VPC englobe tous les éléments internes
        with Cluster("Private subnet"):
            lb = NLB("LoadBalancer")
            with Cluster("AutoScalingGroup"):
                workers = [EC2("worker1"),
                           EC2("worker2"),
                           EC2("worker3"),
                           EC2("worker4"),
                           EC2("worker5")]

        with Cluster("Isolated subnet"):  # Le sous-réseau isolé pour les bases de données
            with Cluster("DBCluster"):
                db_write = RDS("Databases Write")
                db_read = [RDS("Databases Read")]

    # Définition des flux
    dns >> cdn >> lb  # CDN se connecte directement au LoadBalancer
    lb >> workers  # LoadBalancer distribue le trafic vers les workers
    workers >> db_write  # Les workers interagissent directement avec la base de données principale
    db_write >> db_read  # La base de données Write réplique les données vers les Read
