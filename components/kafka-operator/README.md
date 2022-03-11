# Strimzi Kafka Operator

```bash
kubectl create namespace kafka

curl -o manifest.yaml https://strimzi.io/install/latest?namespace=kafka

kubectl apply -f manifest.yaml -n kafka
```
