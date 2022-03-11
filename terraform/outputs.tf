output "kubernetes_cluster_name" {
  value       = google_container_cluster.workshop.name
  description = "GKE Cluster Name"
}

output "kubernetes_cluster_host" {
  value       = google_container_cluster.workshop.endpoint
  description = "GKE Cluster Host"
}

output "region" {
  value       = var.region
}
