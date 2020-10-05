// +build e2e

/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package e2e_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/kind/pkg/cluster"
)

var (
	clusterName    string
	kubeConfigPath string
	provider       *cluster.Provider
	err            error
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

var _ = BeforeSuite(func() {
	kubeConfigPath = "./testdata/kubeConfigPath"
	clusterName = "kind"

	provider, err = createCluster(clusterName, kubeConfigPath)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err = provider.Delete(clusterName, kubeConfigPath)
	Expect(err).NotTo(HaveOccurred())
})

func createCluster(name, kubeConfigPath string) (*cluster.Provider, error) {
	provider := cluster.NewProvider()
	err := provider.Create(name, cluster.CreateWithKubeconfigPath(kubeConfigPath))
	return provider, err
}
