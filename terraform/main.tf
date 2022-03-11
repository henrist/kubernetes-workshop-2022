terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "4.13.0"
    }
  }
}

provider "google" {
  credentials = file("../keys/fluent-buckeye-343615-6016d84b8b9f.json")

  project = "fluent-buckeye-343615"
  region  = var.region
  zone    = "europe-west1-b"
}


provider "google-beta" {
  credentials = file("../keys/fluent-buckeye-343615-6016d84b8b9f.json")

  project = "fluent-buckeye-343615"
  region  = var.region
  zone    = "europe-west1-b"
}
