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
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/go-logr/logr"
	teleportv1 "github.com/zoetrope/teleport-plus/api/v1"
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

	// your logic here
	var github teleportv1.GitHub
	err := r.Get(ctx, req.NamespacedName, &github)
	if err != nil {
		log.Error(err, "unable to fetch GitHub")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconcile a github resource: ", "name", github.GetName())

	// if github.DeletionTimestamp.IsZero() {

	// }

	res := &GithubConnectorV3{
		Kind:    "github",
		Version: "v3",
		Metadata: Metadata{
			Name: github.GetName(),
		},
		Spec: github.Spec,
	}
	out, err := yaml.Marshal(res)
	if err != nil {
		log.Error(err, "unable to marshal github resource")
		return ctrl.Result{}, err
	}
	tmpfile, err := ioutil.TempFile("", "github.yaml")
	if err != nil {
		log.Error(err, "unable to create temp file")
		return ctrl.Result{}, err
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	if _, err := tmpfile.Write(out); err != nil {
		log.Error(err, "unable to write to temp file")
		return ctrl.Result{}, err
	}
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "/tctl", "-c", "/etc/teleport/teleport.yaml", "create", "-f", tmpfile.Name())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Error(err, "failed to create github connector", "stdout", string(stdout.Bytes()), "stderr", string(stderr.Bytes()), "yaml", string(out))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *GitHubReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&teleportv1.GitHub{}).
		Complete(r)
}
