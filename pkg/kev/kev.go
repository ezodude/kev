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

package kev

import (
	"github.com/appvia/kev/pkg/kev/config"
	"github.com/appvia/kev/pkg/kev/converter"
	"github.com/appvia/kev/pkg/kev/log"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

const (
	// ManifestName main application manifest
	ManifestName = "kev.yaml"
	defaultEnv   = "dev"
)

// Init initialises a kev manifest including source compose files and environments.
// A default environment will be allocated if no environments were provided.
func Init(composeSources, envs []string, workingDir string) (*Manifest, error) {
	m, err := NewManifest(composeSources, workingDir)
	if err != nil {
		return nil, err
	}

	if _, err := m.CalculateSourcesBaseOverride(); err != nil {
		return nil, err
	}

	return m.MintEnvironments(envs), nil
}

// PrepareForSkaffold initialises a skaffold manifest for kev project.
func PrepareForSkaffold(envs []string) (*SkaffoldManifest, error) {
	s, err := NewSkaffoldManifest(envs)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Reconcile reconciles changes with docker-compose sources against deployment environments.
func Reconcile(workingDir string) (*Manifest, error) {
	m, err := LoadManifest(workingDir)
	if err != nil {
		return nil, err
	}
	if _, err := m.ReconcileConfig(); err != nil {
		return nil, errors.Wrap(err, "Could not reconcile project latest")
	}
	return m, err
}

// DetectSecrets detects any potential secrets defined in environment variables
// found either in sources or override environments.
// Any detected secrets are logged using a warning log level.
func DetectSecrets(workingDir string) error {
	m, err := LoadManifest(workingDir)
	if err != nil {
		return err
	}
	m.DetectSecretsInSources(config.SecretMatchers)
	m.DetectSecretsInEnvs(config.SecretMatchers)
	return nil
}

// Render renders k8s manifests for a kev app. It returns an app definition with rendered manifest info
func Render(workingDir string, format string, singleFile bool, dir string, envs []string) error {
	manifest, err := LoadManifest(workingDir)
	if err != nil {
		log.Error("Unable to load app manifest")
		return err
	}

	_, err = manifest.RenderWithConvertor(converter.Factory(format), dir, singleFile, envs)
	return err
}

// Watch continuously watches source compose files & environment overrides and notifies changes to a channel
func Watch(workDir string, envs []string, change chan<- string) error {
	manifest, err := LoadManifest(workDir)
	if err != nil {
		log.Error("Unable to load app manifest")
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					change <- event.Name
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Error(err)
			}
		}
	}()

	files := manifest.GetSourcesFiles()
	filteredEnvs, err := manifest.GetEnvironments(envs)
	for _, e := range filteredEnvs {
		files = append(files, e.File)
	}

	for _, f := range files {
		err = watcher.Add(f)
		if err != nil {
			return err
		}
	}

	<-done

	return nil
}
