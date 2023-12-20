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
	t.Cleanup(func() {
		os.RemoveAll("./testData/.terraform")
		os.RemoveAll("./testData/.terraform.lock.hcl")
		os.RemoveAll("./testData/terraform.tfstate")
		os.RemoveAll("./testData/tfplan")
	})
	err := (&InitConfig{
		ChDir: "./testData",
	}).Init()
	ok(t, err)

	r, err := (&PlanConfig{
		ChDir: "./testData",
	}).Plan()

	ok(t, err)
	equals(t, r.Plan.PlannedValues.Outputs["foo"].Type, []any{"list", "string"})
	for _, v := range r.Plan.PlannedValues.Outputs["foo"].Value.([]any) {
		assert(t, strings.HasPrefix(v.(string), "vpc"), "Expected value to start with vpc")
	}
	logrus.Debugf("%v", r.Plan.PlannedValues.Outputs)
}
