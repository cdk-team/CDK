/*
Copyright 2022 The Authors of https://github.com/CDK-TEAM/CDK .

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package evaluate

import "testing"

func TestEvaluatorProfiles(t *testing.T) {
	evaluator := NewEvaluator()
	profiles := evaluator.Profiles()
	required := map[string]bool{
		ProfileBasic:      false,
		ProfileExtended:   false,
		ProfileAdditional: false,
	}
	for _, profile := range profiles {
		if _, ok := required[profile.ID]; ok {
			required[profile.ID] = true
		}
	}
	for id, seen := range required {
		if !seen {
			t.Fatalf("expected profile %q to be registered", id)
		}
	}
}

func TestProfileComposition(t *testing.T) {
	evaluator := NewEvaluator()
	basic, ok := evaluator.Profile(ProfileBasic)
	if !ok {
		t.Fatalf("basic profile not registered")
	}
	extended, ok := evaluator.Profile(ProfileExtended)
	if !ok {
		t.Fatalf("extended profile not registered")
	}
	if len(extended.Categories) <= len(basic.Categories) {
		t.Fatalf("expected extended profile to contain more categories than basic profile")
	}
}

func TestRunUnknownProfile(t *testing.T) {
	evaluator := NewEvaluator()
	if err := evaluator.RunProfile("unknown", nil); err == nil {
		t.Fatalf("expected an error when running an unknown profile")
	}
}

func TestCallWrappers(t *testing.T) {
	CallBasics()
	CallAddedFunc()
	if err := NewEvaluator().RunProfile(ProfileExtended, nil); err != nil {
		t.Fatalf("extended profile run failed: %v", err)
	}
}
