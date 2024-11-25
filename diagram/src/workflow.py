from diagrams import Cluster, Diagram
from diagrams.custom import Custom

with Diagram("TheTipTop - Workflow", filename="workflow", show=False, outformat="png"):
    Custom("Git", "./assets/github.svg")