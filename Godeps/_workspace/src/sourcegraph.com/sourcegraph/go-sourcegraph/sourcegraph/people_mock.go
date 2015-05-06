package sourcegraph

type MockPeopleService struct {
	Get_ func(person PersonSpec) (*Person, Response, error)
}

func (s MockPeopleService) Get(person PersonSpec) (*Person, Response, error) { return s.Get_(person) }
