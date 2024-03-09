package component

import (
	"fmt"
	"testing"
)

func CompareContainers(t *testing.T, got, want *Container) {
	t.Helper()

	if fmt.Sprintf("%v", got.layout) != fmt.Sprintf("%v", want.layout) {
		t.Errorf("Container.UnmarshalXML() got = %v, want = %v", got.layout, want.layout)
	}

	if got.width != want.width {
		t.Errorf("Container.UnmarshalXML() got = %v, want = %v", got.width, want.width)
	}

	if got.height != want.height {
		t.Errorf("Container.UnmarshalXML() got = %v, want = %v", got.height, want.height)
	}

	if got.padding != want.padding {
		t.Errorf("Container.UnmarshalXML() got = %v, want = %v", got.padding, want.padding)
	}
}
