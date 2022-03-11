resource "google_compute_network" "workshop" {
  name = "workshop-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "workshop" {
  name          = "workshop-subnetwork"
  ip_cidr_range = "10.111.0.0/16"
  region        = var.region
  network       = google_compute_network.workshop.id

  # secondary_ip_range {
  #   range_name    = "services-range"
  #   ip_cidr_range = "192.168.1.0/24"
  # }

  # secondary_ip_range {
  #   range_name    = "pod-ranges"
  #   ip_cidr_range = "192.168.64.0/22"
  # }
}
