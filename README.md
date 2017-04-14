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

- Add controller.Run call in [cmd.go](controller/cmd/controller/cmd.go) 
- Run `make`

  **Errors:**

```
go install github.com/dlbewley/openshift-controller-hacking/controller/cmd/controller
# github.com/dlbewley/openshift-controller-hacking/controller/cmd/controller
controller/cmd/controller/cmd.go:4: imported and not used: "fmt"
controller/cmd/controller/cmd.go:30: cannot use openshiftClient (type client.Interface) as type *client.Client in argument to controller.NewController: need type assertion
controller/cmd/controller/cmd.go:31: not enough arguments in call to c.Run
      have ()
        want (<-chan struct {})
make: *** [all] Error 2
```

Presumably because I don't have openshift running on the mac?
I do have a working client config.

```
$ oc version
oc v3.4.1.5
kubernetes v1.4.0+776c994
features: Basic-Auth

Server https://openshift.example.com:8443
openshift v3.4.1.10
kubernetes v1.4.0+776c994
```

- Try building and starting origin following https://github.com/openshift/origin/blob/master/CONTRIBUTING.adoc#openshift-development

```bash
$ cd $GOPATH/src/github.com/openshift/origin
$ make clean build
    rm -rf _output Godeps/_workspace/pkg
    hack/build-go.sh
    ++ Building go targets for darwin/amd64: cmd/openshift cmd/oc
    ++ Placing binaries
    hack/build-go.sh took 65 seconds
```

- Add `172.30.0.0/16` as insecure registry

- Start openshift. Doesn't seem to work on Mac

```bash
cd $GOPATH/src/github.com/openshift/origin/_output/local/bin/darwin/amd64
sudo ./openshift start
...
Created node config for fakenews in openshift.local.config/node-fakenews
I0413 18:56:37.233921   31490 plugins.go:71] No cloud provider specified.
I0413 18:56:37.233945   31490 start_node.go:288] Starting node fakenews (v1.2.2-4-g6f611f3)
    I0413 18:56:37.234937   31490 start_node.go:297] Connecting to API server https://192.168.1.41:8443
    I0413 18:56:37.238519   31490 node.go:131] Connecting to Docker at unix:///var/run/docker.sock
    W0413 18:56:37.240558   31490 iptables.go:144] Error checking iptables version, assuming version at least 1.4.11: executable file not found in $PATH
    I0413 18:56:37.240637   31490 node.go:349] Using iptables Proxier.
    F0413 18:56:37.240703   31490 node.go:359] error: Could not initialize Kubernetes Proxy. You must run this process as root to use the service proxy: can't set sysctl net/ipv4/conf/all/route_localnet: open /proc/sys/net/ipv4/conf/all/route_localnet: no such file or directory
    W0413 18:56:37.281518   31490 server.go:349] setting OOM scores is unsupported in this build
```

