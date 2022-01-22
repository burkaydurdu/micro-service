# Micro Service

### Minikube start
```shell
minikube start --profile=custom
```

### Bind Minikube to docker
If you exec docker command on local, the command will run in the **minikube** environment.
```shell
// Fish Shell
eval (minikube -p custom docker-env)
// Bash Shell
eval $(minikube -p custom docker-env)

```
### Minikube SSH
```shell
minikube -p custom ssh
```
### List of docker image in minikube
```shell
docker images
```

### Skaffold Install
```shell
brew install skaffold
```

### Skaffold Init
```shell
skaffold init --generate-manifests
```

### Install K9s and Run
```shell
brew install k9s
k9s
```

### Run Skaffold on Dev
```shell
skaffold dev
```

### Install Hazelcast with helm

```shell
helm repo add hazelcast https://hazelcast-charts.s3.amazonaws.com/
helm repo update
helm install hz-hazelcast hazelcast/hazelcast --set cluster.memberCount=1 --set mancenter.enabled=false
```

### K9s Keymap
- ``d`` => describe
- ``Shift + f`` => Port forward

### Requests

``curl -v http://localhost:3001/orders``
response:
```
> GET /orders HTTP/1.1
> Host: localhost:3001
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Sat, 22 Jan 2022 15:49:16 GMT
< Content-Length: 62
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact
  I send create payment request and got back {"message": "paid"}* Closing connection 0
```

---
🎉 Thank you! [Hüseyin Babal](https://github.com/huseyinbabal)