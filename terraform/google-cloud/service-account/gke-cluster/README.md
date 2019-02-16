# Terraform - Create Service Account suitable for GKE Kubernetes Cluster create

A Terraform config to create a Google Cloud Platform (GCP) Service Account with all the required roles
so that it can then be used to create a GKE (Google Kubernetes Engine) Kubernetes Cluster (in another terraform config for example).

Demonstrates how you can create & manage Google Cloud Service Accounts via `terraform`.

## How to use

You have to specify two variables:

- `gcp_project_id`, the Project ID in GCP
- and `google_service_account_json_for_iam_create`, the Service Account **already registered** on GCP, which has the roles to create & manage Service Accounts:
    - `Service Account Admin`
    - `Project IAM Admin`

You can provide these by any means supported by `terraform`.

If you just run `terraform apply` in this directory then `terraform` will ask for these variables (if it is an interactive shell, which if you run this manually, you almost certainly run it in a Terminal/Command Line which is an interactive shell).


Alternatively, if you want to avoid the interactive input, probably the easiest is to just set them as env vars before running terraform:

```
export TF_VAR_google_service_account_json_for_iam_create='...'
export TF_VAR_gcp_project_id='...'
```

Then (in this directory, where the `main.tf` and this `README.md` are located) run:

```
terraform apply
```

Once `terraform apply` is finished the Service Account will be registered and ready to use.
You can locate it at https://console.cloud.google.com/iam-admin/serviceaccounts , and **create a key** for it,
which you can then use with other tools/configs to create and manage a GKE Kubernetes Cluster.
