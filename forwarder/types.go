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

package main

import "fmt"

// report is a type of single report sent by Reporting API
//
// TODO(yoshifumi): change the type of Body to concrete struct
type report struct {
	Typ       string `json:"type"`
	URL       string `json:"url"`
	UserAgent string `json:"user_agent"`
	Body      struct {
		// common to all types
		Disposition string

		// common to csp-violation, coep
		BlockedURL string

		// csp-violation
		DocumentURL        string
		EffectiveDirective string
		OriginalPolicy     string
		Referrer           string
		Sample             string
		StatusCode         int

		// common to deprecation, permissions-policy-violation, document-policy-violation, intervention
		ColumnNumber int
		LineNumber   int
		ID           string
		Message      string
		SourceFile   string

		// common to permissions-policy-violation, document-policy-violation
		PolicyID string

		// coep
		Destination string
	} `json:"body"`
	Age int `json:"age"`
}

// print is for debug use
func (r *report) print() {
	fmt.Printf(`
Type:	%v
URL:	%v
UA:	%v
Age:	%v
Body:	%v
	`, r.Typ, r.URL, r.UserAgent, r.Age, r.Body)
}
