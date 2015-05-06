package sourcegraph

type MockOrgsService struct {
	Get_            func(org OrgSpec) (*Org, Response, error)
	ListMembers_    func(org OrgSpec, opt *OrgListMembersOptions) ([]*User, Response, error)
	GetSettings_    func(org OrgSpec) (*OrgSettings, Response, error)
	UpdateSettings_ func(org OrgSpec, settings OrgSettings) (Response, error)
}

func (s MockOrgsService) Get(org OrgSpec) (*Org, Response, error) { return s.Get_(org) }

func (s MockOrgsService) ListMembers(org OrgSpec, opt *OrgListMembersOptions) ([]*User, Response, error) {
	return s.ListMembers_(org, opt)
}

func (s MockOrgsService) GetSettings(org OrgSpec) (*OrgSettings, Response, error) {
	return s.GetSettings_(org)
}

func (s MockOrgsService) UpdateSettings(org OrgSpec, settings OrgSettings) (Response, error) {
	return s.UpdateSettings_(org, settings)
}
