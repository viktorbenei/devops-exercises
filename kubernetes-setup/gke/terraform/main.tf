#
# Docs: https://www.terraform.io/docs/providers/google/r/container_cluster.html
#

# empty username & password disables "basic auth" (https://www.terraform.io/docs/providers/google/r/container_cluster.html)
variable "gke_master_username" {}

variable "gke_master_password" {}

variable "google_service_account_json" {}

# The ID of the project on Google Cloud in which the Kubernetes Cluster should be created.
variable "google_project_id" {}

# Name of the Kubernetes cluster to be created
variable "k8s_cluster_name" {}

provider "google" {
  # provider plugin version
  version = "~> 1.20"

  credentials = "${var.google_service_account_json}"
  project     = "${var.google_project_id}"
  region      = "us-central1"
  zone        = "us-central1-c"
}

resource "google_container_cluster" "primary" {
  name               = "${var.k8s_cluster_name}"
  zone               = "us-central1-a"
  min_master_version = "1.11.5"

  # If additional zones are configured, the number of nodes specified in initial_node_count is created in all specified zones.
  initial_node_count = 1

  additional_zones = [
    "us-central1-b",
    "us-central1-c",
  ]

  master_auth {
    # empty username & password disables "basic auth" (https://www.terraform.io/docs/providers/google/r/container_cluster.html)
    username = "${var.gke_master_username}"
    password = "${var.gke_master_password}"

    #
    client_certificate_config {
      issue_client_certificate = false
    }
  }

  node_config {
    machine_type = "n1-standard-1"

    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    labels {
      foo = "bar"
    }

    tags = ["foo", "bar"]
  }
}

# The following outputs allow authentication and connectivity to the GKE Cluster.
# output "client_certificate" {
#   value = "${google_container_cluster.primary.master_auth.0.issue_client_certificate}"
# }


# output "client_key" {
#   value = "${google_container_cluster.primary.master_auth.0.client_key}"
# }


# output "cluster_ca_certificate" {
#   value = "${google_container_cluster.primary.master_auth.0.cluster_ca_certificate}"
# }

