// +build linux solaris

package List

import (
	"testing"
)

func TestUnixProcess_impl(t *testing.T) {
	var _ Process = new(UnixProcess)
}
