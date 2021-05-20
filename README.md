# GithubIssue Operator Task
## Goal
- Implement an operator that will allow creating, editing and closing a github issue.

## The Reconciliation behaviour
- Fetch the k8s github object by the req.NamespacedName.
- Fetch all the github issues of the repository. 
- Find one with the exact title we need.
  - If it doesn't exist we will create a github issue with the title and description.
  - If it exists we will update the description.
- Update the k8s status with the real github issue state.
- A delete of the k8s object, triggers the github issue to be closed.
- Resyncs every 1 minute.


>>>>>>> 2344754eb238fc929f6c20c6e493910cf19d1071
