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
	"io/ioutil"
	"os"

	"github.com/appvia/kev/pkg/kev"
	"github.com/appvia/kev/pkg/kev/converter"
	"github.com/appvia/kev/pkg/kev/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("wordpress-mysql", func() {
	var (
		workingDir   string
		sources      []string
		envs         []string
		manifestsDir string
		format       string
	)

	BeforeEach(func() {
		workingDir = "./testdata/wordpress-mysql"
		format = "kubernetes"
		manifestsDir = mustCreateTempDir()
		mustGenerateManifests(workingDir, manifestsDir, format, sources, envs)
	})

	AfterEach(func() {
		err := os.RemoveAll(manifestsDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Deployments", func() {
		It("creates a valid wordpress deployment", func() {
			Expect(true).ToNot(BeTrue())
		})
	})
})

func mustGenerateManifests(workingDir, manifestsDir, format string, sources, envs []string) {
	manifest, err := kev.Init(sources, envs, workingDir)
	if err != nil {
		log.Errorf("kev.Init for [%s] has failed", workingDir)
		panic(err)
	}

	if _, err := manifest.ReconcileConfig(); err != nil {
		log.Errorf("kev.Reconcile for [%s] has failed", workingDir)
		panic(err)
	}

	if _, err := manifest.RenderWithConvertor(converter.Factory(format), manifestsDir, false, envs); err != nil {
		log.Errorf("kev.RenderWithConvertor for [%s] has failed", workingDir)
		panic(err)
	}
}

func mustCreateTempDir() string {
	dirPath, err := ioutil.TempDir(".", "k8s-manifests")
	if err != nil {
		panic(err)
	}
	return dirPath
}
