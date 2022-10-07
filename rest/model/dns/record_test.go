package dns

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var marshalRecordCases = []struct {
	name    string
	record  *Record
	answers []*Answer
	out     []byte
}{
	{
		"marshalCAARecord",
		NewRecord("example.com", "caa.example.com", "CAA"),
		[]*Answer{NewCAAAnswer(0, "issue", "letsencrypt.org")},
		[]byte(`{"meta":{},"zone":"example.com","domain":"caa.example.com","type":"CAA","answers":[{"meta":{},"answer":["0","issue","letsencrypt.org"]}],"filters":[]}`),
	},
	{
		"marshalURLFWDRecord",
		NewRecord("example.com", "fwd.example.com", "URLFWD"),
		[]*Answer{
			NewURLFWDAnswer("/net", "https://example.net", 301, 1, 1),
			NewURLFWDAnswer("/org", "https://example.org", 302, 2, 0),
		},
		[]byte(`{"answers":[{"answer":["/net","https://example.net",301,1,1],"meta":{}},{"answer":["/org","https://example.org",302,2,0],"meta":{}}],"meta":{},"zone":"example.com","domain":"fwd.example.com","type":"URLFWD","filters":[]}`),
	},
}

func TestMarshalRecords(t *testing.T) {
	for _, tt := range marshalRecordCases {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.answers {
				tt.record.AddAnswer(tt.answers[i])
			}
			result, err := json.Marshal(tt.record)
			if err != nil {
				t.Error(err)
			}
			if !bytes.Equal(result, tt.out) {
				t.Errorf("got %q, want %q", result, tt.out)
			}
		})
	}
}
func TestMarshalRecordsOverrideTTL(t *testing.T) {
	trueb := true
	falseb := false
	var marshalALIASRecordCases = []struct {
		name         string
		record       *Record
		override_ttl *bool
		out          []byte
	}{
		{
			"marshalOverrideTTLNil",
			NewRecord("example.com", "example.com", "ALIAS"),
			nil,
			[]byte(`{"meta":{},"zone":"example.com","domain":"example.com","type":"ALIAS","answers":[],"filters":[]}`),
		},
		{
			"marshalOverrideTTLTrue",
			NewRecord("example.com", "example.com", "ALIAS"),
			&trueb,
			[]byte(`{"meta":{},"zone":"example.com","domain":"example.com","type":"ALIAS","override_ttl":true,"answers":[],"filters":[]}`),
		},
		{
			"marshalOverrideTTLFalse",
			NewRecord("example.com", "example.com", "ALIAS"),
			&falseb,
			[]byte(`{"meta":{},"zone":"example.com","domain":"example.com","type":"ALIAS","override_ttl":false,"answers":[],"filters":[]}`),
		},
	}
	for _, tt := range marshalALIASRecordCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.record.Override_TTL = tt.override_ttl
			result, err := json.Marshal(tt.record)
			if err != nil {
				t.Error(err)
			}
			if !bytes.Equal(result, tt.out) {
				t.Errorf("got %q, want %q", result, tt.out)
			}
		})
	}
}

func TestNewRecord(t *testing.T) {
	var CapitalLettersCases = []struct {
		name           string
		domain         string
		zone           string
		ExpectedDomain string
		ExpectedZone   string
	}{
		{
			"no cap letters",
			"testcase1.case",
			"testcase1.case",
			"testcase1.case",
			"testcase1.case",
		},
		{
			"domain as title",
			"Testcase2.case",
			"testcase2.case",
			"Testcase2.case",
			"testcase2.case",
		},
		{
			"zone with cap letters",
			"testcase3.case",
			"TeStCase3.case",
			"testcase3.case",
			"TeStCase3.case",
		},
		{
			"Domain no cap letters without zone",
			"test",
			"testcase4.case",
			"test.testcase4.case",
			"testcase4.case",
		},
		{
			"Domain with cap letters without zone",
			"TEST",
			"testcase4.case",
			"TEST.testcase4.case",
			"testcase4.case",
		},
	}
	for _, tt := range CapitalLettersCases {
		t.Run(tt.name, func(t *testing.T) {
			record := NewRecord(tt.zone, tt.domain, "A")
			assert.Equal(t, tt.ExpectedDomain, record.Domain)
			assert.Equal(t, tt.ExpectedZone, record.Zone)
		})
	}
}
