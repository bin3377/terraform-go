package terraform

import (
	"github.com/sirupsen/logrus"
)

type InitConfig struct {
	ChDir         string
	Upgrade       bool
	Reconfigure   bool
	MigrateState  bool
	BackendConfig map[string]string
}

func (c *InitConfig) getArgs() []string {
	args := []string{"init"}
	if c.ChDir != "" {
		args = append([]string{"-chdir=" + c.ChDir}, args...)
	}
	if c.Upgrade {
		args = append(args, "-upgrade")
	}
	if c.Reconfigure {
		args = append(args, "-reconfigure")
	}
	if c.MigrateState {
		args = append(args, "-migrate-state")
	}
	if len(c.BackendConfig) > 0 {
		for k, v := range c.BackendConfig {
			args = append(args, "-backend-config="+k+"="+v)
		}
	}
	return args
}

func (c *InitConfig) Init() error {
	args := c.getArgs()
	out, err := runTFWithEnv(args)
	if err != nil {
		return err
	}
	logrus.Debug(out)
	return nil
}
