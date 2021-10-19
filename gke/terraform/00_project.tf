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

resource "random_string" "project_suffix" {
  special = false
  upper   = false
  lower   = true
  number  = true
  length  = 6
}

data "google_billing_account" "default" {
  display_name = var.billing_account
}

locals {
  project_id = "reporting-api-${random_string.project_suffix.result}"
}

resource "google_project" "demo_project" {
  project_id      = local.project_id
  name            = local.project_id
  billing_account = data.google_billing_account.default.id
}
