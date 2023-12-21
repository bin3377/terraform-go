package terraform

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func assert(tb testing.TB, condition bool, msg string, v ...any) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		logrus.Printf("%s:%d: "+msg+"\n\n", append([]any{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		logrus.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func equals(tb testing.TB, exp, act any) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		logrus.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_Plan(t *testing.T) {
	testPath := "./testData/vpcs"
	t.Cleanup(func() {
		cleanup(testPath)
	})
	err := (&InitConfig{
		ChDir: testPath,
	}).Init()
	ok(t, err)

	r, err := (&PlanConfig{
		ChDir: testPath,
	}).Plan()

	ok(t, err)
	equals(t, r.Plan.PlannedValues.Outputs["foo"].Type, []any{"list", "string"})
	for _, v := range r.Plan.PlannedValues.Outputs["foo"].Value.([]any) {
		assert(t, strings.HasPrefix(v.(string), "vpc"), "Expected value to start with vpc")
	}
	logrus.Debugf("%v", r.Plan.PlannedValues.Outputs)
}

func Test_PlanWithVars(t *testing.T) {
	testPath := "./testData/input"
	t.Cleanup(func() {
		cleanup(testPath)
	})
	err := (&InitConfig{
		ChDir: testPath,
	}).Init()
	ok(t, err)

	r, err := (&PlanConfig{
		ChDir: testPath,
		Vars: map[string]any{
			"foo_in":      "bar",
			"foo_in_list": []string{"test", "me"},
			"foo_in_map": map[string]string{
				"test": "me",
			},
		},
	}).Plan()

	ok(t, err)
	equals(t, r.Plan.PlannedValues.Outputs["foo_out"].Type, "string")
	equals(t, r.Plan.PlannedValues.Outputs["foo_out"].Value, "bar")
	equals(t, r.Plan.PlannedValues.Outputs["foo_out_list"].Type, []any{"list", "string"})
	equals(t, r.Plan.PlannedValues.Outputs["foo_out_list"].Value, []any{"test", "me"})
	equals(t, r.Plan.PlannedValues.Outputs["foo_out_map"].Type, []any{"map", "string"})
	equals(t, r.Plan.PlannedValues.Outputs["foo_out_map"].Value, map[string]any{"test": "me"})
	logrus.Debugf("%v", r.Plan.PlannedValues.Outputs)
}

func cleanup(path string) {
	patterns := []string{
		"*.tfplan",
		"*.tfstate",
		"*.tfstate.backup",
		"*.tfstate.migrate",
		".terraform",
		".terraform.lock.hcl",
	}
	for _, pattern := range patterns {
		files, err := filepath.Glob(filepath.Join(path, pattern))
		if err != nil {
			logrus.Errorf("Error cleaning up %s: %s", pattern, err.Error())
			continue
		}
		for _, file := range files {
			logrus.Debugf("Removing %s", file)
			os.RemoveAll(file)
		}
	}
}
