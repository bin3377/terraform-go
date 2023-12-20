package terraform

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type PlanConfig struct {
	ChDir   string
	Destroy bool
	Targets []string
	Vars    map[string]string
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
		args = append(args, "-var-file="+tempDir+"/vars.tfvar")
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

	args = append(args, "-input=false", "-json", "-out="+planOutputPath)
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

func createVarFile(vars map[string]string, dir string) error {
	f, err := os.Create(dir + "/vars.tfvar")
	if err != nil {
		return err
	}
	defer f.Close()
	for k, v := range vars {
		_, err := f.WriteString(k + "=" + v + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
