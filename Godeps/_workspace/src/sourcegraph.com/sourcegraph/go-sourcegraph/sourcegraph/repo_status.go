package sourcegraph

import (
	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/go-github/github"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

type RepoStatus struct {
	github.RepoStatus
}

type CombinedStatus struct {
	github.CombinedStatus
}

func (s *repositoriesService) GetCombinedStatus(spec RepoRevSpec) (*CombinedStatus, Response, error) {
	url, err := s.client.URL(router.RepoCombinedStatus, spec.RouteVars(), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var status CombinedStatus
	resp, err := s.client.Do(req, &status)
	if err != nil {
		return nil, resp, err
	}

	return &status, resp, nil
}

func (s *repositoriesService) CreateStatus(spec RepoRevSpec, st RepoStatus) (*RepoStatus, Response, error) {
	url, err := s.client.URL(router.RepoStatusCreate, spec.RouteVars(), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", url.String(), st)
	if err != nil {
		return nil, nil, err
	}

	var created RepoStatus
	resp, err := s.client.Do(req, &created)
	if err != nil {
		return nil, resp, err
	}

	return &created, resp, nil
}
