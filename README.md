Following http://www.mikeda.me/hacking-controller-openshiftkubernetes/

- First

_Mac_

```bash
brew install golang
export GOPATH=~/go
export OS_OUTPUT_GOPATH=1
brew install mercurial
```

_Linux_

```bash
sudo yum install golang mercurial
export GOPATH=~/go
export OS_OUTPUT_GOPATH=1
```

- Then

```bash
go get github.com/tools/godep

git clone git://github.com/openshift/origin $GOPATH/src/github.com/openshift/origin
cd $GOPATH/src/github.com/openshift/origin
git checkout release-1.2
 
git clone git://github.com/kubernetes/kubernetes $GOPATH/src/k8s.io/kubernetes
cd $GOPATH/src/k8s.io/kubernetes
git remote add openshift git://github.com/openshift/kubernetes
git fetch openshift

COMMIT_ID=$(cat Godeps/Godeps.json | jq -r '.Deps[] | select(.ImportPath=="k8s.io/kubernetes/pkg/api") .Comment')
echo $COMMIT_ID
v1.2.0-36-g4a3f9c5
git checkout v1.2.0-36-g4a3f9c5
 
git clone https://github.com/go-inf/inf.git $GOPATH/src/speter.net/go/exp/math/dec/inf
 
cd $GOPATH/src/github.com/openshift/origin
godep restore
```

- Add `172.30.0.0/16` as insecure registry and restart docker
- Build and start openshift in Docker

```bash
$ cd $GOPATH/src/github.com/openshift/origin
$ make clean build

cd $GOPATH/src/github.com/openshift/origin/_output/local/bin/linux/amd64
sudo ./openshift start
mv ~/.kube{,.bak}
export oc=/home/dlbewley/go/src/github.com/openshift/origin/_output/local/bin/linux/amd64/oc
$oc login
 developer / developer
```

- Make new project on github called `openshift-controller-hacking`

```bash
git clone git@github.com:$USER/openshift-controller-hacking.git $GOPATH/src/github.com/$USER/openshift-controller-hacking

cd $GOPATH/src/github.com/$USER/openshift-controller-hacking
mkdir -p controller/{cmd,pkg}/controller
```

- Create [$GOPATH/src/github.com/$USER/openshift-controller-hacking/controller/cmd/controller/cmd.go](controller/cmd/controller/cmd.go)

- Save deps

```bash
cd $GOPATH/src/github.com/$USER/openshift-controller-hacking
godep save ./...
git add .
git commit -am firstsies
git push
```

- Create [Makefile](Makefile)

- Create [$GOPATH/src/github.com/$USER/openshift-controller-hacking/controller/pkg/controller/controller.go](controller/pkg/controller/controller.go) 

    **Refs:**

    - https://godoc.org/github.com/openshift/origin/pkg/client
    - https://godoc.org/github.com/openshift/origin/pkg/client#ProjectInterface
    - https://godoc.org/github.com/openshift/origin/pkg/project/api#Project
    - https://godoc.org/k8s.io/kubernetes/pkg/api#ObjectMeta

- Add controller.Run call in [cmd.go](controller/cmd/controller/cmd.go) 
- Run `make`

  **Errors:**



```
go install github.com/dlbewley/openshift-controller-hacking/controller/cmd/controller
# github.com/dlbewley/openshift-controller-hacking/controller/cmd/controller
controller/cmd/controller/cmd.go:4: imported and not used: "fmt"
controller/cmd/controller/cmd.go:30: cannot use openshiftClient (type client.Interface) as type *client.Client in argument to controller.NewController: need type assertion
controller/cmd/controller/cmd.go:31: not enough arguments in call to c.Run
make: *** [all] Error 2
```

I do have a working client config and server running..

```
$oc version
oc v1.2.2-4-g6f611f3
kubernetes v1.2.0-36-g4a3f9c5
```

