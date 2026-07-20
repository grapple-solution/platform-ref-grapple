#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

echo "Registering version-compatible Modelplane Custom Resource Definitions (CRDs)..."
kubectl apply -f "${SCRIPT_DIR}/modelplane-crds.yaml"

echo "Modelplane CRDs successfully registered."
