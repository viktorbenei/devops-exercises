# empty username & password disables "basic auth" (https://www.terraform.io/docs/providers/google/r/container_cluster.html)
variable "gke_master_username" {
    description = "GKE Master Username. Note: if both username and password set to empty value that disables basic auth ( https://www.terraform.io/docs/providers/google/r/container_cluster.html )."
    default = ""
}

variable "gke_master_password" {
    description = "GKE Master Password. Note: if both username and password set to empty value that disables basic auth ( https://www.terraform.io/docs/providers/google/r/container_cluster.html )."
    default = ""
}

variable "google_service_account_json" {
    description = "Google Cloud (GCP) Service Account JSON."
}

# The ID of the project on Google Cloud in which the Kubernetes Cluster should be created.
variable "google_project_id" {
    description = "Google Cloud (GCP) **Project ID** (the one you can find in the URL, or on the Dashboard of the GCP project - https://console.cloud.google.com/home/dashboard )."
}

# Name of the Kubernetes cluster to be created
variable "k8s_cluster_name" {
    description = "Name of the Kubernetes Cluster you want to create."
}