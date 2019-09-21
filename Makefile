
# Image URL to use all building/pushing image targets
IMG ?= teleport-plus:v1
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

CLUSTER_NAME=teleport-plus
KUBECTL=env KUBECONFIG="$(shell kind get kubeconfig-path --name=${CLUSTER_NAME})" kubectl

all: manager

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | ${KUBECTL} apply -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | ${KUBECTL} apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths="./..."

# Build the docker image
docker-build: manager
	docker build . -t ${IMG}

# Push the docker image
docker-push: docker-build
	kind load docker-image ${IMG} --name ${CLUSTER_NAME}
	#docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

start:
	kind create cluster --name ${CLUSTER_NAME} --config ./e2e/cluster.yaml --wait=300s

stop:
	kind delete cluster --name ${CLUSTER_NAME}

e2e: start docker-push setup-cert-manager setup-teleport deploy 

setup-cert-manager:
	${KUBECTL} create namespace cert-manager
	${KUBECTL} label namespace cert-manager certmanager.k8s.io/disable-validation=true
	${KUBECTL} apply -f https://github.com/jetstack/cert-manager/releases/download/v0.10.0/cert-manager.yaml

setup-teleport:
	${KUBECTL} create namespace teleport
	${KUBECTL} -n cert-manager wait --for=condition=available --timeout=60s deployment/cert-manager-webhook
	# cert-manager-webhook is not yet available
	sleep 5
	${KUBECTL} apply -n teleport -f ./e2e/certificate.yaml
	${KUBECTL} apply -n teleport -f ./e2e/teleport.yaml

reload: docker-push
	${KUBECTL} delete pod -n teleport -l app=teleport


.PHONY: start stop e2e setup-cert-manager setup-teleport
