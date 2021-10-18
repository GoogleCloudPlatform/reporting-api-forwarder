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

variable "billing_account" {
  description = "Google Cloud Billing Account for the demo project"
  type        = string
}

variable "project_id" {
  description = "Google Cloud Porject ID to use for demo"
  type        = string
}

variable "zone" {
  description = "Google Cloud zone for GKE and Artifact Registry. https://cloud.google.com/compute/docs/regions-zones"
  type        = string
}

variable "custom_domain" {
  description = "Custom domain for the reserved static external IP address"
  type        = string
}
