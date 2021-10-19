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

terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 3.88"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 3.88"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 3.1"
    }

    local = {
      source  = "hashicorp/local"
      version = "~> 2.1"
    }
  }
}

data "google_billing_account" "account" {
  count        = var.billing_account != null ? 1 : 0
  display_name = var.billing_account
}

locals {
  region = join("-", slice(split("-", var.zone), 0, 2))
}

resource "google_project_service" "iam" {
  project                    = google_project.demo_project.project_id
  service                    = "iam.googleapis.com"
  disable_dependent_services = true
}

resource "google_project_service" "gce" {
  project = google_project.demo_project.project_id
  service = "compute.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
}

resource "google_project_service" "gke" {
  project = google_project.demo_project.project_id
  service = "container.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
}

resource "google_project_service" "artifact_registry" {
  project = google_project.demo_project.project_id
  service = "artifactregistry.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
}

resource "google_project_service" "logging" {
  project = google_project.demo_project.project_id
  service = "logging.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
}

resource "google_project_service" "monitoring" {
  project = google_project.demo_project.project_id
  service = "monitoring.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
}
