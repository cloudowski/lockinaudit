// Package scanner wires the directory walk and the source detectors into a
// single fact-collection step, shared by `lockinaudit inspect` and developer tools
// such as tools/calibrate.
package scanner

import (
	"context"

	"github.com/cloudowski/lockinaudit/internal/model"
	"github.com/cloudowski/lockinaudit/internal/source"
	"github.com/cloudowski/lockinaudit/internal/source/cdk"
	"github.com/cloudowski/lockinaudit/internal/source/cicd"
	"github.com/cloudowski/lockinaudit/internal/source/cloudformation"
	"github.com/cloudowski/lockinaudit/internal/source/code"
	"github.com/cloudowski/lockinaudit/internal/source/dockerfile"
	"github.com/cloudowski/lockinaudit/internal/source/kubernetes"
	"github.com/cloudowski/lockinaudit/internal/source/packages"
	"github.com/cloudowski/lockinaudit/internal/source/pulumi"
	"github.com/cloudowski/lockinaudit/internal/source/terraform"
	"github.com/cloudowski/lockinaudit/internal/walk"
)

// Detectors returns the first-stage and dependent detectors.
func Detectors() (first, dependent []source.Detector) {
	first = []source.Detector{
		terraform.New(),
		cloudformation.New(),
		kubernetes.New(),
		kubernetes.NewHelm(),
		kubernetes.NewKustomize(),
		dockerfile.New(),
		cicd.New(),
		packages.New(),
		pulumi.NewYAML(),
	}
	dependent = []source.Detector{code.New(), cdk.New(), pulumi.New()}
	return first, dependent
}

// Collect walks opts.Root (honoring the extra ignore globs) and runs all
// detectors over it. The only error is a failed walk.
func Collect(ctx context.Context, ignores []string, opts source.Options) (model.Facts, error) {
	idx, err := walk.Walk(opts.Root, ignores)
	if err != nil {
		return model.Facts{}, err
	}
	first, dep := Detectors()
	return source.RunAll(ctx, first, dep, idx, opts), nil
}
