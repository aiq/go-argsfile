package argsfile

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {

	cnt := `
	
--longopt
value
# comment
-opt
file name with spaces.txt

# other comment
--para=
| pval
--arg
|= Full

$ --fs "sys rem" -para="hello world"

# flow text
lorem ipsum
|s lorem ipsum
|t 123
|t 456
|n muspi merol
	
`

	args, err := Read(strings.NewReader(cnt))
	exp := []string{"--longopt", "value", "-opt",
		"file name with spaces.txt", "--para=pval", "--arg=Full", "--fs",
		"sys rem", "-para=hello world",
		"lorem ipsum lorem ipsum\t123\t456\nmuspi merol",
	}

	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}

	for _, m := range cmpStrSlc(args, exp) {
		t.Error(m)
	}
}
