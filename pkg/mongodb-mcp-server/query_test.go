package mongodbmcpserver

import "testing"

func TestNormalizeSelector_ObjectId(t *testing.T) {
	in := `{"_id":ObjectId("64f7e2d2ffe1269f269fa039")}`
	got := normalizeSelector(in)
	want := `{"_id":{"$oid":"64f7e2d2ffe1269f269fa039"}}`
	if got != want {
		t.Fatalf("normalizeSelector(%q) = %q, want %q", in, got, want)
	}
}

func TestParseFilter_ObjectIdAndEJSON(t *testing.T) {
	oid := `{"_id":ObjectId("64f7e2d2ffe1269f269fa039")}`
	f, err := parseFilter(oid)
	if err != nil {
		t.Fatal(err)
	}
	if f["_id"] == nil {
		t.Fatal("expected _id")
	}

	ejson := `{"_id":{"$oid":"64f7e2d2ffe1269f269fa039"}}`
	f2, err := parseFilter(ejson)
	if err != nil {
		t.Fatal(err)
	}
	_ = f2
}
