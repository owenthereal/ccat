package sourcegraph

import (
	"fmt"
	"html/template"
	"log"
	"path"
	"time"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-nnz/nnz"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/graph"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/store"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/unit"
)

// DefsService communicates with the def- and graph-related endpoints in
// the Sourcegraph API.
type DefsService interface {
	// Get fetches a def.
	Get(def DefSpec, opt *DefGetOptions) (*Def, Response, error)

	// List defs.
	List(opt *DefListOptions) ([]*Def, Response, error)

	// ListRefs lists references to def.
	ListRefs(def DefSpec, opt *DefListRefsOptions) ([]*Ref, Response, error)

	// ListExamples lists examples for def.
	ListExamples(def DefSpec, opt *DefListExamplesOptions) ([]*Example, Response, error)

	// ListExamples lists people who committed parts of def's definition.
	ListAuthors(def DefSpec, opt *DefListAuthorsOptions) ([]*AugmentedDefAuthor, Response, error)

	// ListClients lists people who use def in their code.
	ListClients(def DefSpec, opt *DefListClientsOptions) ([]*AugmentedDefClient, Response, error)

	// ListDependents lists repositories that use def in their code.
	ListDependents(def DefSpec, opt *DefListDependentsOptions) ([]*AugmentedDefDependent, Response, error)

	// ListVersions lists all available versions of a definition in
	// the various repository commits in which it has appeared.
	//
	// TODO(sqs): how to deal with renames, etc.?
	ListVersions(def DefSpec, opt *DefListVersionsOptions) ([]*Def, Response, error)
}

// DefSpec specifies a def.
type DefSpec struct {
	Repo     string
	CommitID string
	UnitType string
	Unit     string
	Path     string
}

func (s *DefSpec) RouteVars() map[string]string {
	m := map[string]string{"RepoSpec": s.Repo, "UnitType": s.UnitType, "Unit": s.Unit, "Path": s.Path}
	if s.CommitID != "" {
		m["Rev"] = s.CommitID
	}
	return m
}

// DefKey returns the def key specified by s, using the Repo, UnitType,
// Unit, and Path fields of s.
func (s *DefSpec) DefKey() graph.DefKey {
	if s.Repo == "" {
		panic("Repo is empty")
	}
	if s.UnitType == "" {
		panic("UnitType is empty")
	}
	if s.Unit == "" {
		panic("Unit is empty")
	}
	return graph.DefKey{
		Repo:     s.Repo,
		CommitID: s.CommitID,
		UnitType: s.UnitType,
		Unit:     s.Unit,
		Path:     s.Path,
	}
}

// NewDefSpecFromDefKey returns a DefSpec that specifies the same
// def as the given key.
func NewDefSpecFromDefKey(key graph.DefKey) DefSpec {
	return DefSpec{
		Repo:     key.Repo,
		CommitID: key.CommitID,
		UnitType: key.UnitType,
		Unit:     key.Unit,
		Path:     key.Path,
	}
}

// defsService implements DefsService.
type defsService struct {
	client *Client
}

var _ DefsService = &defsService{}

// Def is a code def returned by the Sourcegraph API.
type Def struct {
	graph.Def

	Stat graph.Stats `json:",omitempty"`

	DocHTML string `json:",omitempty"`

	FmtStrings *DefFormatStrings `json:",omitempty"`
}

// DefFormatStrings contains the various def format strings from the
// srclib def formatter.
type DefFormatStrings struct {
	Name                 QualFormatStrings
	Type                 QualFormatStrings
	NameAndTypeSeparator string
	Language             string `json:",omitempty"`
	DefKeyword           string `json:",omitempty"`
	Kind                 string `json:",omitempty"`
}

// QualFormatStrings holds the formatted string for each Qualification
// for a def (for either Name or Type).
type QualFormatStrings struct {
	Unqualified             string `json:",omitempty"`
	ScopeQualified          string `json:",omitempty"`
	DepQualified            string `json:",omitempty"`
	RepositoryWideQualified string `json:",omitempty"`
	LanguageWideQualified   string `json:",omitempty"`
}

// DefSpec returns the DefSpec that specifies s.
func (s *Def) DefSpec() DefSpec {
	spec := NewDefSpecFromDefKey(s.Def.DefKey)
	return spec
}

