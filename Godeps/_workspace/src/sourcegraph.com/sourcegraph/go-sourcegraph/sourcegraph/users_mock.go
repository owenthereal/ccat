package sourcegraph

type MockUsersService struct {
	Get_                   func(user UserSpec, opt *UserGetOptions) (*User, Response, error)
	GetSettings_           func(user UserSpec) (*UserSettings, Response, error)
	UpdateSettings_        func(user UserSpec, settings UserSettings) (Response, error)
	ListEmails_            func(user UserSpec) ([]*EmailAddr, Response, error)
	GetOrCreateFromGitHub_ func(user GitHubUserSpec, opt *UserGetOptions) (*User, Response, error)
	RefreshProfile_        func(userSpec UserSpec) (Response, error)
	ComputeStats_          func(userSpec UserSpec) (Response, error)
	List_                  func(opt *UsersListOptions) ([]*User, Response, error)
	ListAuthors_           func(user UserSpec, opt *UsersListAuthorsOptions) ([]*AugmentedPersonUsageByClient, Response, error)
	ListClients_           func(user UserSpec, opt *UsersListClientsOptions) ([]*AugmentedPersonUsageOfAuthor, Response, error)
	ListOrgs_              func(member UserSpec, opt *UsersListOrgsOptions) ([]*Org, Response, error)
}

func (s MockUsersService) Get(user UserSpec, opt *UserGetOptions) (*User, Response, error) {
	return s.Get_(user, opt)
}

func (s MockUsersService) GetSettings(user UserSpec) (*UserSettings, Response, error) {
	return s.GetSettings_(user)
}

func (s MockUsersService) UpdateSettings(user UserSpec, settings UserSettings) (Response, error) {
	return s.UpdateSettings_(user, settings)
}

func (s MockUsersService) ListEmails(user UserSpec) ([]*EmailAddr, Response, error) {
	return s.ListEmails_(user)
}

func (s MockUsersService) GetOrCreateFromGitHub(user GitHubUserSpec, opt *UserGetOptions) (*User, Response, error) {
	return s.GetOrCreateFromGitHub_(user, opt)
}

func (s MockUsersService) RefreshProfile(userSpec UserSpec) (Response, error) {
	return s.RefreshProfile_(userSpec)
}

func (s MockUsersService) ComputeStats(userSpec UserSpec) (Response, error) {
	return s.ComputeStats_(userSpec)
}

func (s MockUsersService) List(opt *UsersListOptions) ([]*User, Response, error) { return s.List_(opt) }

func (s MockUsersService) ListAuthors(user UserSpec, opt *UsersListAuthorsOptions) ([]*AugmentedPersonUsageByClient, Response, error) {
	return s.ListAuthors_(user, opt)
}

func (s MockUsersService) ListClients(user UserSpec, opt *UsersListClientsOptions) ([]*AugmentedPersonUsageOfAuthor, Response, error) {
	return s.ListClients_(user, opt)
}

func (s MockUsersService) ListOrgs(member UserSpec, opt *UsersListOrgsOptions) ([]*Org, Response, error) {
	return s.ListOrgs_(member, opt)
}
