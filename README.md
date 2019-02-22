# devops-exercises

DevOps Exercises.

## Important notes

These are examples/exercises **not production ready setups** in most cases! Good for getting started when you need some help/pointers while you're learning any of the tools but don't just copy-paste e.g. a Terraform config and expect it to be production ready!

That said I try to update these when I can, but **I don't aim for keeping these up-to-date all the time**. These are examples used in own exercises/explorations, and short demonstrations how certain things can be done.

## The list

- `kubernetes-setup/gke/terraform`: Create a GKE Kubernetes Cluster with Terraform.
- `kubernetes-setup/gke/gcloud`: Create a GKE Kubernetes Cluster with `gcloud` (Google Cloud CLI).
- `simple-code/echo-server`: Minimal server example, used for testing infra/kubernetes features.
    - The server port can be changed via `PORT` env var.
    - Kubernetes Deploy included.
    - Endpoints:
        - `/` : Simple "Welcome" message, with version included.
        - `/hi?name=Someone` : Returns a message that includes the specified `name`.
        - `/auth-via-kubernetes-secret` : Auth test.
            - Send the auth token in header `Authorization` with value `token TheToken`
            - `TheToken` has to match with the Secret set in the `echo-server-auth-secret-token` k8s secret's `token` data, otherwise you'll get an "Unauthorized".
- `terraform/google-cloud/service-account/gke-cluster`: A Terraform config to create a Google Cloud Platform (GCP) Service Account with all the required roles so that it can then be used to create a GKE (Google Kubernetes Engine) Kubernetes Cluster (in another terraform config for example).
    - Demonstrates how you can create & manage Google Cloud Service Accounts via `terraform`.
