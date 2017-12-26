package key

import (
	"bytes"
	"testing"

	"github.com/dtynn/winston/internal/pb"
	"github.com/golang/protobuf/proto"
)

type testUnmarshaler struct {
	pb.Result
}

func (u *testUnmarshaler) UnmarshalBinary(data []byte) error {
	return proto.Unmarshal(data, &u.Result)
}

func TestKey(t *testing.T) {
	cases := []struct {
		prefix []byte
		str    string
	}{
		{
			nil,
			"1234",
		},
		{
			[]byte(""),
			"1234",
		},
		{
			[]byte("_node_"),
			"1234",
		},
	}

	for i, c := range cases {
		key := Key(c.prefix, bytes.NewBufferString(c.str))
		expected := string(c.prefix) + c.str
		if got := string(key); got != expected {
			t.Errorf("#%d expected %s, got %s", i+1, expected, got)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	res := &testUnmarshaler{
		pb.Result{
			Code:    pb.ResultCode_ResultCodeUnknown,
			Message: "testing",
		},
	}

	data, err := proto.Marshal(&res.Result)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer(data)

	malformedPrefix := []byte("_prefix_")

	prefixes := [][]byte{
		nil,
		[]byte(""),
		[]byte("_node_"),
	}

	for i, prefix := range prefixes {
		key := Key(prefix, buf)
		got := testUnmarshaler{}

		if perr := Unmarshal(key, malformedPrefix, &got); perr != ErrMalformedKeyPrefix {
			t.Fatalf("#%d unexpected error for malformed prefix %s", i+1, err)
		}

		err := Unmarshal(key, prefix, &got)
		if err != nil {
			t.Fatalf("#%d unexpected unmarshal error %v", i+1, err)
		}

		if got.Code != res.Code || got.Message != res.Message {
			t.Fatalf("#%d unexpected message %v", i+1, got)
		}
	}
}
