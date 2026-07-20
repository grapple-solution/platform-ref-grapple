#!/bin/bash
set -euo pipefail

echo "Registering version-compatible Modelplane Custom Resource Definitions (CRDs)..."

kubectl apply -f - <<EOF
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: modeldeployments.modelplane.ai
spec:
  group: modelplane.ai
  versions:
    - name: v1alpha1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          x-kubernetes-preserve-unknown-fields: true
  scope: Namespaced
  names:
    plural: modeldeployments
    singular: modeldeployment
    kind: ModelDeployment
    shortNames:
    - modeldeploy
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: modelservices.modelplane.ai
spec:
  group: modelplane.ai
  versions:
    - name: v1alpha1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          x-kubernetes-preserve-unknown-fields: true
  scope: Namespaced
  names:
    plural: modelservices
    singular: modelservice
    kind: ModelService
    shortNames:
    - modelsvc
EOF

echo "Modelplane CRDs successfully registered."
