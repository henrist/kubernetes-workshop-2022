resource "google_artifact_registry_repository" "workshop" {
  provider = google-beta

  location = var.region
  repository_id = "workshop"
  description = "workshop docker repository"
  format = "DOCKER"
}
