
Create cluster:

```
kind create cluster --config=kind-config.yaml
```

Install MySQL:

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install --set auth.rootPassword=password --set auth.rootPassword=password --set auth.database=users mysql bitnami/mysql
```

Set NodePort to expose MySQL:

```
kubectl get service mysql -o yaml | sed 's/type: ClusterIP/type: NodePort/g' | kubectl replace -f -

kubectl get svc -l='app.kubernetes.io/name=mysql' -o go-template='{{range .items}}{{range.spec.ports}}{{if .nodePort}}{{.nodePort}}{{"\n"}}{{end}}{{end}}{{end}}'

kubectl get service mysql -o yaml | sed 's/nodePort: 30367/nodePort: 30000/g' | kubectl replace -f -
```

Install users-DB on kind cluster:
```
kubectl apply -f ../kube/
```

Add user:

```
curl -d @sample_user.json -X POST http://localhost:30001/api/users
```

Create operator skeleton

```
mkdir operator
cd operator
operator-sdk init --domain example.com --repo github.com/mateuszdyminski/operator
```

Create a UsersDB Operator:

```
operator-sdk create api --group cache --version v1alpha1 --kind UsersDB --resource --controller
```

Set proper resources for UsersDB:

```
// UsersDBSpec defines the desired state of UsersDB
type UsersDBSpec struct {
	//+kubebuilder:validation:Minimum=0
	// Size is the size of the usersDB deployment
	Size int32 `json:"size"`
}

// UsersDBStatus defines the observed state of UsersDB
type UsersDBStatus struct {
	// Nodes are the names of the memcached pods
	Nodes []string `json:"nodes"`
}
```