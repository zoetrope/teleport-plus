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

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	teleportv1 "github.com/zoetrope/teleport-plus/api/v1"
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

	// your logic here
	var github teleportv1.GitHub
	err := r.Get(ctx, req.NamespacedName, &github)
	if err != nil {
		log.Error(err, "unable to fetch GitHub")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconcile a github resource: ", github.GetName())

	return ctrl.Result{}, nil
}

func (r *GitHubReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&teleportv1.GitHub{}).
		Complete(r)
}
