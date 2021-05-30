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

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	examplev1alpha1 "github.com/arielireni/example-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	gerrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// GitHubIssueReconciler reconciles a GitHubIssue object
type GitHubIssueReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Repo structure declaration - all data fields for getting a repo's issues list
type Repo struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

// Issue structure declaration - all data fields for a new github issue submission
type Issue struct {
	Title               string `json:"title"`
	Description         string `json:"body"`
	Number              int    `json:"number"`
	State               string `json:"state,omitempty"`
	LastUpdateTimestamp string `json:"updated_at,omitempty"`
}

// Details structure declaration - all owner's details
type Details struct {
	ApiURL string
	Token  string
}

// Error structure decleration - errors with messages
type Error struct {
	ErrorCode error
	Message   string
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
	log.Info("Performs Reconcilation")

	// Get the object from the api request
	ghIssue := examplev1alpha1.GitHubIssue{}
	err := r.Client.Get(ctx, req.NamespacedName, &ghIssue) // Fetch the k8s github object

	if err != nil {
		// Check if we got NotExist (404) error
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Any other error
		return ctrl.Result{}, err
	}

	log.Info("got the gh issue from api server", "gh-issue", ghIssue)

	// Create a github request and create github issues by interacting with github api
	repoData, issueData, detailsData := initDataStructs(ghIssue.Spec.Repo, ghIssue.Spec.Title, ghIssue.Spec.Description)
	issue, returnErr := findIssue(&repoData, &issueData, &detailsData)

	// If we encounter an error we will warn about it and continue running
	if returnErr != nil {
		log.Info("Got error while finding issue")
	}

	// Deletion behavior
	stopRec, delErr := delete(r, &ghIssue, ctx, issue, &issueData, &detailsData)
	if stopRec == true {
		return ctrl.Result{}, delErr
	}

	// Create new issue or update if needed
	if issue == nil {
		issue = createNewIssue(&issueData, &detailsData)
	} else {
		if (issueData.Description != issue.Description) && (issue.State != "closed") {
			editIssue(&issueData, issue, &detailsData)
		}
	}

	// Update the state of the issue instance by the real github issue state

	// Update the k8s status with the real github issue state
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

// Initialize issueData & detailsData
func initDataStructs(repo, title, body string) (Repo, Issue, Details) {
	splittedRepo := strings.Split(repo, "/")
	owner := splittedRepo[1]
	repoName := splittedRepo[0]
	repoData := Repo{Owner: owner, Repo: repoName}

	issueData := Issue{Title: title, Description: body}

	apiURL := "https://api.github.com/repos/" + repo + "/issues"
	token := os.Getenv("TOKEN")
	detailsData := Details{ApiURL: apiURL, Token: token}

	return repoData, issueData, detailsData
}

// Implementation of deletion behavior, will return true if we need to stop reconcilation
func delete(r *GitHubIssueReconciler, ghIssue *examplev1alpha1.GitHubIssue, ctx context.Context, issue *Issue, issueData *Issue, detailsData *Details) (bool, error) {
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
func (r *GitHubIssueReconciler) deleteExternalResources(issueData *Issue, issue *Issue, detailsData *Details) error {

	issueApiURL := detailsData.ApiURL + "/" + fmt.Sprint(issue.Number)
	issue.State = "closed"
	jsonData, _ := json.Marshal(issue)

	// Now update
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

// Creates new issue with issueData's fields
func createNewIssue(issueData *Issue, detailsData *Details) *Issue {
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
		// added at code review session
		gerrors.Wrap(err, "faild to create gh issue")
		// print body as it may contain hints in case of errors
		fmt.Println(string(body))

		log.Fatal(err)
	}
	var issue *Issue
	issueBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(issueBody, &issue)
	return issue
}

// Checks if the input issue exists. If yes, we will return issue, and nil otherwise
func findIssue(repoData *Repo, issueData *Issue, detailsData *Details) (*Issue, *Error) {
	apiURL := detailsData.ApiURL + "?state=all"
	/* API request for all repository's issues */
	jsonData, _ := json.Marshal(&repoData)
	// creating client to set custom headers for Authorization
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	if err != nil {
		typedErr := Error{ErrorCode: err, Message: "GET request from GitHub API faild"}
		return nil, &typedErr
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// print body as it may contain hints in case of errors
	fmt.Println(string(body))

	// Create array with all repository's issues
	var allIssues []Issue

	err = json.Unmarshal(body, &allIssues)
	if err != nil {
		typedErr := Error{ErrorCode: err, Message: "Unmarshal faild"}
		return nil, &typedErr
	}

	// If we found the issue, we'll return it. Otherwise, return nil
	for _, issue := range allIssues {
		if issue.Title == issueData.Title {
			return &issue, nil
		}
	}
	return nil, nil
}

// Edits issue's description, to be equal to issueData's description
func editIssue(issueData *Issue, issue *Issue, detailsData *Details) {
	issue.Description = issueData.Description
	issueApiURL := detailsData.ApiURL + "/" + fmt.Sprint(issue.Number)
	fmt.Printf("URL: " + issueApiURL)
	jsonData, _ := json.Marshal(issue)

	// Now update
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
