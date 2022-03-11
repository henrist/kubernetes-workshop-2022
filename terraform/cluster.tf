resource "google_container_cluster" "workshop" {
  name     = "workshop"
  location = var.region

  network    = google_compute_network.workshop.id
  subnetwork = google_compute_subnetwork.workshop.id

  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "192.168.64.0/20"
    services_ipv4_cidr_block = "192.168.1.0/24"
  }

  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_service_account" "cluster_node" {
  account_id   = "cluster-node"
  display_name = "Cluster Node Service Account"
}

resource "google_project_iam_member" "cluster_node_artifactregistry" {
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${google_service_account.cluster_node.email}"
  project = "fluent-buckeye-343615"
}

resource "google_project_iam_member" "cluster_node_logging" {
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.cluster_node.email}"
  project = "fluent-buckeye-343615"
}

resource "google_project_iam_member" "cluster_node_monitoring" {
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.cluster_node.email}"
  project = "fluent-buckeye-343615"
}

resource "google_container_node_pool" "primary_nodes" {
  name     = "${google_container_cluster.workshop.name}-node-pool"
  location = var.region
  cluster  = google_container_cluster.workshop.name

  initial_node_count = 1

  autoscaling {
    min_node_count = 1
    max_node_count = 5
  }

  node_config {
    # oauth_scopes = [
    #   "https://www.googleapis.com/auth/logging.write",
    #   "https://www.googleapis.com/auth/monitoring",
    # ]

    service_account = google_service_account.cluster_node.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    labels = {
      env = "workshop"
    }

    preemptible  = true
    machine_type = "n1-standard-2"
    tags         = ["gke-node", "workshop-gke"]
    metadata = {
      disable-legacy-endpoints = "true"
    }
  }
}
