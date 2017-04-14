Following http://www.mikeda.me/hacking-controller-openshiftkubernetes/

- First

```bash
brew install golang
export GOPATH=~/go
export OS_OUTPUT_GOPATH=1
brew install mercurial
```

- Then

```bash
go get github.com/tools/godep

cd $GOPATH/src/github.com/openshift
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

