package verifier

import (
	"github.com/bmatcuk/doublestar"
	"github.com/joshdk/go-junit"
	"github.com/kadaan/junit-tools/config"
	"github.com/kadaan/junit-tools/lib/command"
	"github.com/kadaan/junit-tools/lib/errors"
	"k8s.io/klog/v2"
	"os"
	"path/filepath"
)

func NewVerifier() command.Task[config.VerifyConfig] {
	return &verifier{}
}

type verifier struct {
}

func (v *verifier) Run(_ *config.VerifyConfig, args []string) error {
	var totals = &junit.Totals{}
	var errs []error
	workingDirectory := v.getWorkingDirectory()
	for _, arg := range args {
		matches, err := doublestar.Glob(arg)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, match := range matches {
			suites, err := junit.IngestFile(match)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			fileTotals := v.aggregate(suites)
			v.log(1, v.makeRelativePath(workingDirectory, match), fileTotals)
			totals = v.add(totals, fileTotals)
		}
	}
	v.log(0, "JUnit Results", totals)
	if totals.Failed+totals.Error > 0 {
		return errors.NewCommandError("One or more JUnit tests did not succeed")
	}
	return nil
}

func (v *verifier) log(level klog.Level, name string, totals *junit.Totals) {
	klog.V(level).Infof("%s: { Tests: %d, Passed: %d, Skipped: %d, Errored: %d, Failed: %d }",
		name, totals.Tests, totals.Passed, totals.Skipped, totals.Error, totals.Failed)
}

func (v *verifier) getWorkingDirectory() *string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil
	}
	return &workingDirectory
}

func (v *verifier) add(totals1 *junit.Totals, totals2 *junit.Totals) *junit.Totals {
	var result *junit.Totals
	if totals1 == nil {
		result = &junit.Totals{
			Tests:    0,
			Passed:   0,
			Skipped:  0,
			Failed:   0,
			Error:    0,
			Duration: 0,
		}
	} else {
		result = totals1
	}
	if totals2 != nil {
		result.Tests += totals2.Tests
		result.Passed += totals2.Passed
		result.Skipped += totals2.Skipped
		result.Failed += totals2.Failed
		result.Error += totals2.Error
		result.Duration += totals2.Duration
	}
	return result
}

func (v *verifier) aggregate(suites []junit.Suite) *junit.Totals {
	result := &junit.Totals{
		Tests:    0,
		Passed:   0,
		Skipped:  0,
		Failed:   0,
		Error:    0,
		Duration: 0,
	}
	if suites != nil {
		for _, suite := range suites {
			result.Tests += suite.Totals.Tests
			result.Duration += suite.Totals.Duration
			result.Passed += suite.Totals.Passed
			result.Skipped += suite.Totals.Skipped
			result.Failed += suite.Totals.Failed
			result.Error += suite.Totals.Error
		}
	}
	return result
}

func (v *verifier) makeRelativePath(directory *string, path string) string {
	if directory == nil {
		return path
	}
	if relativePath, err := filepath.Rel(*directory, path); err != nil {
		return path
	} else {
		return relativePath
	}
}
