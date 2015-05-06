package graph

type testDefFormatter struct{}

func (f testDefFormatter) RepositoryListing(s *Def) RepositoryListingDef {
	return RepositoryListingDef{
		Name:    s.Name,
		SortKey: s.Name,
	}
}

func (_ testDefFormatter) KindName(s *Def) string                          { return "" }
func (_ testDefFormatter) LanguageName(s *Def) string                      { return "" }
func (_ testDefFormatter) QualifiedName(s *Def, relativeTo *DefKey) string { return "" }
func (f testDefFormatter) TypeString(s *Def) string                        { return "" }

// func TestFormatAndSortDefsForRepositoryListing(t *testing.T) {
// 	RegisterDefFormatter("t", testDefFormatter{})
// 	defer func() {
// 		DefFormatters = nil
// 	}()

// 	defs := []*Def{
// 		{DefKey: DefKey{UnitType: "t"}, Name: "z"},
// 		{DefKey: DefKey{UnitType: "t"}, Name: "a"},
// 	}

// 	want := map[*Def]RepositoryListingDef{
// 		defs[0]: RepositoryListingDef{Name: "z", NameLabel: "", Language: "", SortKey: "z"},
// 		defs[1]: RepositoryListingDef{Name: "a", NameLabel: "", Language: "", SortKey: "a"},
// 	}

// 	fmtDefs := FormatAndSortDefsForRepositoryListing(defs)

// 	// Check that fmtDefs is sorted (was [z,a], should be [a,z]).
// 	if s1 := defs[0]; s1.Name != "a" {
// 		t.Errorf("got sorted def1 name %q, want 'a'", s1.Name)
// 	}
// 	if s2 := defs[1]; s2.Name != "z" {
// 		t.Errorf("got sorted def2 name %q, want 'z'", s2.Name)
// 	}

// 	if !reflect.DeepEqual(fmtDefs, want) {
// 		t.Errorf("got formatted defs map %+v, want %+v", fmtDefs, want)
// 	}
// }
