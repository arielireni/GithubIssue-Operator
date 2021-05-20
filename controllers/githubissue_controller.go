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
	"os"
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	examplev1alpha1 "github.com/arielireni/example-operator/api/v1alpha1"

<<<<<<< HEAD
=======
	/* Imports for new github issue creation */
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// GitHubIssueReconciler reconciles a GitHubIssue object
type GitHubIssueReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

<<<<<<< HEAD
/* Repo structure declaration - all data fields for getting a repo's issues list */
type Repo struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

=======
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071
/* Issue structure declaration - all data fields for a new github issue submission */
type Issue struct {
	Title               string `json:"title"`
	Description         string `json:"body"`
	Number              int    `json:"number"`
	State               string `json:"state,omitempty"`
	LastUpdateTimestamp string `json:"updated_at,omitempty"`
}

/* Details structure declaration - all owner's details */
type Details struct {
	ApiURL string
	Token  string
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
<<<<<<< HEAD
	log.Info("Performs Reconcilation")
=======
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071

	/* Get the object from the api request */

	ghIssue := examplev1alpha1.GitHubIssue{}
	err := r.Client.Get(ctx, req.NamespacedName, &ghIssue) // fetch the k8s github object
<<<<<<< HEAD
=======
	log.Info("Enter logic")
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071

	if err != nil {
		/* Check if we got NotExist (404) error */
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		/* Any other error */
		return ctrl.Result{}, err
	}

	log.Info("got the gh issue from api server", "gh-issue", ghIssue)

<<<<<<< HEAD
	/* Create a github request and create github issues by interacting with github api */

	splittedRepo := strings.Split(ghIssue.Spec.Repo, "/")
	owner := splittedRepo[1]
	repo := splittedRepo[0]
	repoData := Repo{Owner: owner, Repo: repo}

=======
	/* Create a github request & create github issues by interacting with github api */

	repo := ghIssue.Spec.Repo
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071
	title := ghIssue.Spec.Title
	body := ghIssue.Spec.Description
	issueData := Issue{Title: title, Description: body}

<<<<<<< HEAD
	apiURL := "https://api.github.com/repos/" + ghIssue.Spec.Repo + "/issues?state=all"
	token := os.Getenv("TOKEN")
	detailsData := Details{ApiURL: apiURL, Token: token}

	issue := isIssueExist(&repoData, &issueData, &detailsData)
=======
	apiURL := "https://api.github.com/repos/" + repo + "/issues"
	token := os.Getenv("TOKEN")
	detailsData := Details{ApiURL: apiURL, Token: token}

	index := isIssueExist(issueData, detailsData)

	if index == -1 {
		createNewIssue(issueData, detailsData)
	} else {
		editIssue(issueData, detailsData, index)
	}
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071

	/* Implementation of deletion behaviour */
	finalizerName := "example.training.redhat.com/finalizer"
	// examine DeletionTimestamp to determine if object is under deletion
	if ghIssue.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(ghIssue.GetFinalizers(), finalizerName) {
			controllerutil.AddFinalizer(&ghIssue, finalizerName)
			if err := r.Update(ctx, &ghIssue); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(ghIssue.GetFinalizers(), finalizerName) {
			// our finalizer is present, so lets handle any external dependency
<<<<<<< HEAD
			if issue != nil {
				if err := r.deleteExternalResources(&issueData, issue, &detailsData); err != nil {
					// if fail to delete the external dependency here, return with error
					// so that it can be retried
					return ctrl.Result{}, err
				}
			}

=======
			if err := r.deleteExternalResources(issueData, detailsData, index); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071
			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&ghIssue, finalizerName)
			if err := r.Update(ctx, &ghIssue); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

<<<<<<< HEAD
	if issue == nil {
		issue = createNewIssue(&issueData, &detailsData)
	} else {
		//editIssue(issueData, detailsData, allIssues, index)
		if issueData.Description != issue.Description {
			editIssue(&issueData, issue, &detailsData)
		}
	}

	/* Update the state of the issue instance by the real github issue state */
=======
	/* Update the state of the issue instance by the real github issue state */
	issueData = updateIssueStatus(issueData, detailsData)
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071

	/* Update the k8s status with the real github issue state */
	patch := client.MergeFrom(ghIssue.DeepCopy())

<<<<<<< HEAD
	ghIssue.Status.State = issue.State
	ghIssue.Status.LastUpdateTimestamp = issue.LastUpdateTimestamp
=======
	ghIssue.Status.State = issueData.State
	ghIssue.Status.LastUpdateTimestamp = issueData.LastUpdateTimestamp
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071

	err = r.Client.Status().Patch(ctx, &ghIssue, patch)

	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

<<<<<<< HEAD
// SetupWithManager sets up the controller with the Manager.
func (r *GitHubIssueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1alpha1.GitHubIssue{}).
		Complete(r)
}

/* Functions to handle deletion with finalizer */
func (r *GitHubIssueReconciler) deleteExternalResources(issueData *Issue, issue *Issue, detailsData *Details) error {
=======
/* Functions to handle deletion with finalizer */
func (r *GitHubIssueReconciler) deleteExternalResources(issueData Issue, detailsData Details, index int) error {
	all_issues := createIssuesArray(detailsData.ApiURL)
	issue := all_issues[index]
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071
	issueApiURL := detailsData.ApiURL + "/" + fmt.Sprint(issue.Number)
	issue.State = "closed"
	jsonData, _ := json.Marshal(issue)

	/* Now update */
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", issueApiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Response code is is %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		// print body as it may contain hints in case of errors
		fmt.Println(string(body))
		log.Fatal(err)
	}
	return nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

<<<<<<< HEAD
/* Creates new issue with issueData's fields */
func createNewIssue(issueData *Issue, detailsData *Details) *Issue {
=======
// SetupWithManager sets up the controller with the Manager.
func (r *GitHubIssueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1alpha1.GitHubIssue{}).
		Complete(r)
}

func createNewIssue(issueData Issue, detailsData Details) {
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071
	apiURL := detailsData.ApiURL
	// make it json
	jsonData, _ := json.Marshal(issueData)
	// creating client to set custom headers for Authorization
	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("fatal error")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Response code is is %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		// print body as it may contain hints in case of errors
		fmt.Println(string(body))
		log.Fatal(err)
	}
<<<<<<< HEAD
	var issue *Issue
	issueBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(issueBody, &issue)
	return issue
}

/* Checks if the input issue exists. If yes, we will return issue, and nil otherwise */
func isIssueExist(repoData *Repo, issueData *Issue, detailsData *Details) *Issue {
	/* API request for all repository's issues */
	jsonData, _ := json.Marshal(&repoData)
	// creating client to set custom headers for Authorization
	client := &http.Client{}
	req, _ := http.NewRequest("GET", detailsData.ApiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
=======
}

/* Creates an array which contains all repository's issues */
func createIssuesArray(apiURL string) []Issue {
	/* API request for all repository's issues */
	resp, err := http.Get(apiURL)
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
<<<<<<< HEAD
	body, _ := ioutil.ReadAll(resp.Body)
	// print body as it may contain hints in case of errors
	fmt.Println(string(body))

	/* Create array with all repository's issues */
	var allIssues []Issue
	err = json.Unmarshal(body, &allIssues)

	/* If the issue exists, return it. Otherwise, return nil */
	for _, issue := range allIssues {
		if issue.Title == issueData.Title {
			return &issue
		}
	}
	return nil
}

/* Edits issue's description, to be equal to issueData's description */
func editIssue(issueData *Issue, issue *Issue, detailsData *Details) {
	issue.Description = issueData.Description
	issueApiURL := detailsData.ApiURL + "/" + fmt.Sprint(issue.Number)
	jsonData, _ := json.Marshal(issue)

	/* Now update */
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", issueApiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Response code is is %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		// print body as it may contain hints in case of errors
		fmt.Println(string(body))
		log.Fatal(err)
	}
=======

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Response code is is %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		// print body as it may contain hints in case of errors
		fmt.Println(string(body))
		log.Fatal(err)
	}

	/* Create array with all repository's issues */
	var all_issues []Issue
	issues_body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(issues_body, &all_issues)

	return all_issues
}

