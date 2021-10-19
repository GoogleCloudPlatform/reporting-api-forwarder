// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

resource "google_service_account" "default" {
  project      = google_project.demo_project.project_id
  account_id   = "reporting-api-demo"
  display_name = "Reporting API demo service account"
}

resource "google_artifact_registry_repository" "reporting_api" {
  provider = google-beta
  project  = google_project.demo_project.project_id
  depends_on = [
    google_project_service.artifact_registry,
  ]
  location      = local.region
  repository_id = "reporting-api-registry"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository_iam_binding" "binding" {
  provider = google-beta
  depends_on = [
    google_project_service.artifact_registry,
  ]
  project    = google_artifact_registry_repository.reporting_api.project
  location   = google_artifact_registry_repository.reporting_api.location
  repository = google_artifact_registry_repository.reporting_api.name
  role       = "roles/artifactregistry.repoAdmin"
  members = [
    "serviceAccount:${google_service_account.default.email}",
  ]
}

resource "random_string" "cluster_suffix" {
  special = false
  upper   = false
  lower   = true
  number  = true
  length  = 8
}


resource "google_container_cluster" "demo_cluster" {
  name     = "reporting-api-demo-cluster-${random_string.cluster_suffix.result}"
  project  = google_project.demo_project.project_id
  location = var.zone

  depends_on = [
    google_project_service.gce,
    google_project_service.gke,
    google_service_account.default,
    google_artifact_registry_repository.reporting_api,
    google_artifact_registry_repository_iam_binding.binding,
  ]

  remove_default_node_pool = true
  initial_node_count       = 1

  release_channel {
    channel = "STABLE"
  }

  timeouts {
    create = "10m"
    update = "20m"
  }

}

resource "google_container_node_pool" "demo_cluster_preemptible_nodes" {
  name       = "demo-cluster-np-${random_string.cluster_suffix.result}"
  project    = google_project.demo_project.project_id
  location   = var.zone
  cluster    = google_container_cluster.demo_cluster.name
  node_count = 1

  node_config {
    machine_type    = "e2-standard-2"
    preemptible     = true
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }

  autoscaling {
    min_node_count = 2
    max_node_count = 4
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }
}
