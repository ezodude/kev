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

package interpolate

import (
	"bytes"
	"regexp"

	"github.com/tidwall/gjson"
)

const jsonPathKeyPattern = `\$\{(.*\.+.*)\}`

func NewJsonPathResolver() Resolver {
	return jsonPathResolver{
		resolvedKeys: map[string]gjson.Result{},
	}
}

type jsonPathResolver struct {
	resolvedKeys map[string]gjson.Result
}

func (r jsonPathResolver) Resolve(source, target []byte) ([]byte, error) {
	err := r.resolveKeys(source, target)
	if err != nil {
		return nil, err
	}
	return r.resolve(target), nil
}

func (r jsonPathResolver) resolveKeys(source, target []byte) error {
	pattern, err := regexp.Compile(jsonPathKeyPattern)
	if err != nil {
		return err
	}

	found := pattern.FindAllSubmatch(target, -1)
	for _, f := range found {
		key, jsonPath := string(f[0]), string(f[1])
		r.resolvedKeys[key] = gjson.GetBytes(source, jsonPath)
	}
	return nil
}

func (r jsonPathResolver) resolve(target []byte) []byte {
	var resolved = target
	for k, v := range r.resolvedKeys {
		resolved = bytes.ReplaceAll(resolved, []byte(k), []byte(v.String()))
	}
	return resolved
}