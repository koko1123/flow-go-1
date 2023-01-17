package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

const header = `// AUTO-GENERATED BY %s.  DO NOT MODIFY.

package environment

import (
    "github.com/koko1123/flow-go-1/fvm/errors"
    "github.com/koko1123/flow-go-1/fvm/state"
	"github.com/koko1123/flow-go-1/module/trace"
)

func parseRestricted(
    txnState *state.TransactionState,
    spanName trace.SpanName,
) error {
    if txnState.IsParseRestricted() {
        return errors.NewParseRestrictedModeInvalidAccessFailure(spanName)
    }

    return nil
}

// Utility functions used for checking unexpected operation access while
// cadence is parsing programs.
//
// The generic functions are of the form
//      parseRestrict<x>Arg<y>Ret(txnState, spanName, callback, arg1, ..., argX)
// where the callback expects <x> number of arguments, and <y> number of
// return values (not counting error). If the callback expects no argument,
// <x>Arg is omitted, and similarly for return value.`

func generateWrapper(numArgs int, numRets int, content *FileContent) {
	l := content.Line
	push := content.PushIndent
	pop := content.PopIndent

	argsFuncSuffix := ""
	if numArgs > 0 {
		argsFuncSuffix = fmt.Sprintf("%dArg", numArgs)
	}

	argTypes := []string{}
	argNames := []string{}
	for i := 0; i < numArgs; i++ {
		argTypes = append(argTypes, fmt.Sprintf("Arg%dT", i))
		argNames = append(argNames, fmt.Sprintf("arg%d", i))
	}

	retsFuncSuffix := ""
	if numRets > 0 {
		retsFuncSuffix = fmt.Sprintf("%dRet", numRets)
	}

	retTypes := []string{}
	retNames := []string{}
	for i := 0; i < numRets; i++ {
		retTypes = append(retTypes, fmt.Sprintf("Ret%dT", i))
		retNames = append(retNames, fmt.Sprintf("value%d", i))
	}

	//
	// Generate function signature
	//

	l("")
	l("func parseRestrict%s%s[", argsFuncSuffix, retsFuncSuffix)
	push()

	for _, typeName := range append(argTypes, retTypes...) {
		l("%s any,", typeName)
	}

	pop()
	l("](")
	push()

	l("txnState *state.TransactionState,")
	l("spanName trace.SpanName,")

	callbackRet := "error"
	if numRets > 0 {
		callbackRet = "(" + strings.Join(append(retTypes, "error"), ", ") + ")"
	}

	l("callback func(%s) %s,", strings.Join(argTypes, ", "), callbackRet)

	for i, argType := range argTypes {
		l("%s %s,", argNames[i], argType)
	}

	pop()
	if numRets == 0 {
		l(") error {")
	} else {
		l(") (")
		push()

		for _, retType := range retTypes {
			l("%s,", retType)
		}
		l("error,")

		pop()
		l(") {")
	}
	push()

	//
	// Generate parse restrict check
	//

	l("err := parseRestricted(txnState, spanName)")
	l("if err != nil {")
	push()

	for i, retType := range retTypes {
		l("var %s %s", retNames[i], retType)
	}

	l("return %s", strings.Join(append(retNames, "err"), ", "))

	pop()
	l("}")

	//
	// Generate callback invocation
	//

	l("")
	l("return callback(%s)", strings.Join(argNames, ", "))

	pop()
	l("}")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <output file>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	cmd := append([]string{filepath.Base(os.Args[0])}, os.Args[1:]...)

	content := NewFileContent()
	content.Section(header, strings.Join(cmd, " "))

	for numArgs := 1; numArgs < 4; numArgs++ {
		generateWrapper(numArgs, 0, content)
	}

	for _, numArgs := range []int{0, 1, 2, 3, 4, 6} {
		generateWrapper(numArgs, 1, content)
	}

	generateWrapper(1, 2, content)

	buffer := &bytes.Buffer{}
	_, err := content.WriteTo(buffer)
	if err != nil {
		panic(err) // This should never happen
	}

	source, formatErr := format.Source(buffer.Bytes())

	// NOTE: formatting error can occur if the generated code has syntax
	// errors.  We still want to write out the unformatted source for debugging
	// purpose.
	if formatErr != nil {
		source = buffer.Bytes()
	}

	writeErr := os.WriteFile(os.Args[1], source, 0644)
	if writeErr != nil {
		panic(writeErr)
	}

	if formatErr != nil {
		panic(formatErr)
	}
}
