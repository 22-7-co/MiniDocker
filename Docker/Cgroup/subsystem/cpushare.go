package subsystem

type CpuShareSubsystem struct{}

func (c CpuShareSubsystem) Name() string {
	return "cpu"
}

func (c CpuShareSubsystem) Set(containerName string, res *ResourceConfig) error {
	return nil
}

func (c CpuShareSubsystem) Apply(containerName string, pid int) error {
	return nil
}

func (c CpuShareSubsystem) Remove(containerName string) error {
	return nil
}
