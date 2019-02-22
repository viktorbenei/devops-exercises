variable "gcp_project_id" {
    description = "Google Cloud (GCP) **Project ID** (the one you can find in the URL, or on the Dashboard of the GCP project - https://console.cloud.google.com/home/dashboard )."
}
variable "google_service_account_json_for_iam_create" {
    description = "The Service Account **already registered** on GCP, which has the roles to create & manage Service Accounts: 'Service Account Admin' and 'Project IAM Admin'"
}