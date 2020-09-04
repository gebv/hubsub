package hubsub

import (
	"fmt"
	"strings"
	"testing"
)

func TestFilterFromPairs(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		paires  []string
		wantNil bool
		in      map[string]string
		want    bool
	}{
		{[]string{}, true, nil, false},
		{[]string{"a"}, true, nil, false},
		{[]string{"a", "b", "c"}, true, nil, false},

		{[]string{"a", "b"}, false, map[string]string{"a": "b"}, true},
		{[]string{"a", "b", "c", "d"}, false, map[string]string{"a": "b"}, false},
		{[]string{"a", "b"}, false, map[string]string{"a": "b", "c": "d"}, true},
		{[]string{"a", "b"}, false, map[string]string{"a": "c"}, false},
		{[]string{"a", "c"}, false, map[string]string{"a": "b"}, false},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("paires(%s)~(%v)", strings.Join(tt.paires, ","), tt.in)
		t.Run(name, func(t *testing.T) {
			fn := FilterFromPairs(tt.paires...)
			if tt.wantNil {
				if fn != nil {
					t.Errorf("FilterFromPairs(%v) want nil", tt.paires)
				}
			} else {
				got := fn(tt.in)
				if tt.want && !got {
					t.Errorf("fn(%v) != %v, want %v", tt.in, tt.want, got)
					return
				}
			}

		})
	}
}
