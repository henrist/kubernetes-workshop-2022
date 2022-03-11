# Kafka UI

```bash
kubectl create namespace kafka-ui
kubectl apply -f manifest.yaml

kubectl -n kafka-ui port-forward service/kafka-ui 8080:8080
```
