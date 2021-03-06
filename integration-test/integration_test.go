package test

import (
	"os"
	"testing"

	pinejs "github.com/resin-io/pinejs-client-go"
	"github.com/resin-io/pinejs-client-go/resin"
)

func TestGetDevice(t *testing.T) {
	devices := []resin.Device{}

	pine := pinejs.NewClient(os.Getenv("API_ENDPOINT"), os.Getenv("API_KEY"))

	if err := pine.List(&devices, pinejs.NewQueryOptions(pinejs.Filter, "uuid eq '11111ad359ebf8921a131b9781e01bf0fb4ef7d15437304d49efded7507f16'")...); err != nil {
		t.Fatal(err)
	} else {
		t.Log("Success!")
		t.Logf("%+v", devices)
	}
}

func TestGetDeviceWithToken(t *testing.T) {
	devices := []resin.Device{}

	pine := pinejs.NewClientWithToken(os.Getenv("API_ENDPOINT"), os.Getenv("API_TOKEN"))

	if err := pine.List(&devices, pinejs.NewQueryOptions(pinejs.Filter, "uuid eq '11111ad359ebf8921a131b9781e01bf0fb4ef7d15437304d49efded7507f16'")...); err != nil {
		t.Fatal(err)
	} else {
		t.Log("Success!")
		t.Logf("%+v", devices)
	}
}
