## Using minikube for k8s if you are using a virtual environment:
1. In a seperate terminal run `minikube start`
2. Under the `webook` directory, run `make docker` \
This will rebuild the `webook-app` binary and subsequently use it to build the docker image `xjiang91/webook`, version `v0.0.1`. Modify the image prefix in the `Dockerfile` if needed.
3. run `docker tag $(docker images xjiang91/webook -q) xjiang91/webook:v0.0.1 && docker push xjiang91/webook:v0.0.1` \
This will tage the image as `xjiang91/webook:v0.0.1` and push to docker hub.
4. Run the following commands to start the k8s service:\
`kubectl apply -f k8s-mysql-pv.yaml` \
`kubectl apply -f k8s-mysql-pvc.yaml` \
`kubectl apply -f k8s-mysql-service.yaml` \
`kubectl apply -f k8s-mysql-deployment.yaml` \
`kubectl apply -f k8s-redis-service.yaml` \
`kubectl apply -f k8s-redis-deployment.yaml` \
`kubectl apply -f k8s-webook-service.yaml` \
`kubectl apply -f k8s-webook-deployment.yaml` \
`kubectl apply -f k8s-ingress-nginx.yaml`
5. To access the web page, run in terminal `minikube service webook --url`, and grab the address and go to the brower to visit it, e.g. `http://192.168.76.2:32225/hello`
6. To access the mysql database, run in terminal `minikube service webook-mysql --url`, and use the IP address and port when login datbase.
7. Similary, `minikube service webook-redis --url` will give you the IP address and port to access redis.

For more details about `minikube`, please refer to the official documentation: https://minikube.sigs.k8s.io/docs/handbook/
