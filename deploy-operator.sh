kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/crds/operator_v1alpha1_minecraft_crd.yaml
kubectl create -f deploy/operator.yaml

