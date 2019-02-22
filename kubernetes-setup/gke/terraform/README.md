## Apply:

1. Do a terraform plan to see what terraform would do (dry run), in this directory: `terraform plan`
1. If it looks good run `apply` (just replace `plan` with `apply` in the command) to create the GKE Cluster: `terraform apply`
1. To destroy it run the `destroy` command (again, just replace `apply` with `destroy` in the command): `terraform destroy`

Cheat sheets: https://github.com/bitrise-io/cheat-sheets/blob/master/terraform.md

Note: `gke_master_username` and `gke_master_password` have a defined empty `default` value, so `terraform` won't ask for these. You can still provide these if you want to (e.g. by setting `TF_VAR_gke_master_username` and `TF_VAR_gke_master_password`), but you should leave these empty for production kubernetes clusters as that means (as you can see in the `variables.tf` file) that basic auth based authentication will be completely disabled (best practice to disable it).
