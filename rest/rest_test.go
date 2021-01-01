package rest

import (
	"testing"
)

func TestAmountOfEndpoints(t *testing.T) {
	defined := len(handlers)
	registered := registerEndpoints(handlers)
	if defined != registered {
		t.Errorf("There were %d endpoints defined, but %d endpoints registered", defined, registered)
	} else {
		t.Logf("All %d endpoints were registered", defined)
	}
}

