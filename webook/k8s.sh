#!/bin/bash
kubectl apply -f k8s-mysql-pv.yaml
kubectl apply -f k8s-mysql-pvc.yaml
kubectl apply -f k8s-mysql-deployment.yaml
kubectl apply -f k8s-redis-service.yaml
kubectl apply -f k8s-mysql-service.yaml
kubectl apply -f k8s-redis-deployment.yaml
kubectl apply -f k8s-webook-service.yaml
kubectl apply -f k8s-webook-deployment.yaml
kubectl apply -f k8s-ingress-nginx.yaml
