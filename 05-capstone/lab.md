# Project milestones

This is an independent project, not a copy-and-paste lab. Use the earlier module examples to implement each milestone. Run commands against your own kind-capstone context and write down what you expect before making changes.

## 1. Draw the request path

Decide how the browser reaches the frontend, how frontend requests reach the API and how the API authenticates to the database. Define record fields and endpoint responses.

**Q1.** Why should the browser not call `http://api:8080` directly?

**Done when:** you can trace a request and identify the responsible Service at each internal hop.

## 2. Build and test

Implement the frontend, API and SQL persistence; validate input and test behavior. Create images and CI. Use environment variables for database connection details and a Secret for the password.

**Q2.** Which configuration belongs in a ConfigMap, and which must be kept out of Git?

**Done when:** a submitted record can be retrieved through the API and invalid input is rejected.

## 3. Deploy your package

Create the separate three-node cluster, choose Helm or Kustomize, and deploy two application replicas plus the database and PVC. Inspect Pods, events and Services. Access the frontend with port-forward.

**Q3.** What happens to your record when the database Pod is deleted? What happens when its PVC or the whole kind cluster is deleted?

**Done when:** the browser can save/list records and a Pod replacement preserves a saved record.

## 4. Let GitOps manage it

Install exactly one GitOps controller for this application. Configure the repository, revision, path and namespace. Adopt or cleanly replace your initial installation without overlapping ownership; explain the procedure.

**Q4.** How can you prove which Git revision is deployed? Why might a successfully fetched revision still fail to become healthy?

**Done when:** your chosen controller reports success and repairs a manually changed replica count.

## 5. Give the final demonstration

Show the page, save a note, commit a ConfigMap color change and observe it in the browser without rebuilding an image. Publish and promote a new image, then diagnose and recover from one deliberate failure.

**Q5.** How does your frontend observe the ConfigMap change? What delay or restart should the audience expect?

**Q6.** Why must the permanent incident repair be recorded in Git?

**Done when:** all evidence in requirements.md is collected and the deployment is healthy again.

## Cleanup

Document how you stop GitOps management before deleting application resources. Explain which operations preserve database storage and which destroy it. Export any data you need before deleting the cluster.

Submit answers Q1–Q6 with your project. The instructor solution includes reference answers and a separate demonstration runbook.
