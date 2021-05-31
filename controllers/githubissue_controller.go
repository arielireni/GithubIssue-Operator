/*
Copyright 2021.

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
	"github.com/arielireni/example-operator/controllers/clients"

	examplev1alpha1 "github.com/arielireni/example-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// GitHubIssueReconciler reconciles a GitHubIssue object
type GitHubIssueReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	ClientFrame clients.ClientFrame
}

//+kubebuilder:rbac:groups=example.training.redhat.com,resources=githubissues,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.training.redhat.com,resources=githubissues/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.training.redhat.com,resources=githubissues/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GitHubIssue object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *GitHubIssueReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := r.Log.WithValues("name-of-gh-issue", req.NamespacedName)
	log.Info("Performs Reconciliation")

	// Get the object from the api request
	ghIssue := examplev1alpha1.GitHubIssue{}
	err := r.Client.Get(ctx, req.NamespacedName, &ghIssue) // Fetch the k8s clients object

	if err != nil {
		// Check if we got NotExist (404) error
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Any other error
		return ctrl.Result{}, err
	}

	log.Info("got the gh issue from api server", "gh-issue", ghIssue)

	// Create a clients request and create clients issues by interacting with clients api
	repoData, issueData, detailsData := r.ClientFrame.InitDataStructs(ghIssue.Spec.Repo, ghIssue.Spec.Title, ghIssue.Spec.Description)
	issue, returnErr := r.ClientFrame.FindIssue(repoData, issueData, detailsData)

	// If we encounter an error we will warn about it and continue running
	if returnErr.ErrorCode != nil {
		log.Info(returnErr.Message)
	}

	// Deletion behavior
	stopReconcile, delErr := r.delete(&ghIssue, ctx, issue, issueData, detailsData)
	if stopReconcile == true {
		return ctrl.Result{}, delErr
	}

	// Create new issue or update if needed
	if issue == nil {
		issue, returnErr = r.ClientFrame.CreateIssue(issueData, detailsData)
		if returnErr.ErrorCode != nil {
			log.Info(returnErr.Message)
		}
	} else {
		if (issueData.Description != issue.Description) && (issue.State != "closed") {
			returnErr = r.ClientFrame.EditIssue(issueData, issue, detailsData)
			if returnErr.ErrorCode != nil {
				log.Info(returnErr.Message)
			}
		}
	}

	// Update the state of the issue instance by the real clients issue state

	// Update the k8s status with the real clients issue state
	patch := client.MergeFrom(ghIssue.DeepCopy())

	ghIssue.Status.State = issue.State
	ghIssue.Status.LastUpdateTimestamp = issue.LastUpdateTimestamp

	err = r.Client.Status().Patch(ctx, &ghIssue, patch)

	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GitHubIssueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1alpha1.GitHubIssue{}).
		Complete(r)
}

// Implementation of deletion behavior, will return true if we need to stop reconcilation
func (r *GitHubIssueReconciler) delete(ghIssue *examplev1alpha1.GitHubIssue, ctx context.Context, issue *clients.Issue, issueData *clients.Issue, detailsData *clients.Details) (bool, error) {
	finalizerName := "example.training.redhat.com/finalizer"
	// examine DeletionTimestamp to determine if object is under deletion
	if ghIssue.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(ghIssue.GetFinalizers(), finalizerName) {
			controllerutil.AddFinalizer(ghIssue, finalizerName)
			if err := r.Update(ctx, ghIssue); err != nil {
				return true, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(ghIssue.GetFinalizers(), finalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if issue != nil {
				if err := r.deleteExternalResources(issueData, issue, detailsData); err != nil {
					// if fail to delete the external dependency here, return with error
					// so that it can be retried
					return true, err
				}
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(ghIssue, finalizerName)
			if err := r.Update(ctx, ghIssue); err != nil {
				return true, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return true, nil
	}
	return false, nil
}

// Functions to handle deletion with finalizer
func (r *GitHubIssueReconciler) deleteExternalResources(issueData *clients.Issue, issue *clients.Issue, detailsData *clients.Details) error {
	return r.ClientFrame.CloseIssue(issueData, issue, detailsData).ErrorCode
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
