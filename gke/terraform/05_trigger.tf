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

locals {
  paths         = split("/", google_artifact_registry_repository.reporting_api.id)
  registry_path = format("%s-docker.pkg.dev/%s/%s", local.paths[3], local.paths[1], local.paths[5])
}

data "template_file" "init" {
  template = file("${path.module}/templates/init.sh.tpl")

  vars = {
    project_id    = google_project.demo_project.project_id
    cluster_name  = google_container_cluster.demo_cluster.name
    zone          = google_container_cluster.demo_cluster.location
    registry_path = local.registry_path
  }
}

resource "null_resource" "current_project" {
  provisioner "local-exec" {
    command = <<-EOT
      gcloud config set project ${google_project.demo_project.project_id}
      gcloud config set compute/zone ${google_container_cluster.demo_cluster.location}
    EOT
  }
}

resource "local_file" "generate_script" {
  depends_on = [
    google_container_cluster.demo_cluster,
  ]

  content  = data.template_file.init.rendered
  filename = "${path.module}/init.sh"
}

resource "null_resource" "init" {
  depends_on = [
    null_resource.current_project,
    local_file.generate_script,
  ]

  provisioner "local-exec" {
    command = "bash ${path.module}/init.sh"
  }
}
