package gobravia

import (
	// "fmt"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
	// "io/ioutil"
	// "net/http"
	"testing"
)

var (
	testcommand string
)

func init() {
	testcommand = `{
  "id": 10,
  "result": [
    {
      "bundled": true,
      "type": "RM-J1100"
    },
    [
      {
        "name": "PowerOff",
        "value": "AAAAAQAAAAEAAAAvAw=="
      },
      {
        "name": "Input",
        "value": "AAAAAQAAAAEAAAAlAw=="
      }
    ]
  ]
}`

}

func TestGetBravia(t *testing.T) {
	brv := GetBravia("mytv.local", "0001", "FC:FF:FF:F2:FF:FF")

	if brv.Address != "mytv.local" {
		t.Log("Did not stored Address correct")
		t.Fail()
	}
	if brv.Pin != "0001" {
		t.Log("Did not stored Pin correct")
		t.Fail()
	}
}

func TestGetCommands(t *testing.T) {
	defer gock.Off()

	gock.New("http://blub").
		Post("/sony/system").
		MatchHeader("Content-Type", "application/json").
		Reply(200).
		BodyString(testcommand)
		// JSON(map[string]string{"foo": "bar"})

	b := *GetBravia("blub", "0000", "FC:FF:FF:F2:FF:FF")
	b.GetCommands()

	st.Expect(t, b.Commands["poweroff"], "AAAAAQAAAAEAAAAvAw==")
	st.Expect(t, b.Commands["input"], "AAAAAQAAAAEAAAAlAw==")
	st.Expect(t, len(b.Commands), 2)

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)
}

func TestSearchCommands(t *testing.T) {
	defer gock.Off()

	gock.New("http://blub").
		Post("/sony/system").
		MatchHeader("Content-Type", "application/json").
		Reply(200).
		BodyString(testcommand)
		// JSON(map[string]string{"foo": "bar"})

	b := *GetBravia("blub", "0000", "FC:FF:FF:F2:FF:FF")
	b.GetCommands()

	code, ok := b.SearchCode("poweroff")
	st.Expect(t, code, "AAAAAQAAAAEAAAAvAw==")
	st.Expect(t, ok, true)

	code, ok = b.SearchCode("powerpuff")
	st.Expect(t, code, "")
	st.Expect(t, ok, false)

	st.Expect(t, gock.IsDone(), true)
}
