package packit

import "os"

type BuildPlanProvision struct {
	Name string
}

type BuildPlanRequirement struct {
	Name     string
	Version  string
	Metadata interface{}
}

type BuildPlan struct {
	Provides []BuildPlanProvision
	Requires []BuildPlanRequirement
}

type DetectContext struct {
	WorkingDir string
}

type DetectResult struct {
	Plan BuildPlan
}

type DetectFunc func(DetectContext) (DetectResult, error)

func Detect(args []string, f DetectFunc) {
	dir, err := os.Getwd()
	if err != nil {
		//TODO: need to implement ErrorFunc in order to test this
		panic("blah")
	}

	f(DetectContext{
		WorkingDir: dir,
	})

}
