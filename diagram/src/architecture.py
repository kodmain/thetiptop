from diagrams import Cluster, Diagram
from diagrams.aws.compute import EC2
from diagrams.aws.database import RDS
from diagrams.aws.network import NLB, Route53, VPC, CloudFront

with Diagram("TheTipTop - Architecture", filename="architecture", show=False, outformat="svg"):
    dns = Route53("DNS")
    cdn = CloudFront("CDN")

    with Cluster("VPC"): 
        with Cluster("Private subnet"):
            lb = NLB("LoadBalancer")
            with Cluster("AutoScalingGroup"):
                workers = [EC2("worker1"),
                           EC2("worker2"),
                           EC2("worker3"),
                           EC2("worker4"),
                           EC2("worker5")]

        with Cluster("Isolated subnet"): 
            with Cluster("DBCluster"):
                db_write = RDS("Databases Write")
                db_read = [RDS("Databases Read")]

    dns >> cdn >> lb
    lb >> workers
    workers >> db_write
    db_write >> db_read 
