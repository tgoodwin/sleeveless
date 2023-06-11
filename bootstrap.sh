# for ubuntu

# install docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# python is already installed

# install kind
[ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.19.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
# validate the binary
echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check
chmod +x ./kubecl
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
echo "testing kubectl version"
kubectl version --client

# setup kubeconfig
mkdir -p $HOME/.kube
export KUBECONFIG=$HOME/.kube/config

# install go
curl -OL https://golang.org/dl/go1.19.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.19.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo "verifying go installation"
go version

export PATH=$PATH:$(go env GOPATH)/bin
export GOPATH=$(go env GOPATH)

# install pip
sudo apt install python3-pip

# install helm
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
sudo chmod 700 get_helm.sh
./get_helm.sh

# install k9s
curl -LO https://github.com/derailed/k9s/releases/download/v0.27.4/k9s_Linux_amd64.tar.gz
sudo tar -C /usr/local/bin -xzf k9s_Linux_amd64.tar.gz


# install sieve
git clone https://github.com/sieve-project/sieve.git
cd sieve
pip3 install -r requirements.txt

# set up sieve environment

