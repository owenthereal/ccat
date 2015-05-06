package sourcegraph

import (
	"errors"
	"fmt"

	"strconv"
	"strings"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-nnz/nnz"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/db_common"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

// UsersService communicates with the users-related endpoints in the
// Sourcegraph API.
type UsersService interface {
	// Get fetches a user.
	Get(user UserSpec, opt *UserGetOptions) (*User, Response, error)

	// GetSettings fetches a user's configuration settings. If err is
	// nil, then the returned UserSettings must be non-nil.
	GetSettings(user UserSpec) (*UserSettings, Response, error)

	// UpdateSettings updates an user's configuration settings.
	UpdateSettings(user UserSpec, settings UserSettings) (Response, error)

	// ListEmails returns a list of a user's email addresses.
	ListEmails(user UserSpec) ([]*EmailAddr, Response, error)

	// GetOrCreateFromGitHub creates a new user based a GitHub user.
	GetOrCreateFromGitHub(user GitHubUserSpec, opt *UserGetOptions) (*User, Response, error)

	// RefreshProfile updates the user's profile information from external
	// sources, such as GitHub.
	//
	// This operation is performed asynchronously on the server side
	// (after receiving the request) and the API currently has no way
	// of notifying callers when the operation completes.
	RefreshProfile(userSpec UserSpec) (Response, error)

	// ComputeStats recomputes statistics about the user.
	//
	// This operation is performed asynchronously on the server side
	// (after receiving the request) and the API currently has no way
	// of notifying callers when the operation completes.
	ComputeStats(userSpec UserSpec) (Response, error)

	// List users.
	List(opt *UsersListOptions) ([]*User, Response, error)

	// ListAuthors lists users who authored code that user uses.
	ListAuthors(user UserSpec, opt *UsersListAuthorsOptions) ([]*AugmentedPersonUsageByClient, Response, error)

	// ListClients lists users who use code that user authored.
	ListClients(user UserSpec, opt *UsersListClientsOptions) ([]*AugmentedPersonUsageOfAuthor, Response, error)

	// ListOrgs lists organizations that a user is a member of.
	ListOrgs(member UserSpec, opt *UsersListOrgsOptions) ([]*Org, Response, error)
}

// User represents a registered user.
type User struct {
	// UID is the numeric primary key for a user.
	UID int `db:"uid"`

	// GitHubID is the numeric ID of the GitHub user account corresponding to
	// this user.
	GitHubID nnz.Int `db:"github_id"`

	// Login is the user's username, which typically corresponds to the user's
	// GitHub login.
	Login string

	// Name is the (possibly empty) full name of the user.
	Name string

	// Type is either "User" or "Organization".
	Type string

	// AvatarURL is the URL to an avatar image specified by the user.
	AvatarURL string

	// Location is the user's physical location (from their GitHub profile).
	Location string `json:",omitempty"`

	// Company is the user's company (from their GitHub profile).
	Company string `json:",omitempty"`

	// HomepageURL is the user's homepage or blog URL (from their GitHub
	// profile).
	HomepageURL string `db:"homepage_url" json:",omitempty"`

	// UserProfileDisabled is whether the user profile should not be displayed
	// on the Web app.
	UserProfileDisabled bool `db:"user_profile_disabled" json:",omitempty"`

	// RegisteredAt is the date that the user registered. If the user has not
	// registered (i.e., we have processed their repos but they haven't signed
	// into Sourcegraph), it is null.
	RegisteredAt db_common.NullTime `db:"registered_at"`

	// Stat contains statistics about this user. It's only filled in
	// by certain API responses.
	Stat PersonStats `db:"-" json:",omitempty"`
}

func (u *User) Spec() UserSpec {
	return UserSpec{Login: u.Login, UID: u.UID}
}

// GitHubLogin returns the user's Login. They are the same for now, but callers
// that intend to get the GitHub login should call GitHubLogin() so that we can
// decouple the logins in the future if needed.
func (u *User) GitHubLogin() string {
	if u.GitHubID == 0 {
		return ""
	}
	return u.Login
}

// IsOrganization is whether this user represents a GitHub organization
// (which are treated as a subclass of User in GitHub's data model).
func (u *User) IsOrganization() bool { return u.Type == "Organization" }

// AvatarURLOfSize returns the URL to an avatar for the user with the
// given width (in pixels).
func (u *User) AvatarURLOfSize(width int) string {
	return avatarURLOfSize(u.AvatarURL, width)
}

func avatarURLOfSize(avatarURL string, width int) string {
	return avatarURL + fmt.Sprintf("&s=%d", width)
}

// CanOwnRepositories is whether the user is capable of owning repositories
// (e.g., GitHub users can own GitHub repositories).
func (u *User) CanOwnRepositories() bool {
	return u.GitHubLogin() != ""
}

// CanAttributeCodeTo is whether this user can commit code. It is false for
// organizations and true for both users and transient users.
func (u *User) CanAttributeCodeTo() bool {
	return !u.IsOrganization()
}

// Person returns an equivalent Person.
func (u *User) Person() *Person {
	return &Person{
		PersonSpec: PersonSpec{UID: u.UID, Login: u.Login},
		FullName:   u.Name,
		AvatarURL:  u.AvatarURL,
	}
}

// UserSpec specifies a user. At least one of Login, and UID must be
// nonempty.
type UserSpec struct {
	// Login is a user's login.
	Login string

	// UID is a user's UID.
	UID int
}

// PathComponent returns the URL path component that specifies the user.
func (s *UserSpec) PathComponent() string {
	if s.Login != "" {
		return s.Login
	}
	if s.UID > 0 {
		return "$" + strconv.Itoa(s.UID)
	}
	panic("empty UserSpec")
}

func (s *UserSpec) RouteVars() map[string]string {
	return map[string]string{"UserSpec": s.PathComponent()}
}

// ParseUserSpec parses a string generated by (*UserSpec).String() and
// returns the equivalent UserSpec struct.
func ParseUserSpec(pathComponent string) (UserSpec, error) {
	if strings.Contains(pathComponent, "@") {
		return UserSpec{}, fmt.Errorf("UserSpec %q must not contain '@'")
	}
	if strings.HasPrefix(pathComponent, "$") {
		uid, err := strconv.Atoi(pathComponent[1:])
		return UserSpec{UID: uid}, err
	}
	return UserSpec{Login: pathComponent}, nil
}

// ErrUserNotExist is an error indicating that no such user exists.
var ErrUserNotExist = errors.New("user does not exist")

// ErrUserRenamed is an error type that indicates that a user account was renamed
// from OldLogin to NewLogin.
type ErrUserRenamed struct {
	// OldLogin is the previous login name.
	OldLogin string

	// NewLogin is what the old login was renamed to.
	NewLogin string
}

func (e ErrUserRenamed) Error() string {
	return fmt.Sprintf("login %q was renamed to %q; use the new name", e.OldLogin, e.NewLogin)
}

// usersService implements UsersService.
type usersService struct {
	client *Client
}

var _ UsersService = &usersService{}

type UserGetOptions struct {
	// Stats is whether to include statistics about the user in the response.
	Stats bool `url:",omitempty"`
}

func (s *usersService) Get(user_ UserSpec, opt *UserGetOptions) (*User, Response, error) {
	url, err := s.client.URL(router.User, user_.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var user__ *User
	resp, err := s.client.Do(req, &user__)
	if err != nil {
		return nil, resp, err
	}

	return user__, resp, nil
}

// EmailAddr is an email address associated with a user.
type EmailAddr struct {
	Email string // the email address (case-insensitively compared in the DB and API)

	Verified bool // whether this email address has been verified

	Primary bool // indicates this is the user's primary email (only 1 email can be primary per user)

	Guessed bool // whether Sourcegraph inferred via public data that this is an email for the user

	Blacklisted bool // indicates that this email should not be associated with the user (even if guessed in the future)
}

func (s *usersService) ListEmails(user UserSpec) ([]*EmailAddr, Response, error) {
	url, err := s.client.URL(router.UserEmails, user.RouteVars(), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var emails []*EmailAddr
	resp, err := s.client.Do(req, &emails)
	if err != nil {
		return nil, resp, err
	}

	return emails, resp, nil
}

// UserSettings describes a user's configuration settings.
type UserSettings struct {
	// RequestedUpgradeAt is the date on which a user requested an upgrade
	RequestedUpgradeAt db_common.NullTime `json:",omitempty"`

	PlanSettings `json:",omitempty"`
}

// PlanSettings describes the pricing plan that the user or org has selected.
type PlanSettings struct {
	PlanID *string `json:",omitempty"`
}

func (s *usersService) GetSettings(user UserSpec) (*UserSettings, Response, error) {
	url, err := s.client.URL(router.UserSettings, user.RouteVars(), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var settings *UserSettings
	resp, err := s.client.Do(req, &settings)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, nil
}

func (s *usersService) UpdateSettings(user UserSpec, settings UserSettings) (Response, error) {
	url, err := s.client.URL(router.UserSettingsUpdate, user.RouteVars(), nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), settings)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GitHubUserSpec specifies a GitHub user, either by GitHub login or GitHub user
// ID.
type GitHubUserSpec struct {
	Login string
	ID    int
}

func (s GitHubUserSpec) RouteVars() map[string]string {
	if s.ID != 0 {
		panic("GitHubUserSpec ID not supported via HTTP API")
	} else if s.Login != "" {
		return map[string]string{"GitHubUserSpec": s.Login}
	}
	panic("empty GitHubUserSpec")
}

func (s *usersService) GetOrCreateFromGitHub(user GitHubUserSpec, opt *UserGetOptions) (*User, Response, error) {
	url, err := s.client.URL(router.UserFromGitHub, user.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var user__ *User
	resp, err := s.client.Do(req, &user__)
	if err != nil {
		return nil, resp, err
	}

	return user__, resp, nil
}

func (s *usersService) RefreshProfile(user_ UserSpec) (Response, error) {
	url, err := s.client.URL(router.UserRefreshProfile, user_.RouteVars(), nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *usersService) ComputeStats(user_ UserSpec) (Response, error) {
	url, err := s.client.URL(router.UserComputeStats, user_.RouteVars(), nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UsersListOptions specifies options for the UsersService.List method.
type UsersListOptions struct {
	// Query filters the results to only those whose logins match. The
	// search algorithm is an implementation detail (currently it is a
	// prefix match).
	Query string `url:",omitempty" json:",omitempty"`

	Sort      string `url:",omitempty" json:",omitempty"`
	Direction string `url:",omitempty" json:",omitempty"`

	ListOptions
}

func (s *usersService) List(opt *UsersListOptions) ([]*User, Response, error) {
	url, err := s.client.URL(router.Users, nil, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

type PersonUsageByClient struct {
	AuthorUID   nnz.Int    `db:"author_uid"`
	AuthorEmail nnz.String `db:"author_email"`
	RefCount    int        `db:"ref_count"`
}

type AugmentedPersonUsageByClient struct {
	Author *Person
	*PersonUsageByClient
}

// UsersListAuthorsOptions specifies options for the UsersService.ListAuthors
// method.
type UsersListAuthorsOptions UsersListOptions

func (s *usersService) ListAuthors(user UserSpec, opt *UsersListAuthorsOptions) ([]*AugmentedPersonUsageByClient, Response, error) {
	url, err := s.client.URL(router.UserAuthors, user.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var people []*AugmentedPersonUsageByClient
	resp, err := s.client.Do(req, &people)
	if err != nil {
		return nil, resp, err
	}

	return people, resp, nil
}

type PersonUsageOfAuthor struct {
	ClientUID   nnz.Int    `db:"client_uid"`
	ClientEmail nnz.String `db:"client_email"`
	RefCount    int        `db:"ref_count"`
}

type AugmentedPersonUsageOfAuthor struct {
	Client *Person
	*PersonUsageOfAuthor
}

// UsersListClientsOptions specifies options for the UsersService.ListClients
// method.
type UsersListClientsOptions UsersListOptions

func (s *usersService) ListClients(user UserSpec, opt *UsersListClientsOptions) ([]*AugmentedPersonUsageOfAuthor, Response, error) {
	url, err := s.client.URL(router.UserClients, user.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var people []*AugmentedPersonUsageOfAuthor
	resp, err := s.client.Do(req, &people)
	if err != nil {
		return nil, resp, err
	}
	return people, resp, nil
}

type UsersListOrgsOptions struct {
	ListOptions
}

func (s *usersService) ListOrgs(member UserSpec, opt *UsersListOrgsOptions) ([]*Org, Response, error) {
	url, err := s.client.URL(router.UserOrgs, member.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var orgs []*Org
	resp, err := s.client.Do(req, &orgs)
	if err != nil {
		return nil, resp, err
	}

	return orgs, resp, nil
}

var _ UsersService = &MockUsersService{}
