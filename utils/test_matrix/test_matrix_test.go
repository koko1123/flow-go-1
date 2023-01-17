package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Can't have a const []string so resorting to using a test helper function.
func getAllFlowPackages() []string {
	return []string{
		flowPackagePrefix + "abc",
		flowPackagePrefix + "abc/def",
		flowPackagePrefix + "abc/def/ghi",
		flowPackagePrefix + "def",
		flowPackagePrefix + "def/abc",
		flowPackagePrefix + "ghi",
		flowPackagePrefix + "jkl",
		flowPackagePrefix + "mno/abc",
		flowPackagePrefix + "pqr",
		flowPackagePrefix + "stu",
		flowPackagePrefix + "vwx",
		flowPackagePrefix + "vwx/ghi",
		flowPackagePrefix + "yz",
	}
}

func TestListTargetPackages(t *testing.T) {
	targetPackages, seenPackages := listTargetPackages([]string{"abc", "ghi"}, getAllFlowPackages())
	require.Equal(t, 2, len(targetPackages))
	require.Equal(t, 4, len(seenPackages))

	// there should be 3 packages that start with "abc"
	require.Equal(t, 3, len(targetPackages["abc"]))
	require.Contains(t, targetPackages["abc"], flowPackagePrefix+"abc")
	require.Contains(t, targetPackages["abc"], flowPackagePrefix+"abc/def")
	require.Contains(t, targetPackages["abc"], flowPackagePrefix+"abc/def/ghi")

	// there should be 1 package that starts with "ghi"
	require.Equal(t, 1, len(targetPackages["ghi"]))
	require.Contains(t, targetPackages["ghi"], flowPackagePrefix+"ghi")

	require.Contains(t, seenPackages, flowPackagePrefix+"abc")
	require.Contains(t, seenPackages, flowPackagePrefix+"abc/def")
	require.Contains(t, seenPackages, flowPackagePrefix+"abc/def/ghi")
	require.Contains(t, seenPackages, flowPackagePrefix+"ghi")
}

func TestListOtherPackages(t *testing.T) {
	var seenPackages = make(map[string]string)
	seenPackages[flowPackagePrefix+"abc"] = flowPackagePrefix + "abc"
	seenPackages[flowPackagePrefix+"ghi"] = flowPackagePrefix + "ghi"
	seenPackages[flowPackagePrefix+"mno/abc"] = flowPackagePrefix + "mno/abc"
	seenPackages[flowPackagePrefix+"stu"] = flowPackagePrefix + "stu"

	otherPackages := listOtherPackages(getAllFlowPackages(), seenPackages)

	require.Equal(t, 9, len(otherPackages))

	require.Contains(t, otherPackages, flowPackagePrefix+"abc/def")
	require.Contains(t, otherPackages, flowPackagePrefix+"abc/def/ghi")
	require.Contains(t, otherPackages, flowPackagePrefix+"def")
	require.Contains(t, otherPackages, flowPackagePrefix+"def/abc")
	require.Contains(t, otherPackages, flowPackagePrefix+"jkl")
	require.Contains(t, otherPackages, flowPackagePrefix+"pqr")
	require.Contains(t, otherPackages, flowPackagePrefix+"vwx")
	require.Contains(t, otherPackages, flowPackagePrefix+"vwx/ghi")
	require.Contains(t, otherPackages, flowPackagePrefix+"yz")
}

func TestGenerateTestMatrix(t *testing.T) {
	targetPackages, seenPackages := listTargetPackages([]string{"abc", "ghi"}, getAllFlowPackages())
	require.Equal(t, 2, len(targetPackages))
	require.Equal(t, 4, len(seenPackages))

	otherPackages := listOtherPackages(getAllFlowPackages(), seenPackages)

	matrix := generateTestMatrix(targetPackages, otherPackages)

	// should be 3 groups in test matrix: abc, ghi, others
	require.Equal(t, 3, len(matrix))

	require.Contains(t, matrix, testMatrix{
		Name:     "abc",
		Packages: "github.com/koko1123/flow-go-1/abc github.com/koko1123/flow-go-1/abc/def github.com/koko1123/flow-go-1/abc/def/ghi"},
	)
	require.Contains(t, matrix, testMatrix{
		Name:     "ghi",
		Packages: "github.com/koko1123/flow-go-1/ghi"},
	)
	require.Contains(t, matrix, testMatrix{
		Name:     "others",
		Packages: "github.com/koko1123/flow-go-1/def github.com/koko1123/flow-go-1/def/abc github.com/koko1123/flow-go-1/jkl github.com/koko1123/flow-go-1/mno/abc github.com/koko1123/flow-go-1/pqr github.com/koko1123/flow-go-1/stu github.com/koko1123/flow-go-1/vwx github.com/koko1123/flow-go-1/vwx/ghi github.com/koko1123/flow-go-1/yz"},
	)
}
