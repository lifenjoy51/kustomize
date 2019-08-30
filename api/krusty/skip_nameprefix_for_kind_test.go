// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package krusty_test

import (
	"testing"
)

// TestNameTransformation_skippingSecrets validates that
// PrefixSuffixTransformer can skip based on resource kind, and that CustomResourceDefinitions are always skipped.
// TODO: also the nameref stuff?
func TestNameTransformation_skippingSecrets(t *testing.T) {
	th := makeTestHarness(t)

	th.WriteK("/nameandconfig", `
namePrefix: p1-
nameSuffix: -s1
resources:
- resources.yaml
configurations:
- skip-secrets-prefix.yaml
`)

	th.WriteF("/nameandconfig/resources.yaml", `
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm1
---
apiVersion: v1
kind: Secret
metadata:
  name: secret
---
kind: CustomResourceDefinition
metadata:
  name: my-crd
`)

	th.WriteF("/nameandconfig/skip-secrets-prefix.yaml", `
namePrefix:
- path: metadata/name
  apiVersion: v1
  kind: Secret
  skip: true
`)

	m := th.Run("/nameandconfig", th.MakeDefaultOptions())
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: p1-cm1-s1
---
apiVersion: v1
kind: Secret
metadata:
  name: secret
---
kind: CustomResourceDefinition
metadata:
  name: my-crd
`)
}
