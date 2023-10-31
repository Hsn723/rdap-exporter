package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeLabel(t *testing.T) {
	t.Parallel()
	cases := []struct {
		title  string
		value  string
		expect string
	}{
		{
			title:  "AlreadyNormalized",
			value:  "registration",
			expect: "registration",
		},
		{
			title:  "HasSpace",
			value:  "client transfer prohibited",
			expect: "client_transfer_prohibited",
		},
		{
			title:  "SpaceAndCaps",
			value:  "last update of RDAP database",
			expect: "last_update_of_rdap_database",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Helper()
			actual := normalizeLabel(tc.value)
			assert.Equal(t, tc.expect, actual)
		})
	}
}
