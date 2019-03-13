kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/crds/operator_v1alpha1_minecraft_crd.yaml
kubectl delete -f deploy/operator.yaml

