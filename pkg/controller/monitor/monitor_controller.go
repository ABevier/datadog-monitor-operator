package monitor

import (
	"context"
	"strings"

	datadogv1alpha1 "github.com/abevier/datadog-monitor-operator/pkg/apis/datadog/v1alpha1"
	"github.com/abevier/datadog-monitor-operator/pkg/dd"
	"gopkg.in/zorkian/go-datadog-api.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_monitor")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Monitor Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMonitor{client: mgr.GetClient(), scheme: mgr.GetScheme(), datadogClient: dd.NewClient()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("monitor-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Monitor
	err = c.Watch(&source.Kind{Type: &datadogv1alpha1.Monitor{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Monitor
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &datadogv1alpha1.Monitor{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileMonitor implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMonitor{}

// ReconcileMonitor reconciles a Monitor object
type ReconcileMonitor struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client        client.Client
	scheme        *runtime.Scheme
	datadogClient *dd.DataDogClient
}

// Reconcile reads that state of the cluster for a Monitor object and makes changes based on the state read
// and what is in the Monitor.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMonitor) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Monitor")

	// Fetch the Monitor instance
	instance := &datadogv1alpha1.Monitor{}
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

	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if hasMonitorFinalizer(instance) {
			// Delete
			id := instance.Status.ID
			err := r.datadogClient.DeleteMonitor(id)
			if err != nil && !strings.HasPrefix(err.Error(), "API error 404 Not Found") {
				return reconcile.Result{}, err
			}

			removeMonitorFinalizer(instance)
			r.client.Update(context.Background(), instance)

			log.Info("Delete monitor", "id", id)
		}
	} else if instance.Status.ID != 0 {
		// Update
	} else {
		// Add
		id, err := r.datadogClient.CreateMonitor(&instance.Spec)
		if err != nil {
			//TODO: better error handling?
			return reconcile.Result{}, err
		}

		instance.Status.ID = id
		err = r.client.Update(context.Background(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		addMonitorFinalizer(instance)
		err = r.client.Update(context.Background(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		log.Info("Created monitor.", "id", id)
	}

	// // Define a new Pod object
	// pod := newPodForCR(instance)

	// // Set Monitor instance as the owner and controller
	// if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Check if this Pod already exists
	// found := &corev1.Pod{}
	// err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	// 	err = r.client.Create(context.TODO(), pod)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}

	// 	// Pod created successfully - don't requeue
	// 	return reconcile.Result{}, nil
	// } else if err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Pod already exists - don't requeue
	// reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	//return reconcile.Result{RequeueAfter: time.Second * 30}, nil
	return reconcile.Result{}, nil
}

func newMonitor() (int, error) {
	client := datadog.NewClient("secret", "secret")

	m := &datadog.Monitor{}
	m.SetName("test")
	m.SetQuery("test")

	m, err := client.CreateMonitor(m)
	if err != nil {
		return -1, err
	}

	return *m.Id, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *datadogv1alpha1.Monitor) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
