package sourcegraph

type MockDefsService struct {
	Get_            func(def DefSpec, opt *DefGetOptions) (*Def, Response, error)
	List_           func(opt *DefListOptions) ([]*Def, Response, error)
	ListRefs_       func(def DefSpec, opt *DefListRefsOptions) ([]*Ref, Response, error)
	ListExamples_   func(def DefSpec, opt *DefListExamplesOptions) ([]*Example, Response, error)
	ListAuthors_    func(def DefSpec, opt *DefListAuthorsOptions) ([]*AugmentedDefAuthor, Response, error)
	ListClients_    func(def DefSpec, opt *DefListClientsOptions) ([]*AugmentedDefClient, Response, error)
	ListDependents_ func(def DefSpec, opt *DefListDependentsOptions) ([]*AugmentedDefDependent, Response, error)
	ListVersions_   func(def DefSpec, opt *DefListVersionsOptions) ([]*Def, Response, error)
}

func (s MockDefsService) Get(def DefSpec, opt *DefGetOptions) (*Def, Response, error) {
	return s.Get_(def, opt)
}

func (s MockDefsService) List(opt *DefListOptions) ([]*Def, Response, error) { return s.List_(opt) }

func (s MockDefsService) ListRefs(def DefSpec, opt *DefListRefsOptions) ([]*Ref, Response, error) {
	return s.ListRefs_(def, opt)
}

func (s MockDefsService) ListExamples(def DefSpec, opt *DefListExamplesOptions) ([]*Example, Response, error) {
	return s.ListExamples_(def, opt)
}

func (s MockDefsService) ListAuthors(def DefSpec, opt *DefListAuthorsOptions) ([]*AugmentedDefAuthor, Response, error) {
	return s.ListAuthors_(def, opt)
}

func (s MockDefsService) ListClients(def DefSpec, opt *DefListClientsOptions) ([]*AugmentedDefClient, Response, error) {
	return s.ListClients_(def, opt)
}

func (s MockDefsService) ListDependents(def DefSpec, opt *DefListDependentsOptions) ([]*AugmentedDefDependent, Response, error) {
	return s.ListDependents_(def, opt)
}

func (s MockDefsService) ListVersions(def DefSpec, opt *DefListVersionsOptions) ([]*Def, Response, error) {
	return s.ListVersions_(def, opt)
}
