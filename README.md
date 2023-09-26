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
2. If the service type is `LoadBalancer`, you can also open a tunnel to the minikube cluster by running in a seperate terminal: `minikube tunnel` and keep that terminal alive, then use `kubectl get services` to check the EXTERNAL-IP and port.

For more details about `minikube`, please refer to the official documentation: https://minikube.sigs.k8s.io/docs/handbook/

## 方案
1. 采用circular queue储存最近N次SMS请求结果，计算成功率。如果成功率低于阈值，或者当前服务触发限流，则将短信转存到数据库，并切换当前服务商。
2. 启动单独的goroutine进行重试，循环读取发送失败的短信，从最旧时间戳开始依次重试。如果成功则更新数据库，将短信标记为成功，如果失败则更新时间戳。
3. 所有方法都采用读写锁（数据库交互）或者原子操作（结构体本身）以保证并发安全。
方案优点：用于计算成功率的样本大小可以按需调整，对于不同的业务可以根据成功率敏感度灵活调节。
方案缺点：对于每个服务商都需要维护一个请求结果的循环队列，会有一定的性能开销。
