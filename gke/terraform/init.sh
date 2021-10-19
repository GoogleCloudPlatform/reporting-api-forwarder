#!/bin/bash
#
# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

function check_commands () {
  if ! type gcloud > /dev/null; then
    echo "Install gcloud command following the instructions."
    echo "https://cloud.google.com/sdk/docs/install"
    exit 1
  fi

  if ! type kubectl > /dev/null; then
    echo "Install kubectl following the instructions."
    echo "https://kubernetes.io/docs/tasks/tools/"
    echo
    echo "Also you can install kubectl via gcloud command."
    echo "$ gcloud components install kubectl"
    exit 1
  fi

  if ! type skaffold > /dev/null; then
    echo "Install skaffold command via gcloud."
    echo ""
    echo "    $ gcloud components install skaffold"
    echo ""
    echo "Otherwise, follow the official installation guide: https://skaffold.dev/docs/install/"
    exit 1
  fi
}

#######################################
# Check if the current kubectl context is correct
# Arguments:
#   None
#######################################
function set_kubectl_context() {
  gcloud container clusters get-credentials reporting-api-demo-cluster-vqwsdu1r

  context=$(kubectl config current-context)
  guess="gke_reporting-api-s2la0q_us-central1-b_reporting-api-demo-cluster-vqwsdu1r"
  if [[ "${context}" -ne "${guess}" ]]; then
    echo "Set kubectl context to ${guess}"
    kubectl config use-context "${guess}"
  fi
}

#######################################
# Set skaffold's default container repository
# Arguments:
#   None
#######################################
function set_skaffold_default_registry() {
  gcloud auth configure-docker
  skaffold config set default-repo us-central1-docker.pkg.dev/reporting-api-s2la0q/reporting-api-registry
}


set_kubectl_context
set_skaffold_default_registry