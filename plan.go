package terraform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/sirupsen/logrus"
)

type PlanConfig struct {
	ChDir   string
	Destroy bool
	Targets []string
	Vars    map[string]any
}

type PlanResult struct {
	Plan       Plan
	PlanBinary []byte
	PlanPath   string
	Logs       []string
}

var planOutputPath string

func (c *PlanConfig) getArgs() ([]string, error) {
	planOutputPath = tempDir + "/output.tfplan"
	args := []string{"plan"}
	if len(c.Vars) > 0 {
		if err := createVarFile(c.Vars, tempDir); err != nil {
			return nil, err
		}
		args = append(args, "-var-file="+tempDir+"/vars.tfvars")
	}

	if c.ChDir != "" {
		args = append([]string{"-chdir=" + c.ChDir}, args...)
	}

	if c.Destroy {
		args = append(args, "-destroy")
	}

	for _, t := range c.Targets {
		args = append(args, "-target="+t)
	}

	args = append(args, "-input=false", "-no-color", "-out="+planOutputPath)
	return args, nil
}

func (c *PlanConfig) getShowArgs() []string {
	args := []string{"show", "-json", planOutputPath}

	if c.ChDir != "" {
		args = append([]string{"-chdir=" + c.ChDir}, args...)
	}

	return args
}

func (c *PlanConfig) Plan() (*PlanResult, error) {
	args, err := c.getArgs()
	if err != nil {
		return nil, err
	}
	out, err := runTFWithEnv(args)
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(planOutputPath)
	if err != nil {
		logrus.Errorf("fail to read %s", planOutputPath)
		return nil, err
	}

	result := PlanResult{
		PlanBinary: bytes,
		PlanPath:   planOutputPath,
		Logs:       strings.Split(out, "\n"),
	}

	out, err = runTFWithEnv(c.getShowArgs())
	if err != nil {
		return nil, err
	}

	var plan Plan
	if err := json.Unmarshal([]byte(out), &plan); err != nil {
		return nil, err
	}
	result.Plan = plan
	return &result, nil
}

func createVarFile(vars map[string]any, dir string) error {
	f, err := os.Create(dir + "/vars.tfvars")
	if err != nil {
		return err
	}
	defer f.Close()
	bytes, err := toHCL(vars)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}

func toHCL(in any) ([]byte, error) {
	json, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	f, err := hcl.ParseBytes(json)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	err = printer.DefaultConfig.Fprint(buf, f.Node)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(buf)
	return stripeKeyQuotes(buf.Bytes()), nil
}

func stripeKeyQuotes(content []byte) []byte {
	regex := regexp.MustCompile(`(?m)^"(.*)" = (.*)$`)
	result := regex.ReplaceAll(content, []byte("$1 = $2\n"))
	return result
}
