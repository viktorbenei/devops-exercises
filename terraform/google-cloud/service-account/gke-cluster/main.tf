variable "gcp_project_id" {}
variable "google_service_account_json_for_iam_create" {}

provider "google" {
  # provider plugin version
  version = "~> 1.20"

  credentials = "${var.google_service_account_json_for_iam_create}"
  project     = "${var.gcp_project_id}"
}

# create service account
resource "google_service_account" "tf-admin" {
  account_id   = "terraform-admin"
  display_name = "Terraform Admin"
  project      = "${var.gcp_project_id}"
}

# attach roles
resource "google_project_iam_binding" "terraform-admin-role-ServiceAccountUser" {
  role = "roles/iam.serviceAccountUser"

  members = [
    "serviceAccount:${google_service_account.tf-admin.email}",
  ]
}

resource "google_project_iam_binding" "terraform-admin-role-KubernetesEngineAdmin" {
  role = "roles/container.admin"

  members = [
    "serviceAccount:${google_service_account.tf-admin.email}",
  ]
}

resource "google_project_iam_binding" "terraform-admin-role-KubernetesClusterAdmin" {
  role = "roles/container.clusterAdmin"

  members = [
    "serviceAccount:${google_service_account.tf-admin.email}",
  ]
}

resource "google_project_iam_binding" "terraform-admin-role-ComputeAdmin" {
  role = "roles/compute.admin"

  members = [
    "serviceAccount:${google_service_account.tf-admin.email}",
  ]
}
