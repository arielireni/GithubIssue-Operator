package controllers

import (
	"context"
	"fmt"
	examplev1alpha1 "github.com/arielireni/example-operator/api/v1alpha1"
	"github.com/arielireni/example-operator/controllers/clients"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

// Create issue tests
func TestSuccessfulCreate(t *testing.T) {
	// Given a valid ghIssue
	fakeClient := clients.NewFakeClient([]clients.Issue{}, true, nil)

	// Reconciler
	s := scheme.Scheme
	examplev1alpha1.AddToScheme(s)
	fakeK8sClient := newFakeK8sClient()
	r := GitHubIssueReconciler{
		Client:      fakeK8sClient,
		Log:         ctrl.Log,
		Scheme:      s,
		ClientFrame: fakeClient,
	}

	// When creating a real issue
	_, err := r.Reconcile(context.Background(), ctrl.Request{})

	// Then reconcile returns ctrl.Result{} and no error
	if err != nil {
		t.Errorf("Expected nil but got error")
	}
}

func TestFailedCreate(t *testing.T) {
	// Given we fail to create a real issue
	testErr := fmt.Errorf("TestFailedCreate error")
	fakeClient := clients.NewFakeClient([]clients.Issue{}, false, testErr)
	// Reconciler
	s := scheme.Scheme
	examplev1alpha1.AddToScheme(s)
	fakeK8sClient := newFakeK8sClient()
	r := GitHubIssueReconciler{
		Client:      fakeK8sClient,
		Log:         ctrl.Log,
		Scheme:      s,
		ClientFrame: fakeClient,
	}
	_, err := r.Reconcile(context.Background(), ctrl.Request{})
	// Then reconcile returns ctrl.Result{} and error
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func newFakeK8sClient() client.Client {
	issue := examplev1alpha1.GitHubIssue{
		Spec: examplev1alpha1.GitHubIssueSpec{
			Repo:  "arielireni/Issues-Example",
			Title: "title1",
		},
	}
	objects := []runtime.Object{issue.DeepCopyObject()}
	fakeK8sClient := fake.NewClientBuilder().WithRuntimeObjects(objects...).Build()
	return fakeK8sClient
}

// Edit issue tests
func TestSuccessfulEdit(t *testing.T) {
	t.Skip("unimplemented")
}

func TestFailedEdit(t *testing.T) {
	t.Skip("unimplemented")
}

// Close issue tests
func TestSuccessfulClose(t *testing.T) {
	t.Skip("unimplemented")
}

func TestFailedClose(t *testing.T) {
	t.Skip("unimplemented")
}
