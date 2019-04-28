package operator

import (
	"testing"
)

func TestNothing(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
