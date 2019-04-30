package library

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strings"
	"testing"
	"time"
	"github.com/stretchr/testify/require"
)

var (
	WaitPollInterval	= time.Second
	WaitPollTimeout		= 10 * time.Minute
)

func GenerateNameForTest(t *testing.T, prefix string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	require.NoError(t, err)
	name := []byte(fmt.Sprintf("%s%s-%016x", prefix, t.Name(), n.Int64()))
	name = regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAll(name, []byte("-"))
	name = regexp.MustCompile("-+").ReplaceAll(name, []byte("-"))
	return strings.Trim(string(name), "-")
}
