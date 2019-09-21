/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	teleportv1 "github.com/zoetrope/teleport-plus/api/v1"
)

// TeleportResourceReconciler reconciles a TeleportResource object
type TeleportResourceReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=teleport.gravitational.com,resources=teleportresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=teleport.gravitational.com,resources=teleportresources/status,verbs=get;update;patch

func (r *TeleportResourceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("teleportresource", req.NamespacedName)

	res := new(teleportv1.TeleportResource)
	err := r.Get(ctx, req.NamespacedName, res)
	if err != nil {
		if !apierrs.IsNotFound(err) {
			log.Error(err, "unable to fetch TeleportResource")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconcile a teleport resource: ", "name", res.GetName())

	if res.DeletionTimestamp.IsZero() {
		if !containsString(res.Finalizers, finalizerName) {
			err = r.addFinalizer(ctx, log, res)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		err = r.createOrUpdate(ctx, log, res)
		if err == nil {
			err = r.updateStatus(ctx, log, res, ConditionRegistered, "")
		} else {
			_ = r.updateStatus(ctx, log, res, ConditionFailed, err.Error())
		}
		return ctrl.Result{}, err
	} else if containsString(res.Finalizers, finalizerName) {
		err = r.remove(ctx, log, res)
		if err != nil {
			_ = r.updateStatus(ctx, log, res, ConditionFailed, err.Error())
		}
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *TeleportResourceReconciler) addFinalizer(ctx context.Context, log logr.Logger, res *teleportv1.TeleportResource) error {
	log.Info("add finalizer")
	res2 := res.DeepCopy()
	res2.Finalizers = append(res2.Finalizers, finalizerName)
	patch := client.MergeFrom(res)
	if err := r.Patch(ctx, res2, patch); err != nil {
		log.Error(err, "failed to add finalizer", "name", finalizerName)
		return err
	}
	return nil
}

func (r *TeleportResourceReconciler) remove(ctx context.Context, log logr.Logger, res *teleportv1.TeleportResource) error {
	log.Info("remove resource")

	var meta TeleportMeta
	err := yaml.Unmarshal([]byte(res.Spec.Data), &meta)
	if err != nil {
		return err
	}

	resName := meta.Kind + "/" + meta.Metadata.Name
	stdout, stderr, err := execTctl(ctx, "rm", resName)
	if err != nil {
		log.Error(err, "failed to remove github connector", "stdout", string(stdout), "stderr", string(stderr), "name", resName)
		return err
	}
	res2 := res.DeepCopy()
	res2.Finalizers = removeString(res2.Finalizers, finalizerName)
	patch := client.MergeFrom(res)
	if err := r.Patch(ctx, res2, patch); err != nil {
		log.Error(err, "failed to remove finalizer", "name", finalizerName)
		return err
	}
	return nil
}

func (r *TeleportResourceReconciler) createOrUpdate(ctx context.Context, log logr.Logger, res *teleportv1.TeleportResource) error {
	log.Info("create or update resource")
	tmpfile, err := ioutil.TempFile("", "teleport-resource.yaml")
	if err != nil {
		log.Error(err, "unable to create temp file")
		return err
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	if _, err := tmpfile.Write([]byte(res.Spec.Data)); err != nil {
		log.Error(err, "unable to write to temp file")
		return err
	}
	stdout, stderr, err := execTctl(ctx, "create", "-f", tmpfile.Name())
	if err != nil {
		log.Error(err, "failed to create teleport resource", "stdout", string(stdout), "stderr", string(stderr), "yaml", res.Spec.Data)
		return err
	}
	return nil
}

func (r *TeleportResourceReconciler) updateStatus(ctx context.Context, logger logr.Logger, res *teleportv1.TeleportResource, cond string, reason string) error {
	nowTime := metav1.NewTime(time.Now())

	res2 := res.DeepCopy()
	res2.Status.Condition = cond
	res2.Status.Reason = reason
	res2.Status.LastTransitionTime = &nowTime
	patch := client.MergeFrom(res)

	if err := r.Status().Patch(ctx, res, patch); err != nil {
		logger.Error(err, "unable to update status")
		return err
	}
	return nil
}

func (r *TeleportResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&teleportv1.TeleportResource{}).
		Complete(r)
}
