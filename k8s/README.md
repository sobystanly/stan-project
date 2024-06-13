# Local K8s Cluster Setup(KIND)

The following instructions will setup a local node k8s cluster for use.

- We will be using [KIND - K8s in docker](https://kind.sigs.k8s.io/) and also creating a local docker registry that the
  KIND cluster has access to and you can push your application images to.

1. Install KIND using these [instructions](https://kind.sigs.k8s.io/docs/user/quick-start/#installation). 
   On MacOS, the easiest method is using homebrew `brew install kind`.  
   On Linux (or Linux VM), you can install using the [release binaries](https://kind.sigs.k8s.io/docs/user/quick-start/#installing-from-release-binaries).
   Note that if you are on an M1 Mac (ARM arch) but want to perform this exercise in a Linux VM you will need the `arm64` linux binary.

2. Install `kubectl` if you do not have it already installed - [instructions here](https://kubernetes.io/docs/tasks/tools/).  On MacOS, the easiest method is to use homebrew `brew install kubectl`.

3. Setup KIND local cluster and registry

```
   # If your docker install is setup like this: https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user run this:
   make create-cluster

   # If for example you execute docker as `sudo docker ...`, run this:
   make create-cluster-sudo

```
4. Test cluster
```
   make test-cluster
```
Should return something like the below:
```
Kubernetes control plane is running at https://127.0.0.1:36063
CoreDNS is running at https://127.0.0.1:36063/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
NAME                 STATUS   ROLES           AGE   VERSION
kind-control-plane   Ready    control-plane   28m   v1.24.0
```
5.When KIND creates the cluster it should append to your current `~/.kube/config` (if it exists) or create a new one.  If you run as `make create-cluster-sudo` (e.g. your docker setup is not setup as a non root user), it will copy the config to `~/.kube` as `~/.kube/kind-config` and then append to your `$KUBECONFIG` in your current session.  If at any time you run into a problem with the kubeconfig, you can generate it yourself via `sudo kind get kubeconfig > ~/.kube/kind-kubeconfig` and then manually append it to your `$KUBECONFIG` like this `export KUBECONFIG=$KUBECONFIG:~/.kube/kind-kubeconfig`.

## Important info about local images / local registry using KIND

There are some specific workflows to be aware of in order to deploy your containerized application on the local cluster. In order to use the image in the cluster you will need to build and push your image to the local registry created.

For example: build/tag it `docker build -t localhost:5001/image:tag .` and then push it to the registry `docker push localhost:5001/image:foo`.  You can then reference it in the cluster as `localhost:5001/image:tag` since the local cluster will have access to the registry.  For more information see the [documentation](https://kind.sigs.k8s.io/docs/user/local-registry/)

## Cleanup

If you installed the local KIND cluster you can cleanup / delete it with `kind cluster delete kind` or use the included [Makefile](/Makefile)

## Makefile Commands
- `make create-cluster`: Creates the Cluster
- `make create-cluster-sudo`: Creates the Cluster as non-root user
- `make test-cluster:`: Get info about currently running KIND cluster.
- `make remove-cluster`: Removes the cluster.
- `make remove-cluster-sudo`: Removes the cluster as a non-root user.

## Kubernetes manifests
- The [k8s folder](/k8s) consist of all necessary configuration files to deploy the `Risk` service to k8s cluster.

## How to build docker image of `Risk` service to deploy to the k8s cluster

- Navigate to the `stan-project` directory where the `Dockerfile` is located.
- Run the following command to build the docker image

```
    docker build -t localhost:5001/risk:1.0.0 
```
- Push the image to local Docker registry
```
    docker push localhost:5001/risk:1.0.0
```

## Deploying the `Risk` service to the k8s Cluster

- Ensure that the KIND cluster is created 
- Navigate to the [k8s folder](/k8s) in your terminal, we will be using `kubectl` commands from terminal

**Deploy the components in the following order**
  **Create the PostgresSQL Secret**
  - Run the following kubectl command to create the secret for the PostgreSQL database.
  - For simplicity, we use the opaque secret type and base64 encode both the PostgreSQL username and password as postgres.
  ```
  kubectl apply -f postgres-secret.yaml
  ```
  **Create the PostgresSQL ConfigMap**
  - Run the following kubectl command to create the PostgreSQL ConfigMap.
  - This acts as centralized configuration management.
  - In the context of this service, we use it to override the PostgreSQL address with the Kubernetes postgres internal 
    service and to add the Docker init script to the PostgreSQL deployment YAML.
  ```
  kubectl apply -f postgres-configmap.yaml
  ```
  **Deploy PostgresSQL**
  - Run the following kubectl command to deploy PostgreSQL with one replica.
  ```
  kubectl apply -f postgres.yaml
  ```
  **Deploy the `Risk` Service**
  - The tag of the Docker image built earlier is specified in [deployment.yaml](deployment.yaml).
  - Run the following kubectl command to deploy the Risk service.
  ```
  kubectl apply -f deployment.yaml
  ```
  **Create the `Risk` internal service**
  - Run this kubectl command to create `risk` k8s internal service.
  ```
  kubectl apply -f service.yaml
  ```
  **Create the ingress for the `Risk` service**
  - Run this kubectl command to create the ingress for `Risk` service
  - The ingress expose the Risk API with host as `risk-app.com`
  - After deploying this update your /etc/hosts to map 127.0.0.1 to `risk-app.com`
  ```
  kubectl apply -f ingress.yaml
  ```
  - All necessary Kubernetes components for the Risk service are now installed.
  - Run `kubectl get all` to fetch all components' info, and it should return something like this

  - Now the requests to the `Risk` can be made with host `risk-app.com`
  ```http request
      gPOST http://risk-app.com/v1/risk
```
```http request
    GET http://risk-app.com/v1/risk/{id}
```
```http request
    GET http://risk-app.com/v1/risks?offset=0&limit=10&sortBy=title&sortOrder=desc
```
  ## Risk Service Request Flow
    
 ```mermaid
    graph TD
    subgraph User
        U[User]
    end

    subgraph KIND Cluster
        subgraph Ingress
            I[Ingress]
        end

        subgraph Services
            RISK_SVC[Risk Service]
        end

        subgraph Deployments
            subgraph Risk App
                RISK_APP1[Risk App Pod 1]
                RISK_APP2[Risk App Pod 2]
                RISK_APP3[Risk App Pod 3]
            end
        end

        subgraph StatefulSets
            subgraph PostgreSQL
                PG[PostgreSQL Pod]
            end
        end
    end

    U -->|Requests| I
    I -->|Routes to /| RISK_SVC
    RISK_SVC --> RISK_APP1
    RISK_SVC --> RISK_APP2
    RISK_SVC --> RISK_APP3
    RISK_APP1 -->|Database queries| PG
    RISK_APP2 -->|Database queries| PG
    RISK_APP3 -->|Database queries| PG
```