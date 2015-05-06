package sourcegraph

type MockBuildsService struct {
	Get_            func(build BuildSpec, opt *BuildGetOptions) (*Build, Response, error)
	List_           func(opt *BuildListOptions) ([]*Build, Response, error)
	Create_         func(repoRev RepoRevSpec, opt *BuildCreateOptions) (*Build, Response, error)
	Update_         func(build BuildSpec, info BuildUpdate) (*Build, Response, error)
	ListBuildTasks_ func(build BuildSpec, opt *BuildTaskListOptions) ([]*BuildTask, Response, error)
	CreateTasks_    func(build BuildSpec, tasks []*BuildTask) ([]*BuildTask, Response, error)
	UpdateTask_     func(task TaskSpec, info TaskUpdate) (*BuildTask, Response, error)
	GetLog_         func(build BuildSpec, opt *BuildGetLogOptions) (*LogEntries, Response, error)
	GetTaskLog_     func(task TaskSpec, opt *BuildGetLogOptions) (*LogEntries, Response, error)
	DequeueNext_    func() (*Build, Response, error)
}

func (s MockBuildsService) Get(build BuildSpec, opt *BuildGetOptions) (*Build, Response, error) {
	return s.Get_(build, opt)
}

func (s MockBuildsService) List(opt *BuildListOptions) ([]*Build, Response, error) {
	return s.List_(opt)
}

func (s MockBuildsService) Create(repoRev RepoRevSpec, opt *BuildCreateOptions) (*Build, Response, error) {
	return s.Create_(repoRev, opt)
}

func (s MockBuildsService) Update(build BuildSpec, info BuildUpdate) (*Build, Response, error) {
	return s.Update_(build, info)
}

func (s MockBuildsService) ListBuildTasks(build BuildSpec, opt *BuildTaskListOptions) ([]*BuildTask, Response, error) {
	return s.ListBuildTasks_(build, opt)
}

func (s MockBuildsService) CreateTasks(build BuildSpec, tasks []*BuildTask) ([]*BuildTask, Response, error) {
	return s.CreateTasks_(build, tasks)
}

func (s MockBuildsService) UpdateTask(task TaskSpec, info TaskUpdate) (*BuildTask, Response, error) {
	return s.UpdateTask_(task, info)
}

func (s MockBuildsService) GetLog(build BuildSpec, opt *BuildGetLogOptions) (*LogEntries, Response, error) {
	return s.GetLog_(build, opt)
}

func (s MockBuildsService) GetTaskLog(task TaskSpec, opt *BuildGetLogOptions) (*LogEntries, Response, error) {
	return s.GetTaskLog_(task, opt)
}

func (s MockBuildsService) DequeueNext() (*Build, Response, error) { return s.DequeueNext_() }
