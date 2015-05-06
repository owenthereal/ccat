package sourcegraph

type MockDeltasService struct {
	Get_                    func(ds DeltaSpec, opt *DeltaGetOptions) (*Delta, Response, error)
	ListUnits_              func(ds DeltaSpec, opt *DeltaListUnitsOptions) ([]*UnitDelta, Response, error)
	ListDefs_               func(ds DeltaSpec, opt *DeltaListDefsOptions) (*DeltaDefs, Response, error)
	ListDependencies_       func(ds DeltaSpec, opt *DeltaListDependenciesOptions) (*DeltaDependencies, Response, error)
	ListFiles_              func(ds DeltaSpec, opt *DeltaListFilesOptions) (*DeltaFiles, Response, error)
	ListAffectedAuthors_    func(ds DeltaSpec, opt *DeltaListAffectedAuthorsOptions) ([]*DeltaAffectedPerson, Response, error)
	ListAffectedClients_    func(ds DeltaSpec, opt *DeltaListAffectedClientsOptions) ([]*DeltaAffectedPerson, Response, error)
	ListAffectedDependents_ func(ds DeltaSpec, opt *DeltaListAffectedDependentsOptions) ([]*DeltaAffectedRepo, Response, error)
	ListReviewers_          func(ds DeltaSpec, opt *DeltaListReviewersOptions) ([]*DeltaReviewer, Response, error)
	ListIncoming_           func(rr RepoRevSpec, opt *DeltaListIncomingOptions) ([]*Delta, Response, error)
}

func (s MockDeltasService) Get(ds DeltaSpec, opt *DeltaGetOptions) (*Delta, Response, error) {
	return s.Get_(ds, opt)
}

func (s MockDeltasService) ListUnits(ds DeltaSpec, opt *DeltaListUnitsOptions) ([]*UnitDelta, Response, error) {
	return s.ListUnits_(ds, opt)
}

func (s MockDeltasService) ListDefs(ds DeltaSpec, opt *DeltaListDefsOptions) (*DeltaDefs, Response, error) {
	return s.ListDefs_(ds, opt)
}

func (s MockDeltasService) ListDependencies(ds DeltaSpec, opt *DeltaListDependenciesOptions) (*DeltaDependencies, Response, error) {
	return s.ListDependencies_(ds, opt)
}

func (s MockDeltasService) ListFiles(ds DeltaSpec, opt *DeltaListFilesOptions) (*DeltaFiles, Response, error) {
	return s.ListFiles_(ds, opt)
}

func (s MockDeltasService) ListAffectedAuthors(ds DeltaSpec, opt *DeltaListAffectedAuthorsOptions) ([]*DeltaAffectedPerson, Response, error) {
	return s.ListAffectedAuthors_(ds, opt)
}

func (s MockDeltasService) ListAffectedClients(ds DeltaSpec, opt *DeltaListAffectedClientsOptions) ([]*DeltaAffectedPerson, Response, error) {
	return s.ListAffectedClients_(ds, opt)
}

func (s MockDeltasService) ListAffectedDependents(ds DeltaSpec, opt *DeltaListAffectedDependentsOptions) ([]*DeltaAffectedRepo, Response, error) {
	return s.ListAffectedDependents_(ds, opt)
}

func (s MockDeltasService) ListReviewers(ds DeltaSpec, opt *DeltaListReviewersOptions) ([]*DeltaReviewer, Response, error) {
	return s.ListReviewers_(ds, opt)
}

func (s MockDeltasService) ListIncoming(rr RepoRevSpec, opt *DeltaListIncomingOptions) ([]*Delta, Response, error) {
	return s.ListIncoming_(rr, opt)
}
