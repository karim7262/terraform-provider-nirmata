# Terraform Provider for Nirmata

## Building

````bash
go build
````

## Executing the samples

The samples are available in the [samples](samples) folder.

To run the samples, initialize the Terraform provider and then run the `plan` and `apply` commands. 

Here is an example of how to run the GKE cluster provisioning:

1. Clone this repository

```bash
git clone 
```

2. Build the Nirmata Terraform Provider

Golang needs to be installed to build the provider from source. See: https://golang.org/doc/install for Golang installation instructions and then run:

```bash
go build
```

3. Set your `NIRMATA_TOKEN` environment variable to contain your Nirmata API key. You can optionally set `NIRMATA_URL` to point to the Nirmata address (defaults to https://nirmata.io.)

4. Edit the sample Terraform config file `samples/cloud_provider/gke/gke.tf` and include your credentials, and desired region, machinetype, and disksize.

In Nirmata, a `ClusterType` is a reusable configuration that you can use to create several clusters. 

The example file first creates a ClusterType and then creates a single node `Cluster` using that type. Optionally, you can create the ClusterType via the Nirmata web console, or using 
[nctl](https://downloads.nirmata.io/nctl/downloads/), and then use the Terraform provider to create clusters of that type.

5. Initialize the Terraform provider with the correct directory

**NOTE: this assumes Terraform 0.12.x or lower. With Terraform 0.13.x local providers need to be initialized differently.** 

```bash
terraform init examples/gke
```

6. Run `plan` to build the execution plan:

```bash
terraform plan examples/gke
```

7. Run `apply` to execute the plan:

```bash
terraform apply examples/gke
```

8. Run `show` to see the created resources:

```bash
terraform show
```

9. Run `destroy` to delete the created resources:

````bash
terraform destroy samples/cloud_provider/gke
````

## Troubleshooting

Set the TF_LOG environment variable to `debug` or `trace`.


## Documentation

The provider documentation is available in the [docs](./docs) folder.


## Examples

Examples are available in the [examples](./examples) folder.
