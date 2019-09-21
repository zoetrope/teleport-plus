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
	teleportv1 "github.com/zoetrope/teleport-plus/api/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

// GitHubReconciler reconciles a GitHub object
type GitHubReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=teleport.gravitational.com,resources=githubs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=teleport.gravitational.com,resources=githubs/status,verbs=get;update;patch

func (r *GitHubReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("github", req.NamespacedName)

	gh := new(teleportv1.GitHub)
	err := r.Get(ctx, req.NamespacedName, gh)
	if err != nil {
		if !apierrs.IsNotFound(err) {
			log.Error(err, "unable to fetch GitHub")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	ns, err := ownNamespace()
	if err != nil {
		log.Error(err, "unable to get own namespace")
		return ctrl.Result{}, err
	}
	if ns != gh.Namespace {
		return ctrl.Result{}, err
	}

	log.Info("Reconcile a github resource: ", "name", gh.GetName())

	if gh.DeletionTimestamp.IsZero() {
		if !containsString(gh.Finalizers, finalizerName) {
			err = r.addFinalizer(ctx, log, gh)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		err = createOrUpdate(ctx, log, gh)
		if err == nil {
			err = r.updateStatus(ctx, log, gh, ConditionRegistered, "")
		} else {
			_ = r.updateStatus(ctx, log, gh, ConditionFailed, err.Error())
		}
		return ctrl.Result{}, err
	} else if containsString(gh.Finalizers, finalizerName) {
		err = r.remove(ctx, log, gh)
		if err != nil {
			_ = r.updateStatus(ctx, log, gh, ConditionFailed, err.Error())
		}
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *GitHubReconciler) addFinalizer(ctx context.Context, log logr.Logger, gh *teleportv1.GitHub) error {
	log.Info("add finalizer")
	gh2 := gh.DeepCopy()
	gh2.Finalizers = append(gh2.Finalizers, finalizerName)
	patch := client.MergeFrom(gh)
	if err := r.Patch(ctx, gh2, patch); err != nil {
		log.Error(err, "failed to add finalizer", "name", finalizerName)
		return err
	}
	return nil
}

func (r *GitHubReconciler) remove(ctx context.Context, log logr.Logger, gh *teleportv1.GitHub) error {
	log.Info("remove resource")
	stdout, stderr, err := execTctl(ctx, "rm", "github/"+gh.Name)
	if err != nil {
		log.Error(err, "failed to remove github connector", "stdout", string(stdout), "stderr", string(stderr), "name", string(gh.Name))
		return err
	}
	gh2 := gh.DeepCopy()
	gh2.Finalizers = removeString(gh2.Finalizers, finalizerName)
	patch := client.MergeFrom(gh)
	if err := r.Patch(ctx, gh2, patch); err != nil {
		log.Error(err, "failed to remove finalizer", "name", finalizerName)
		return err
	}
	return nil
}

func createOrUpdate(ctx context.Context, log logr.Logger, gh *teleportv1.GitHub) error {
	log.Info("create or update resource")
	res := &GithubConnectorV3{
		Kind:    "github",
		Version: "v3",
		Metadata: Metadata{
			Name: gh.GetName(),
		},
		Spec: gh.Spec,
	}
	out, err := yaml.Marshal(res)
	if err != nil {
		log.Error(err, "unable to marshal github resource")
		return err
	}
	tmpfile, err := ioutil.TempFile("", "github.yaml")
	if err != nil {
		log.Error(err, "unable to create temp file")
		return err
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	if _, err := tmpfile.Write(out); err != nil {
		log.Error(err, "unable to write to temp file")
		return err
	}
	stdout, stderr, err := execTctl(ctx, "create", "-f", tmpfile.Name())
	if err != nil {
		log.Error(err, "failed to create github connector", "stdout", string(stdout), "stderr", string(stderr), "yaml", string(out))
		return err
	}
	return nil
}

func (r *GitHubReconciler) updateStatus(ctx context.Context, logger logr.Logger, gh *teleportv1.GitHub, cond string, reason string) error {
	nowTime := metav1.NewTime(time.Now())

	gh2 := gh.DeepCopy()
	gh2.Status.Condition = cond
	gh2.Status.Reason = reason
	gh2.Status.LastTransitionTime = &nowTime
	patch := client.MergeFrom(gh)

	if err := r.Status().Patch(ctx, gh2, patch); err != nil {
		logger.Error(err, "unable to update GitHub status")
		return err
	}
	return nil
}

func (r *GitHubReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&teleportv1.GitHub{}).
		Complete(r)
}
