// Copyright (c) OpenFaaS Project 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"errors"
	"fmt"
	"testing"

	"github.com/openfaas/faas-cli/stack"
)

func Test_build(t *testing.T) {

	aTests := [][]string{
		{"build"},
		{"build", "--image=my_image"},
		{"build", "--image=my_image", "--handler=/path/to/fn/"},
	}

	for _, aTest := range aTests {
		faasCmd.SetArgs(aTest)
		err := faasCmd.Execute()
		if err == nil {
			t.Fatalf("No error found while testing \n%v", err)
		}
	}
}

func Test_parseBuildArgs_ValidParts(t *testing.T) {
	mapped, err := parseBuildArgs([]string{"k=v"})

	if err != nil {
		t.Errorf("err was supposed to be nil but was: %s", err.Error())
		t.Fail()
	}

	if mapped["k"] != "v" {
		t.Errorf("value for 'k', want: %s got: %s", "v", mapped["k"])
		t.Fail()
	}
}

func Test_parseBuildArgs_NoSeparator(t *testing.T) {
	_, err := parseBuildArgs([]string{"kv"})

	want := "each build-arg must take the form key=value"
	if err != nil && err.Error() != want {
		t.Errorf("Expected an error due to missing seperator")
		t.Fail()
	}
}

func Test_parseBuildArgs_EmptyKey(t *testing.T) {
	_, err := parseBuildArgs([]string{"=v"})

	want := "build-arg must have a non-empty key"
	if err == nil {
		t.Errorf("Expected an error due to missing key")
		t.Fail()
	} else if err.Error() != want {
		t.Errorf("missing key error want: %s, got: %s", want, err.Error())
		t.Fail()
	}
}

func Test_parseBuildArgs_MultipleSeparators(t *testing.T) {
	mapped, err := parseBuildArgs([]string{"k=v=z"})

	if err != nil {
		t.Errorf("Expected second separator to be included in value")
		t.Fail()
	}

	if mapped["k"] != "v=z" {
		t.Errorf("value for 'k', want: %s got: %s", "v=z", mapped["k"])
		t.Fail()
	}
}

func Test_build_NoLanguage(t *testing.T) {
	functions := map[string]stack.Function{
		"first": {
			Handler:  "first_handler",
			Image:    "first_image",
			Language: "",
		},
		"second": {
			Handler:  "second_handler",
			Image:    "second_image",
			Language: "",
		},
	}
	services := &stack.Services{
		Functions: functions,
	}

	mockbuildImageFunc := func(image string, handler string, functionName string, language string, nocache bool, squash bool, shrinkwrap bool, buildArgMap map[string]string, buildOptions []string, tag string) error {
		return nil
	}

	err := build(services, 2, true, mockbuildImageFunc)
	if err == nil {
		fmt.Errorf("Expected error")
	}
}

func Test_build_SkipBuilds(t *testing.T) {
	functions := map[string]stack.Function{
		"first": {
			Handler:   "first_handler",
			Image:     "first_image",
			Language:  "first_language",
			SkipBuild: true,
		},
		"second": {
			Handler:   "second_handler",
			Image:     "second_image",
			Language:  "second_language",
			SkipBuild: true,
		},
	}
	services := &stack.Services{
		Functions: functions,
	}

	mockbuildImageFunc := func(image string, handler string, functionName string, language string, nocache bool, squash bool, shrinkwrap bool, buildArgMap map[string]string, buildOptions []string, tag string) error {
		return nil
	}
	err := build(services, 2, true, mockbuildImageFunc)
	if err != nil {
		fmt.Errorf("Did not expected error")
	}
}

func Test_build_BuildImageNoErrors(t *testing.T) {
	functions := map[string]stack.Function{
		"first": {
			Handler:  "first_handler",
			Image:    "first_image",
			Language: "first_language",
		},
		"second": {
			Handler:  "second_handler",
			Image:    "second_image",
			Language: "second_language",
		},
	}
	services := &stack.Services{
		Functions: functions,
	}

	mockbuildImageFunc := func(image string, handler string, functionName string, language string, nocache bool, squash bool, shrinkwrap bool, buildArgMap map[string]string, buildOptions []string, tag string) error {
		return nil
	}
	err := build(services, 2, true, mockbuildImageFunc)
	if err != nil {
		fmt.Errorf("Did not expected error")
	}
}

func Test_build_BuildImageOneError(t *testing.T) {
	functions := map[string]stack.Function{
		"first": {
			Handler:  "first_handler",
			Image:    "first_image",
			Language: "first_language",
		},
		"second": {
			Handler:  "second_handler",
			Image:    "second_image",
			Language: "second_language",
		},
	}
	services := &stack.Services{
		Functions: functions,
	}

	mockbuildImageFunc := func(image string, handler string, functionName string, language string, nocache bool, squash bool, shrinkwrap bool, buildArgMap map[string]string, buildOptions []string, tag string) error {
		if functionName == "first" {
			return nil
		}
		return errors.New("Build error")
	}
	err := build(services, 2, true, mockbuildImageFunc)
	if err == nil {
		fmt.Errorf("Expected error")
	}
}

func Test_build_BuildImageAllErrors(t *testing.T) {
	functions := map[string]stack.Function{
		"first": {
			Handler:  "first_handler",
			Image:    "first_image",
			Language: "first_language",
		},
		"second": {
			Handler:  "second_handler",
			Image:    "second_image",
			Language: "second_language",
		},
	}
	services := &stack.Services{
		Functions: functions,
	}

	mockbuildImageFunc := func(image string, handler string, functionName string, language string, nocache bool, squash bool, shrinkwrap bool, buildArgMap map[string]string, buildOptions []string, tag string) error {
		return errors.New("Build error")
	}
	err := build(services, 2, true, mockbuildImageFunc)
	if err == nil {
		fmt.Errorf("Expected error")
	}
}

func Test_build_BuildImagePanics(t *testing.T) {
	functions := map[string]stack.Function{
		"first": {
			Handler:  "first_handler",
			Image:    "first_image",
			Language: "first_language",
		},
		"second": {
			Handler:  "second_handler",
			Image:    "second_image",
			Language: "second_language",
		},
	}
	services := &stack.Services{
		Functions: functions,
	}
	mockbuildImageFunc := func(image string, handler string, functionName string, language string, nocache bool, squash bool, shrinkwrap bool, buildArgMap map[string]string, buildOptions []string, tag string) error {
		panic("build error")
	}
	err := build(services, 2, false, mockbuildImageFunc)
	if err == nil {
		fmt.Errorf("Expected error")
	}

}
