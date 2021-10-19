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

data "template_file" "ingress_yaml" {
  template = file("${path.module}/templates/ingress.yaml.tpl")

  vars = {
    static_ip_name = google_compute_global_address.static.name,
  }
}

resource "local_file" "ingress_yaml" {
  depends_on = [
    google_compute_global_address.static,
    data.template_file.ingress_yaml,
  ]

  content  = data.template_file.ingress_yaml.rendered
  filename = "${path.module}/../manifests/ingress.yaml"
}

data "template_file" "certificate_yaml" {
  template = file("${path.module}/templates/certificate.yaml.tpl")

  vars = {
    custom_domain = var.custom_domain,
  }
}

resource "local_file" "certificate_yaml" {
  depends_on = [
    data.template_file.certificate_yaml,
  ]

  content  = data.template_file.certificate_yaml.rendered
  filename = "${path.module}/../manifests/certificate.yaml"
}
