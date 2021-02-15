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
	composego "github.com/compose-spec/compose-go/types"
)

// Manifest contains the tracked project's docker-compose sources and deployment environments
type Manifest struct {
	Id           string       `yaml:"id,omitempty" json:"id,omitempty"`
	Sources      *Sources     `yaml:"compose,omitempty" json:"compose,omitempty"`
	Environments Environments `yaml:"environments,omitempty" json:"environments,omitempty"`
	Skaffold     string       `yaml:"skaffold,omitempty" json:"skaffold,omitempty"`
}

// Sources tracks a project's docker-compose sources
type Sources struct {
	Files    []string `yaml:"-" json:"-"`
	override *composeOverride
}

// Environments tracks a project's deployment environments
type Environments []*Environment

// Environment is a deployment environment
type Environment struct {
	Name     string `yaml:"-" json:"-"`
	File     string `yaml:"-" json:"-"`
	override *composeOverride
}

// composeOverride augments a compose project with labels and env vars to produce
// k8s deployment config
type composeOverride struct {
	Version  string   `yaml:"version,omitempty" json:"version,omitempty" diff:"version"`
	Services Services `json:"services" diff:"services"`
	Volumes  Volumes  `yaml:",omitempty" json:"volumes,omitempty" diff:"volumes"`
}

// ComposeProject wrapper around a compose-go Project. It also provides the original
// compose file version.
type ComposeProject struct {
	version string
	*composego.Project
}

// ServiceConfig is a shallow version of a compose-go ServiceConfig
type ServiceConfig struct {
	Name        string                      `yaml:"-" json:"-" diff:"name"`
	Labels      composego.Labels            `yaml:",omitempty" json:"labels,omitempty" diff:"labels"`
	Environment composego.MappingWithEquals `yaml:",omitempty" json:"environment,omitempty" diff:"environment"`
	Extensions  map[string]interface{}      `yaml:",inline" json:"-"`
}

type secretHit struct {
	svcName     string
	envVar      string
	description string
}

// Services is a list of ServiceConfig
type Services []ServiceConfig

// Volumes is a mapping of volume name to VolumeConfig
type Volumes map[string]VolumeConfig

// VolumeConfig is a shallow version of a compose-go VolumeConfig
type VolumeConfig struct {
	Name   string           `yaml:",omitempty" json:"name,omitempty" diff:"name"`
	Labels composego.Labels `yaml:",omitempty" json:"labels,omitempty" diff:"labels"`
}

// changeset tracks changes made to a version, services and volumes
type changeset struct {
	version  change
	services []change
	volumes  []change
}

// change describes a create, update or delete modification
// targeting an attribute in a version, service or volume.
type change struct {
	Type   string
	Value  interface{}
	Parent string
	Target string
	Index  interface{}
}

// ErrorHandler is a callback function that handles error and returns error
type ErrorHandler func(error) error

// ChangeHandler is a callback function that handles change and returns error, e.g. change event when in dev mode
type ChangeHandler func(string) error

// RunFunc is a callback function expected to run before/after the current command
type RunFunc func() error

// DevOptions contains parameters required for Dev loop
type DevOptions struct {
	Skaffold      bool
	Namespace     string
	Kubecontext   string
	Kevenv        string
	Tail          bool
	ManualTrigger bool
	Verbose       bool
}
