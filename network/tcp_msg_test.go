package network

import "testing"

func TestMsgParser_PackMsgWithMsg(t *testing.T) {
	_, err := NewMsgParser().PackMsgWithTrace([]byte("123"))
	if err != nil {
		t.Fatal(err)
		return
	}
}
