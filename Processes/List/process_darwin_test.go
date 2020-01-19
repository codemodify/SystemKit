// +build darwin

package List

import (
	"testing"
)

func TestDarwinProcess_impl(t *testing.T) {
	var _ Process = new(DarwinProcess)
}
