# The Art of Infrastructure Management with Kubernetes API Machinery
![Screen Shot 2023-09-19 at 10 54 34 AM (1)](https://github.com/kellybw12/art-of-infrastructure-management/assets/72891670/d91b2e16-02f5-4bd5-b353-508943723c64)


The goal of this Level Up Lab demo is to illustrate how Kubernetes can be leveraged in real-world situations.

## Description

Infrastructure management is both an art and a science, requiring both technical expertise and creative problem-solving skills. Kubernetes provides powerful tools to implement a declarative style API to manage systems. However, many users of Kubernetes may not be fully aware of its capabilities or use it effectively. In this session, speakers will run through an infrastructure management use-case leveraging the Kubernetes API Machinery to help you gain insights into using Kubernetes control plane to manage resources.

Kubernetes API Machinery is attractive for three reasons: it is a powerful toolkit (e.g. provides a customizable declarative API and resolution engine), it’s a rapidly evolving technology and it’s a valuable skill to have for software engineers working with Kubernetes. Speakers will divide the workshop into three parts, each building upon the previous one to gradually introduce you to building an infrastructure automation tool using the Kubernetes control plane.

## Getting Started
### Setup Environment
You’ll need a Kubernetes cluster to run against. And localstack setup to simulate the AWS services on your local system.
There's setup scripts available for ubuntu and macos. Windows user please setup ubuntu virtual machine first then follow ubuntu setup script.

#### MacOS

```bash
# If you don't have docker installed, please run below command to install docker
# If you already have it installed, please open docker desktop
bash scripts/macos/setup_docker.sh 

# Please run below command to setup kind cluster and kubebuilder
bash scripts/macos/setup_kind_env.sh

# Please run below command to setup localstack and awscli
bash scripts/macos/setup_localstack.sh
```

#### Ubuntu

```bash
# If you don't have docker installed, please run below command to install docker
bash scripts/ubuntu/setup_docker.sh 

# Please run below command to setup kind cluster and kubebuilder
bash scripts/ubuntu/setup_kind_env.sh

# Please run below command to setup localstack and awscli
bash scripts/ubuntu/setup_localstack.sh 
```

**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

## Useful commands

1. List the current s3 buckets (AWS API Server)

```sh
aws s3 ls --endpoint=http://localhost:4566
```

2. List the current s3 bucket crds

```sh
kubectl get s3buckets.bucket.my.domain
```

## Demo Part 2: Simple Example

1. Run the controllers

```sh
make run
```

2. In a separate terminal, update the spec.count of the BucketGroup to 4

```sh
kubectl patch s3bucketgroups.bucketgroup.my.domain s3bucketgroup-sample --patch '{"spec": {"desiredBucketCount":4}}' --type=merge
```

3. In the same terminal as Step 2, run the following commands to delete my-s3-bucket-1 and my-s3-bucket-2

```sh
aws s3 rb s3://my-s3-bucket-1 --endpoint=http://localhost:4566
aws s3 rb s3://my-s3-bucket-2 --endpoint=http://localhost:4566
```

## Demo Part 3: Complex Example

1. **Comment** lines 74-77 and **uncomment** lines 80-83 in /internal/controllers/s3bucketgroup_controller.go

```sh
    74 // result, err := DoPart2(r, ctx, req)
    75 // if err != nil {
    76 // 	log.Log.Error(err, "error occurred when running part 2")
    77 // }
    78
    79
    80 result, err := DoPart3(r, ctx, req)
    81 if err != nil {
    82	log.Log.Error(err, "error occurred when running part 3")
    83 }
```

2. Run the controllers

```sh
make run
```

3. In a separate terminal, delete my-s3-bucket-1

```sh
aws s3 rb s3://my-s3-bucket-1 --endpoint=http://localhost:4566
```

### Running on the cluster

1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/art-of-infrastructure-management:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/art-of-infrastructure-management:tag
```

### Uninstall CRDs

To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller

UnDeploy the controller from the cluster:

```sh
make undeploy
```

### How it works

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out

1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## Start Your Project With Kubebuilder
[Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) provides powerful libraries and tools to simplify building and publishing Kubernetes APIs from scratch.
Strongly recommend follow the [quick start](https://book.kubebuilder.io/quick-start) to build your project and explore more options.


## More Examples of KRM-based APIs
* [Crossplane](https://docs.crossplane.io/v1.11/getting-started/introduction/) - KRM for control of cloud resources
* [Google Cloud Config Connector](https://cloud.google.com/config-connector/docs/overview)
* [Operator Hub](https://operatorhub.io/) - Directory of k8s operators.
* [Cluster-API](https://cluster-api.sigs.k8s.io/) - an Operator to manage other k8s cluster
* [Cruster-API](https://github.com/rudoi/cruster-api) - a CRD that orders pizza!

## More Resources of Kubernetes
* Kubernetes: Up and Running by Brendan Burns, Joe Beda, Kelsey Hightower
* Kubernetes in Action by Marko Lukša
* [The Mechanics of Kubernetes](https://dominik-tornow.medium.com/the-mechanics-of-kubernetes-ac8112eaa302) by Andrew Chen and Dominik Tornow
* For learning more about extending Kubernetes:
* Programming Kubernetes by Michael Hausenblas and Stefan Schimanski

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
