# DataEngine Hello World stack

Terraform coverage of Data Engine **workload** resources (infra stays in
`../vms_setup` / `modules/dataengine_vms_setup`):

| Resource | Purpose |
|----------|---------|
| `vastdata_de_function` | Register container image |
| `vastdata_de_schedule_trigger` | Cron / quartz schedule |
| `vastdata_de_element_trigger` | S3 object events on a view/bucket |
| `vastdata_de_kubernetes_secret` | Upsert a K8s secret via DE API |
| `vastdata_de_pipeline` | Wire both triggers → function, deploy |

**Out of band (image build — not managed here):**

```bash
vastde functions init python-pip hello-world -t ~/test/functions/
vastde functions build -T latest hello-world -t ~/test/functions/hello-world
docker tag hello-world:latest fnn45/hello-world:latest
docker push fnn45/hello-world:latest
```

## Prerequisites

- `vms_setup` applied (engine, k8s cluster, registry, broker, topics)
- Function image pushed (`artifact_source`)
- Docker Hub: `skip_image_validation = true` (create-time probe often 401)

## Usage

```bash
cp terraform.tfvars.example terraform.tfvars
# fill from: cd ../vms_setup && terraform output
terraform init
terraform plan
terraform apply
```

Verify:

```bash
vastde functions list
vastde triggers list
vastde pipelines list
kubectl -n vast-pipelines get secret de-hello-app-secret
```

## Notes

- Updating the pipeline manifest triggers a **redeploy**. Lab Knative Revisions
  often fail with `RevisionMissing` / queue-proxy `ErrImagePull` when Zarf
  rewrites `queue-sidecar-image` to a tag/digest missing from the local
  registry. Fix once (pins KnativeServing + retags queue image + pull secrets):

  ```bash
  export KUBECONFIG=../rke2/out/kubeconfig
  ../scripts/fix-knative-queue-pull.sh
  ```

- Element trigger listens on `element_source_bucket` (default = `broker_name`).
- Kubernetes secret has no GET API; TF tracks identity from config.
