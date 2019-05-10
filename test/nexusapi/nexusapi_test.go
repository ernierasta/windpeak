package nexusapitest

import (
	"testing"

	"github.com/ernierasta/windpeak/nexusapi"
)

// TestMain is not real test, i used it to design api
func TestMain(t *testing.T) {

	apikey, err := nexusapi.Register("Windpeak")
	if err != nil {
		t.Error(err)
	}
	_ = apikey

}