func (s *Def) XRefs() int { return s.Stat["xrefs"] }
func (s *Def) RRefs() int { return s.Stat["rrefs"] }
func (s *Def) URefs() int { return s.Stat["urefs"] }

// TotalRefs is the number of unique references of all kinds to s. It
// is computed as (xrefs + rrefs), omitting urefs to avoid double-counting
// references in the same repository.
//
// The number of examples for s is usually TotalRefs() - 1, since the definition
// of a def counts as a ref but not an example.
func (s *Def) TotalRefs() int { return s.XRefs() + s.RRefs() }

func (s *Def) TotalExamples() int { return s.TotalRefs() - 1 }

// DefGetOptions specifies options for DefsService.Get.
type DefGetOptions struct {
	Doc bool `url:",omitempty"`

	// Stats is whether the Def response object should include statistics.
	Stats bool `url:",omitempty"`
}

func (s *defsService) Get(def DefSpec, opt *DefGetOptions) (*Def, Response, error) {
	url, err := s.client.URL(router.Def, def.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var def_ *Def
	resp, err := s.client.Do(req, &def_)
	if err != nil {
		return nil, resp, err
	}

	return def_, resp, nil
}

// DefListOptions specifies options for DefsService.List.
type DefListOptions struct {
	Name string `url:",omitempty" json:",omitempty"`

	// Specifies a search query for defs. If specified, then the Sort and Direction options are ignored
	Query string `url:",omitempty" json:",omitempty"`

	// ByteStart and ByteEnd will restrict the results to only definitions that overlap with the specified
	// start and end byte offsets. This filter is only applied if both values are set.
	ByteStart, ByteEnd uint32

	// DefKeys, if set, will return the definitions that match the given DefKey
	DefKeys []*graph.DefKey

	// RepoRevs constrains the results to a set of repository
	// revisions (given by their URIs plus an optional "@" and a
	// revision specifier). For example, "repo.com/foo@revspec".
	RepoRevs []string `url:",omitempty,comma" json:",omitempty"`

	UnitType string `url:",omitempty" json:",omitempty"`
	Unit     string `url:",omitempty" json:",omitempty"`

	Path string `url:",omitempty" json:",omitempty"`

	// File, if specified, will restrict the results to only defs defined in
	// the specified file.
	File string `url:",omitempty" json:",omitempty"`

	// FilePathPrefix, if specified, will restrict the results to only defs defined in
	// files whose path is underneath the specified prefix.
	FilePathPrefix string `url:",omitempty" json:",omitempty"`

	Kinds    []string `url:",omitempty,comma" json:",omitempty"`
	Exported bool     `url:",omitempty" json:",omitempty"`
	Nonlocal bool     `url:",omitempty" json:",omitempty"`

	// IncludeTest is whether the results should include definitions in test
	// files.
	IncludeTest bool `url:",omitempty" json:",omitempty"`

	// Enhancements
	Doc   bool `url:",omitempty" json:",omitempty"`
	Stats bool `url:",omitempty" json:",omitempty"`
	Fuzzy bool `url:",omitempty" json:",omitempty"`

	// Sorting
	Sort      string `url:",omitempty" json:",omitempty"`
	Direction string `url:",omitempty" json:",omitempty"`

	// Paging
	ListOptions
}

func (o *DefListOptions) DefFilters() []store.DefFilter {
	var fs []store.DefFilter
	if o.DefKeys != nil {
		fs = append(fs, store.DefFilterFunc(func(def *graph.Def) bool {
			for _, dk := range o.DefKeys {
				if (def.Repo == "" || def.Repo == dk.Repo) && (def.CommitID == "" || def.CommitID == dk.CommitID) &&
					(def.UnitType == "" || def.UnitType == dk.UnitType) && (def.Unit == "" || def.Unit == dk.Unit) &&
					def.Path == dk.Path {
					return true
				}
			}
			return false
		}))
	}
	if o.Name != "" {
		fs = append(fs, store.DefFilterFunc(func(def *graph.Def) bool {
			return def.Name == o.Name
		}))
	}
	if o.ByteEnd != 0 {
		fs = append(fs, store.DefFilterFunc(func(d *graph.Def) bool {
			return d.DefStart == o.ByteStart && d.DefEnd == o.ByteEnd
		}))
	}
	if o.Query != "" {
		fs = append(fs, store.ByDefQuery(o.Query))
	}
	if len(o.RepoRevs) > 0 {
		vs := make([]store.Version, len(o.RepoRevs))
		for i, repoRev := range o.RepoRevs {
			repo, commitID := ParseRepoAndCommitID(repoRev)
			if len(commitID) != 40 {
				log.Printf("WARNING: In DefListOptions.DefFilters, o.RepoRevs[%d]==%q has no commit ID or a non-absolute commit ID. No defs will match it.", i, repoRev)
			}
			vs[i] = store.Version{Repo: repo, CommitID: commitID}
		}
		fs = append(fs, store.ByRepoCommitIDs(vs...))
	}
	if o.Unit != "" && o.UnitType != "" {
		fs = append(fs, store.ByUnits(unit.ID2{Type: o.UnitType, Name: o.Unit}))
	}
	if (o.UnitType != "" && o.Unit == "") || (o.UnitType == "" && o.Unit != "") {
		log.Println("WARNING: DefListOptions.DefFilter: must specify either both or neither of --type and --name (to filter by source unit)")
	}
	if o.File != "" {
		fs = append(fs, store.ByFiles(path.Clean(o.File)))
	}
	if o.FilePathPrefix != "" {
		fs = append(fs, store.ByFiles(path.Clean(o.FilePathPrefix)))
	}
	if len(o.Kinds) > 0 {
		fs = append(fs, store.DefFilterFunc(func(def *graph.Def) bool {
			for _, kind := range o.Kinds {
				if def.Kind == kind {
					return true
				}
			}
			return false
		}))
	}
	if o.Exported {
		fs = append(fs, store.DefFilterFunc(func(def *graph.Def) bool {
			return def.Exported
		}))
	}
	if o.Nonlocal {
		fs = append(fs, store.DefFilterFunc(func(def *graph.Def) bool {
			return !def.Local
		}))
	}
	if !o.IncludeTest {
		fs = append(fs, store.DefFilterFunc(func(def *graph.Def) bool {
			return !def.Test
		}))
	}
	switch o.Sort {
	case "key":
		fs = append(fs, store.DefsSortByKey{})
	case "name":
		fs = append(fs, store.DefsSortByName{})
	}
	return fs
}

func (s *defsService) List(opt *DefListOptions) ([]*Def, Response, error) {
	url, err := s.client.URL(router.Defs, nil, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var defs []*Def
	resp, err := s.client.Do(req, &defs)
	if err != nil {
		return nil, resp, err
	}

	return defs, resp, nil
}

type Ref struct {
	graph.Ref
	Authorship *AuthorshipInfo
}

type Refs []*Ref

func (r *Ref) sortKey() string     { return fmt.Sprintf("%+v", r) }
func (vs Refs) Len() int           { return len(vs) }
func (vs Refs) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs Refs) Less(i, j int) bool { return vs[i].sortKey() < vs[j].sortKey() }

type DefListRefsOptions struct {
	Authorship bool   `url:",omitempty"` // whether to fetch authorship info about the refs
	Repo       string `url:",omitempty"` // only fetch refs from this repository URI
	ListOptions
}

func (s *defsService) ListRefs(def DefSpec, opt *DefListRefsOptions) ([]*Ref, Response, error) {
	url, err := s.client.URL(router.DefRefs, def.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var defRefs []*Ref
	resp, err := s.client.Do(req, &defRefs)
	if err != nil {
		return nil, resp, err
	}

	return defRefs, resp, nil
}

// Example is a usage example of a def.
type Example struct {
	graph.Ref

	// SrcHTML is the formatted HTML source code of the example, with links to
	// definitions.
	SrcHTML template.HTML `json:",omitempty"`

	// SourceCode contains the parsed source for this example, if requested via
	// DefListExamplesOptions.
	SourceCode *SourceCode `json:",omitempty"`

	// The line that the given example starts on
	StartLine int

	// The line that the given example ends on
	EndLine int

	// Error is whether an error occurred while fetching this example.
	Error bool
}

type Examples []*Example

func (r *Example) sortKey() string     { return fmt.Sprintf("%+v", r) }
func (vs Examples) Len() int           { return len(vs) }
func (vs Examples) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs Examples) Less(i, j int) bool { return vs[i].sortKey() < vs[j].sortKey() }

// DefListExamplesOptions specifies options for DefsService.ListExamples.
type DefListExamplesOptions struct {
	Formatted bool `url:",omitempty"`

	// Filter by a specific Repo URI
	Repo string `url:",omitempty"`

	// TokenizedSource requests that the source code be returned as a tokenized data
	// structure rather than an (annotated) string.
	//
	// This is useful when the client wants to take full control of rendering and manipulating
	// the contents.
	TokenizedSource bool `url:",omitempty"`

	ListOptions
}

func (s *defsService) ListExamples(def DefSpec, opt *DefListExamplesOptions) ([]*Example, Response, error) {
	url, err := s.client.URL(router.DefExamples, def.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var examples []*Example
	resp, err := s.client.Do(req, &examples)
	if err != nil {
		return nil, resp, err
	}

	return examples, resp, nil
}

type AuthorshipInfo struct {
	AuthorEmail    string    `db:"author_email"`
	LastCommitDate time.Time `db:"last_commit_date"`

	// LastCommitID is the commit ID of the last commit that this author made to
	// the thing that this info describes.
	LastCommitID string `db:"last_commit_id"`
}

type DefAuthorship struct {
	AuthorshipInfo

	// Exported is whether the def is exported.
	Exported bool

	Bytes           int
	BytesProportion float64
}

type DefAuthor struct {
	UID   nnz.Int
	Email nnz.String
	DefAuthorship
}

type DefAuthorsByBytes []*DefAuthor

func (v DefAuthorsByBytes) Len() int           { return len(v) }
func (v DefAuthorsByBytes) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v DefAuthorsByBytes) Less(i, j int) bool { return v[i].Bytes < v[j].Bytes }

type AugmentedDefAuthor struct {
	Person *Person
	*DefAuthor
}

// DefListAuthorsOptions specifies options for DefsService.ListAuthors.
type DefListAuthorsOptions struct {
	ListOptions
}

func (s *defsService) ListAuthors(def DefSpec, opt *DefListAuthorsOptions) ([]*AugmentedDefAuthor, Response, error) {
	url, err := s.client.URL(router.DefAuthors, def.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var authors []*AugmentedDefAuthor
	resp, err := s.client.Do(req, &authors)
	if err != nil {
		return nil, resp, err
	}

	return authors, resp, nil
}

type DefClient struct {
	UID   nnz.Int
	Email nnz.String

	AuthorshipInfo

	// UseCount is the number of times this person referred to the def.
	UseCount int `db:"use_count"`
}

type AugmentedDefClient struct {
	Person *Person
	*DefClient
}

// DefListClientsOptions specifies options for DefsService.ListClients.
type DefListClientsOptions struct {
	ListOptions
}

func (s *defsService) ListClients(def DefSpec, opt *DefListClientsOptions) ([]*AugmentedDefClient, Response, error) {
	url, err := s.client.URL(router.DefClients, def.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var clients []*AugmentedDefClient
	resp, err := s.client.Do(req, &clients)
	if err != nil {
		return nil, resp, err
	}

	return clients, resp, nil
}

type DefDependent struct {
	FromRepo string `db:"from_repo"`
	Count    int
}

type AugmentedDefDependent struct {
	Repo *Repo
	*DefDependent
}

// DefListDependentsOptions specifies options for DefsService.ListDependents.
type DefListDependentsOptions struct {
	ListOptions
}

func (s *defsService) ListDependents(def DefSpec, opt *DefListDependentsOptions) ([]*AugmentedDefDependent, Response, error) {
	url, err := s.client.URL(router.DefDependents, def.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var dependents []*AugmentedDefDependent
	resp, err := s.client.Do(req, &dependents)
	if err != nil {
		return nil, resp, err
	}

	return dependents, resp, nil
}

type DefListVersionsOptions struct {
	ListOptions
}

func (s *defsService) ListVersions(def DefSpec, opt *DefListVersionsOptions) ([]*Def, Response, error) {
	url, err := s.client.URL(router.DefVersions, def.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var defVersions []*Def
	resp, err := s.client.Do(req, &defVersions)
	if err != nil {
		return nil, resp, err
	}

	return defVersions, resp, nil
}

var _ DefsService = &MockDefsService{}
