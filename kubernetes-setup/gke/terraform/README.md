## Apply:

1. Create a `secrets.tfvars` file here (make sure it's `.gitignore`d!!)
    - You can find an example below.
1. Make sure that all the `variable`s from `main.tf` are defined in the `secrets.tfvars` otherwise `terraform` will ask for those interactively.
1. Do a terraform plan to see what terraform would do (dry run), in this directory: `terraform plan -var-file=secrets.tfvars`
1. If it looks good run `apply` (just replace `plan` with `apply` in the command) to create the GKE Cluster: `terraform apply -var-file=secrets.tfvars`
1. To destroy it run the `destroy` command (again, just replace `apply` with `destroy` in the command): `terraform destroy -var-file=secrets.tfvars`

Cheat sheets: https://github.com/bitrise-io/cheat-sheets


## Example `secrets.tfvars`

```
# empty username & password disables "basic auth" (https://www.terraform.io/docs/providers/google/r/container_cluster.html)
gke_master_username = ""

gke_master_password = ""

google_service_account_json = <<EOF
...
EOF

google_project_id = "..."

k8s_cluster_name = "..."

```