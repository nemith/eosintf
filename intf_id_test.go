package eosintf

import "testing"

func TestIntf(t *testing.T) {
	tt := []struct {
		input int
		want  string
	}{
		{0x000c0202, "Ethernet3/1/2"},
		{0x01ffffff, "Ethernet127/511/511"},
		{0x00000001, "Ethernet1"},
	}

	for _, tc := range tt {
		t.Run(tc.want, func(t *testing.T) {
			intf := Intf(tc.input)
			got := intf.String()

			if got != tc.want {
				t.Errorf("unexpected interface name (want %q, got %q)", tc.want, got)
			}
		})
	}
}
