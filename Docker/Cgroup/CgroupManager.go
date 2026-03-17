package cgroup

import "mini-docker/Docker/Cgroup/subsystem"

type CgroupManager struct {
	CgroupName string
	Resouce    *subsystem.ResourceConfig
}

func NewCgroupManager(cgroupName string) *CgroupManager {
	return &CgroupManager{
		CgroupName: cgroupName,
	}
}

func (c *CgroupManager) Apply(pid int) error {
	for _, ins := range subsystem.Subsystems {
		if err := ins.Apply(c.CgroupName, pid); err != nil {
			return err
		}
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystem.ResourceConfig) error {
	for _, ins := range subsystem.Subsystems {
		if err := ins.Set(c.CgroupName, res); err != nil {
			return err
		}
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _, ins := range subsystem.Subsystems {
		if err := ins.Remove(c.CgroupName); err != nil {
			return err
		}
	}
	return nil
}
