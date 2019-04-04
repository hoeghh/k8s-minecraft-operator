# How-to
First off, to deploy this operator, simply run the script ```./deploy-operator.sh``` and wait for the operator to be ready.

Then run the script ```./deploy-minecraft-cr.sh```. This will create a Custom Resource of type minecraft, and the operator will react and reconcile. Soon, you should have a minecraft server running in a pod. 

Run ```kubectl port-forward example-minecraft-pod 25565:25565``` to forward local trafic to the pod on port 25565 (default minecraft port). Now start minecraft and chose multiplayer. Add a server and choose localhost for the hostname. Connect to the server and you are now playing minecraft on a Operator managed Minecraft deployment in Kubernetes.

Guide :
https://github.com/operator-framework/operator-sdk

Good example :
https://github.com/operator-framework/operator-sdk-samples/blob/master/memcached-operator/pkg/controller/memcached/memcached_controller.go

Complete data overview:
https://godoc.org/k8s.io/api/core/v1

Edit controller :
pkg/controller/minecraft/minecraft_controller.go
