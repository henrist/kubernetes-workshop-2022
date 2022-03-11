# Kubernetes Workshop

## Setup: Get authenticated to the cluster

Sign in with your email:

```bash
gcloud init
```

To verify you have been properly authenticated run:

```bash
gcloud projects list
```

You should see the `fluent-buckeye-343615` listed.

Get credentials for `kubectl`:

```bash
gcloud container clusters get-credentials workshop --region europe-west1
```

If successful you will have a new context in `kubectl` that is now active:

```bash
kubectl config get-contexts
```

## Setup: Test cluster connection

```bash
kubectl get nodes
```

This list will show the nodes that are part of the cluster.

## Setup: Docker registry

To be able to push and pull images we need to be authenticated to Aritfact Registry:

```bash
gcloud auth configure-docker europe-west1-docker.pkg.dev
```

This will make Docker auth for this registry go via `gcloud`.

## Create your own namespace and default to it

Create namespace:

```bash
# replace NAME with e.g. henrik etc
kubectl create namespace NAME
```

To make `kubectl` command using your namespace by default (instead of the namespace called `default`) you can set it as a default:

```bash
kubectl config set-contexts --current --namespace=NAME
```

You could also add `--namespace NAME` (or `-n NAME` for short) to every `kubectl` command.

## Run your first pod

This is an imperative way of using the cluster and usually not recommended. But it's an easy way to get started for now.

```bash
kubectl run demo --image=nginx
```

See the status of it by running:

```bash
kubectl get pods
```

You can connect its port by forwarding the port via a tunnel:

```bash
kubectl port-forward demo 8081:80
```

Open http://localhost:8081

Play a bit with some commands to interact with the pod - see https://kubernetes.io/docs/reference/kubectl/cheatsheet/:

- To try get the logs from the running service (hint: `kubectl logs`)
- Can you follow the logs in realtime and see your requests?
- Describe the pod to show details about it
- Can you exec into the pod to see what files are there?

Delete the pod when you've finished:

```bash
kubectl delete pod demo
```

## A small queue worker-based system

A simple queue based system allows us to play a bit with deployments
and scaling. For this workshop there is a producer service that will
write to the queue and a consumer service that will poll the queue and
simulate a compute intensive task.

We'll be using Kafka for this and a Kafka cluster is already running within the Kubernetes cluster.
This is achieved by using the Strimzi Kubernetes Operator that manages this for us.
The setup can be seen in `kafka-operator` and `kafka-cluster` directories.

### Create your own Kafka Topic

We can create a Kafka Topic by using a Kubernetes resource thanks to
the Kafka Operator that's installed.

Note that this resource will be created in the `kafka` namespace.

Edit manifest for `kafka-topic` and change `name`.

Deploy it:

```bash
kubectl apply -f manifest.yaml
```

It should be created immediately and you can check with:

```bash
kubectl get kafkatopic -n kafka
```

There is an instance of `kafka-ui` running in the cluster that
gives simple insight to the Kafka cluster. You can connect to it by:

```bash
kubectl port-forward service/kafka-ui -n kafka-ui 8080:8080
```

http://localhost:8080

### Create a reusable config map

See `config-map` directory.

A config map is useful to extract some configuration. The example
we use in the workshop isn't that interesting, but consider e.g.
credentials and connection details for a database.

For this workshop we use a config map to make the manifest
decoupled from the kafka topic setup.

Modify `manifest.yaml` and change `KAFKA_TOPIC` and `KAFKA_CONSUMER_GROUP_ID`.
The value of `KAFKA_CONSUMER_GROUP_ID` should be set to e.g. `NAME-worker`
where `NAME` is the same as your topic name.

Deploy the config map:

```bash
kubectl apply manifest.yaml
```

### Deploy the producer

See `queue-producer` directory.

The producer uses a Deployment resource.
https://kubernetes.io/docs/concepts/workloads/controllers/deployment/

Have a look through `manifest.yaml` to see how it looks.

Deploy the application:

```bash
kubectl apply -f manifest.yaml
```

Verify that the application is running. To get an idea of the various
resources that has been created you can run `kubectl get all`.

The application is a very simplified Go program. If you'd like you can
modify the source of it and use the `build-and-push.sh <tag>` script to upload
a separate Docker image by using your own personal tag. You can then run that
build by changing `latest` to your tag in the manifest.

### Deploy the consumer / worker

See `queue-worker` directory.

Follow the same steps as for the producer above.

When the worker has been deployed you should be able to identify the
consumer group in Kafka UI and see the backlog.

### Scale up the producer

Modify the manifest of the producer and scale it up to 2-4 instances.

You should see that the backlog is increasing in Kafka UI.

Since this consumer is CPU bound you should be able to verify that
it is using all available resources. Locate the workload in the
Kubernetes Engine dashboard in Google Cloud Platform to learn a bit
what metrics is availble and other details that can be seen.
You will also be able to retrive the logs there.

### Modify the producer

Change the message sent on the queue (main.go, line 13) and
deploy the updated application. You will have to build and push
the Docker image by running `./bulid-and-push.sh TAG` (change `TAG`
to something of your choice) and modify `latest` tag in the manifest
file to what you specified when building.

The workshop examples uses `latest`. Normally it is recommended to
avoid using `latest` as any scaling event might be pulling a different
image.

Note that Kubernetes will by default pull `latest` tag, but other tags
are not checked. The [`imagePullPolicy`](https://kubernetes.io/docs/concepts/containers/images/)
option can be used to change this behaviour. Normally, when set up with
proper CI/CD, you will use an immutable tag.

### Scale down the producer to 0 directly

To illustrate other use of the CLI you can directly modify the
scaling of a deployment:

```bash
kubectl scale --replicas=0 deployment/producer
```

Be aware that such commands are usually not adviced, since it might
cause a different state than in the source manifest/source control.

Scale up the producer again to its previous value.

### Set up an horizontal auto scaler for the consumer

Configure an auto scaler that will allow you to drain the queue
and keep ahead with the rate of messages.

Be aware that we are running Kubernetes 1.21 in the workshop cluster.

https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/

https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/

## Exercise if time

Create a Deployment for the `hello-world` directory.
Set up an external load balancer and required resources to make it
publicly accessible.

https://cloud.google.com/kubernetes-engine/docs/how-to/load-balance-ingress

## End: Destroy your resources

Since most of your resources are created in your seperate namespace you
can delete the namespace directly:

```bash
kubectl delete namespace NAME
```

Also delete the Kafka topic:

```bash
kubectl delete -n kafkatopic TOPIC-NAME
```

## Useful resources

Google Cloud Platform dashboard: https://console.cloud.google.com/ (sign in with Miles email)

Running workloads: https://console.cloud.google.com/kubernetes/workload/overview?project=fluent-buckeye-343615

kubectl cheatsheet:

- https://kubernetes.io/docs/reference/kubectl/cheatsheet/
- https://dockerlabs.collabnix.com/kubernetes/cheatsheets/kubectl.html
