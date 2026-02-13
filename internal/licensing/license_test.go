package licensing

import "testing"

func TestHasModuleExactMatch(t *testing.T) {
	l := License{Modules: "m-a,m-aa,m-b"}

	if !l.HasModule("m-a") {
		t.Fatalf("expected module m-a to be found")
	}
	if l.HasModule("m") {
		t.Fatalf("did not expect partial module match")
	}
	if l.HasModule("m-c") {
		t.Fatalf("did not expect unknown module")
	}
}

func TestHasModuleTrimmedMatch(t *testing.T) {
	l := License{Modules: "m-a, m-b , m-c"}

	if !l.HasModule("m-b") {
		t.Fatalf("expected trimmed module m-b to be found")
	}
	if l.HasModule("") {
		t.Fatalf("did not expect empty module to be found")
	}
}
