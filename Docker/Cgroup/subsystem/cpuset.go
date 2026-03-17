package subsystem

type CpuSetSubsystem struct{}

func (c CpuSetSubsystem) Name() string {
	return "cpuset"
}

func (c CpuSetSubsystem) Set(containerName string, res *ResourceConfig) error {
	return nil
}

func (c CpuSetSubsystem) Apply(containerName string, pid int) error {
	return nil
}

func (c CpuSetSubsystem) Remove(containerName string) error {
	return nil
}
