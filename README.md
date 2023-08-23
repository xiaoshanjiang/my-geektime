## Using minikube for k8s if you are using a virtual environment:
1. In a seperate terminal run `minikube start`
2. Under the `webook` directory, run `make docker` \
This will rebuild the `webook-app` binary and subsequently use it to build the docker image `xjiang91/webook`, version `v0.0.1`. Modify the image prefix in the `Dockerfile` if needed.
3. run `docker tag $(docker images xjiang91/webook -q) xjiang91/webook:v0.0.1 && docker push xjiang91/webook:v0.0.1` \
This will tage the image as `xjiang91/webook:v0.0.1` and push to docker hub.
4. Run the following commands to start the k8s service:\
`chmod 700 k8s.sh && ./k8s.sh`

There are two ways to access the services:
1. To access a service (e.g.webook) via nodePort, run in terminal `minikube service webook --url`, and grab the address and go to the brower to visit it, e.g. `http://192.168.76.2:32225/hello`. \
Similary, you can access the mysql service by `minikube service webook-mysql --url`, and the redis service by `minikube service webook-redis --url`. You can then use the IP address and port to login.
2. If the service type is `LoadBalancer`, you can also open a tunnel to the minikube cluster by running in a seperate terminal: `minikube tunnel` and keep that terminal alive, then use `kubectl get services` to check the EXTERNAL-IP and port.\

For more details about `minikube`, please refer to the official documentation: https://minikube.sigs.k8s.io/docs/handbook/
