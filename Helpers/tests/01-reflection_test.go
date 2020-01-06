package tests

import (
	"fmt"
	"testing"

	helpersReflect "github.com/codemodify/SystemKit/Helpers"
)

func Test_01(t *testing.T) {
	var thisFuncName = fmt.Sprintf("THIS FUNC: %s", helpersReflect.GetThisFuncName())

	t.Log(thisFuncName)
	fmt.Println(thisFuncName)
}
