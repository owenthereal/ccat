package sourcegraph

import (
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/unit"
)

// UnitsService communicates with the source unit-related endpoints in
// the Sourcegraph API.
type UnitsService interface {
	// Get fetches a unit.
	Get(spec UnitSpec) (*unit.RepoSourceUnit, Response, error)

	// List units.
	List(opt *UnitListOptions) ([]*unit.RepoSourceUnit, Response, error)
}

// UnitSpec specifies a source unit.
type UnitSpec struct {
	RepoRevSpec
	UnitType string
	Unit     string
}

func NewUnitSpecFromUnit(u *unit.RepoSourceUnit) UnitSpec {
	return UnitSpec{
		RepoRevSpec: RepoRevSpec{
			RepoSpec: RepoSpec{URI: u.Repo},
			Rev:      u.CommitID,
			CommitID: u.CommitID,
		},
		UnitType: u.UnitType,
		Unit:     u.Unit,
	}
}

func UnmarshalUnitSpec(vars map[string]string) (UnitSpec, error) {
	repoRevSpec, err := UnmarshalRepoRevSpec(vars)
	if err != nil {
		return UnitSpec{}, err
	}
	return UnitSpec{
		RepoRevSpec: repoRevSpec,
		UnitType:    vars["UnitType"],
		Unit:        vars["Unit"],
	}, nil
}

func (s UnitSpec) RouteVars() map[string]string {
	v := s.RepoRevSpec.RouteVars()
	v["UnitType"] = s.UnitType
	v["Unit"] = s.Unit
	return v
}

// unitsService implements UnitsService.
type unitsService struct {
	client *Client
}

var _ UnitsService = &unitsService{}

// UnitListOptions specifies options for UnitsService.List.
type UnitListOptions struct {
	// RepoRevs constrains the results to a set of repository
	// revisions (given by their URIs plus an optional "@" and a
	// revision specifier). For example, "repo.com/foo@revspec".
	RepoRevs []string `url:",omitempty,comma" json:",omitempty"`

	UnitType string `url:",omitempty"`
	Unit     string `url:",omitempty"`

	// NameQuery specifies a full-text search query over the unit
	// name.
	NameQuery string `url:",omitempty" json:",omitempty"`

	// Query specifies a full-text search query over the repo URI,
	// unit name, and unit data.
	Query string `url:",omitempty" json:",omitempty"`

	// Paging
	ListOptions
}

func (s *unitsService) Get(spec UnitSpec) (*unit.RepoSourceUnit, Response, error) {
	url, err := s.client.URL(router.Unit, spec.RouteVars(), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var u unit.RepoSourceUnit
	resp, err := s.client.Do(req, &u)
	if err != nil {
		return nil, resp, err
	}

	return &u, resp, nil
}

func (s *unitsService) List(opt *UnitListOptions) ([]*unit.RepoSourceUnit, Response, error) {
	url, err := s.client.URL(router.Units, nil, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var units []*unit.RepoSourceUnit
	resp, err := s.client.Do(req, &units)
	if err != nil {
		return nil, resp, err
	}

	return units, resp, nil
}

var _ UnitsService = &MockUnitsService{}