/* Checks if the input issue exists.
If yes, we will return the index of issue at the issues array, and -1 otherwise */
func isIssueExist(issueData Issue, detailsData Details) int {

	apiURL := detailsData.ApiURL
	all_issues := createIssuesArray(apiURL)

	/* Loop over all repository's issues */
	for i := 0; i < len(all_issues); i++ {
		if all_issues[i].Title == issueData.Title {
			return i
		}
	}

	/* If we have reached this point, then the issue doesn't exist yet */
	return -1
}

func editIssue(issueData Issue, detailsData Details, index int) {

	all_issues := createIssuesArray(detailsData.ApiURL)
	issue := all_issues[index]
	issueApiURL := detailsData.ApiURL + "/" + fmt.Sprint(issue.Number)
	if issue.Description != issueData.Description {
		issue.Description = issueData.Description
		jsonData, _ := json.Marshal(issue)

		/* Now update */
		client := &http.Client{}
		req, _ := http.NewRequest("PATCH", issueApiURL, bytes.NewReader(jsonData))
		req.Header.Set("Authorization", "token "+detailsData.Token)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Response code is is %d\n", resp.StatusCode)
			body, _ := ioutil.ReadAll(resp.Body)
			// print body as it may contain hints in case of errors
			fmt.Println(string(body))
			log.Fatal(err)
		}
	}

}

/* NOTE: we do not want to do the update in the updateIssue function since someone can also change the
issue status from the outside */

func updateIssueStatus(issueData Issue, detailsData Details) Issue {
	all_issues := createIssuesArray(detailsData.ApiURL)

	/* Loop over all repository's issues */
	for i := 0; i < len(all_issues); i++ {
		if all_issues[i].Title == issueData.Title {
			return all_issues[i]
		}
	}
	return issueData
>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071
}
