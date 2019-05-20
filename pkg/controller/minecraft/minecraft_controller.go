package minecraft

import (
	"context"

	operatorv1alpha1 "github.com/softica/minecraft-operator/pkg/apis/operator/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_minecraft")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Minecraft Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMinecraft{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("minecraft-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Minecraft
	err = c.Watch(&source.Kind{Type: &operatorv1alpha1.Minecraft{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Minecraft
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorv1alpha1.Minecraft{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMinecraft{}

// ReconcileMinecraft reconciles a Minecraft object
type ReconcileMinecraft struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Minecraft object and makes changes based on the state read
// and what is in the Minecraft.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMinecraft) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Minecraft")

	// Fetch the Minecraft instance
	instance := &operatorv1alpha1.Minecraft{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}




	ingress := newIngressForCR(instance)

        // Set Minecraft instance as the owner and controller of the ingress
        if err := controllerutil.SetControllerReference(instance, ingress, r.scheme); err != nil {
                return reconcile.Result{}, err
        }

        // Check if this Ingress already exists
        ingfound := &extensionsv1beta1.Ingress{}
        err = r.client.Get(context.TODO(), types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, ingfound)
        if err != nil && errors.IsNotFound(err) {
            err = r.client.Create(context.TODO(), ingress)
            if err != nil {
                return reconcile.Result{}, err
            }
        } else if err != nil {
            return reconcile.Result{}, err
        }






        // Define a new Service object
        service := newServiceForCR(instance)

        // Set Minecraft instance as the owner and controller of the service
        if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
                return reconcile.Result{}, err
        }

        // Check if this Service already exists
        srvfound := &corev1.Service{}
        err = r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, srvfound)
        if err != nil && errors.IsNotFound(err) {
            err = r.client.Create(context.TODO(), service)
            if err != nil {
                return reconcile.Result{}, err
            }
        } else if err != nil {
            return reconcile.Result{}, err
        }






	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set Minecraft instance as the owner and controller of the pod
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil


}

// newIngressForCR returns a minecraft ingress with the same name/namespace as the cr
func newIngressForCR(cr *operatorv1alpha1.Minecraft) *extensionsv1beta1.Ingress {
	// https://godoc.org/k8s.io/api/extensions/v1beta1
        labels := map[string]string{
                "app": cr.Name,
                "version": cr.Spec.Version,
                "uela": cr.Spec.Uela,
        }
        return &extensionsv1beta1.Ingress{
                ObjectMeta: metav1.ObjectMeta{
                        Name:      cr.Name + "-service",
                        Namespace: cr.Namespace,
                        Labels:    labels,
                },
		Spec: extensionsv1beta1.IngressSpec{
		// https://godoc.org/k8s.io/api/extensions/v1beta1#IngressSpec
			Rules: []extensionsv1beta1.IngressRule{
				{
					Host: "static-host.com",  //TODO : Need to get hostname from CR
				},
			},
		// https://godoc.org/k8s.io/api/extensions/v1beta1#IngressRule
		},
	}
}

// newServiceForCR returns a minecraft service with the same name/namespace as the cr
func newServiceForCR(cr *operatorv1alpha1.Minecraft) *corev1.Service {
        labels := map[string]string{
                "app": cr.Name,
                "version": cr.Spec.Version,
                "uela": cr.Spec.Uela,
        }
        return &corev1.Service{
		// https://godoc.org/k8s.io/api/core/v1#Service
                ObjectMeta: metav1.ObjectMeta{
                        Name:      cr.Name + "-service",
                        Namespace: cr.Namespace,
                        Labels:    labels,
                },
		Spec: corev1.ServiceSpec{
			// https://godoc.org/k8s.io/api/core/v1#ServiceSpec
//			ServiceType: ClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: "minecraft",
					Port: 25565,
//					TargetPort: 25565,
				},
			},
			Selector: labels,
		},
	}
}

// newPodForCR returns a minecraft pod with the same name/namespace as the cr
func newPodForCR(cr *operatorv1alpha1.Minecraft) *corev1.Pod {
	var envVars []corev1.EnvVar
	envVars = []corev1.EnvVar {
		corev1.EnvVar {
		Name: "EULA",
		Value: cr.Spec.Uela,
		},
	}

	var fsType int64
	fsType = int64(1000)

	labels := map[string]string{
		"app": cr.Name,
                "version": cr.Spec.Version,
		"uela": cr.Spec.Uela,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			SecurityContext: &corev1.PodSecurityContext{
				FSGroup: &fsType,
				RunAsUser: &fsType,
			},
			Containers: []corev1.Container{
				{
					Name:    "minecraft",
					Image:   "hoeghh/minecraft:" + cr.Spec.Version,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 25565,
							Name:         "minecraft",
						},
					},
					Env: envVars,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "minecraft-volume",
							MountPath: "/minecraft-data",
							ReadOnly: false,
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "minecraft-volume",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
 	                                        	Path: "/data",
                                                },
					},
				},
			},
		},
	}
}

