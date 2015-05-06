package sourcegraph

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/db_common"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

func TestBuildsService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &Build{BID: 1}

	var called bool
	mux.HandleFunc(urlPath(t, router.Build, map[string]string{"BID": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	build, _, err := client.Builds.Get(BuildSpec{BID: 1}, nil)
	if err != nil {
		t.Errorf("Builds.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeBuildTime(build, want)
	if !reflect.DeepEqual(build, want) {
		t.Errorf("Builds.Get returned %+v, want %+v", build, want)
	}
}

func TestBuildsService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*Build{{BID: 1}}

	var called bool
	mux.HandleFunc(urlPath(t, router.Builds, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	builds, _, err := client.Builds.List(nil)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeBuildTime(builds...)
	normalizeBuildTime(want...)
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}

func TestBuildsService_Create(t *testing.T) {
	setup()
	defer teardown()

	config := &BuildCreateOptions{BuildConfig: BuildConfig{Import: true, Queue: true}, Force: true}
	want := &Build{BID: 123, Repo: 456}

	var called bool
	mux.HandleFunc(urlPath(t, router.RepoBuildsCreate, map[string]string{"RepoSpec": "r.com/x", "Rev": "c"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"Import":true,"Queue":true,"UseCache":false,"Priority":0,"Force":true}`+"\n")

		writeJSON(w, want)
	})

	build_, _, err := client.Builds.Create(RepoRevSpec{RepoSpec: RepoSpec{URI: "r.com/x"}, Rev: "c"}, config)
	if err != nil {
		t.Errorf("Builds.Create returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeBuildTime(build_)
	normalizeBuildTime(want)
	if !reflect.DeepEqual(build_, want) {
		t.Errorf("Builds.Create returned %+v, want %+v", build_, want)
	}
}

func TestBuildsService_Update(t *testing.T) {
	setup()
	defer teardown()

	update := BuildUpdate{Host: String("h")}
	want := &Build{BID: 123, Repo: 456}

	var called bool
	mux.HandleFunc(urlPath(t, router.BuildUpdate, map[string]string{"BID": "123"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"StartedAt":null,"EndedAt":null,"HeartbeatAt":null,"Host":"h","Success":null,"Purged":null,"Failure":null,"Killed":null,"Priority":null}`+"\n")

		writeJSON(w, want)
	})

	build, _, err := client.Builds.Update(BuildSpec{BID: 123}, update)
	if err != nil {
		t.Errorf("Builds.Update returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeBuildTime(build)
	normalizeBuildTime(want)
	if !reflect.DeepEqual(build, want) {
		t.Errorf("Builds.Update returned %+v, want %+v", build, want)
	}
}

func TestBuildsService_UpdateTask(t *testing.T) {
	setup()
	defer teardown()

	update := TaskUpdate{Success: Bool(true)}
	want := &BuildTask{BID: 123, TaskID: 456, CreatedAt: db_common.NullTime{}}

	var called bool
	mux.HandleFunc(urlPath(t, router.BuildTaskUpdate, map[string]string{"BID": "123", "TaskID": "456"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"StartedAt":null,"EndedAt":null,"Success":true,"Failure":null}`+"\n")

		writeJSON(w, want)
	})

	task, _, err := client.Builds.UpdateTask(TaskSpec{BuildSpec: BuildSpec{BID: 123}, TaskID: 456}, update)
	if err != nil {
		t.Errorf("Builds.UpdateTask returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}
	if !reflect.DeepEqual(task, want) {
		t.Errorf("Builds.UpdateTask returned %+v, want %+v", task, want)
	}
}

func TestBuildsService_CreateTasks(t *testing.T) {
	setup()
	defer teardown()

	create := []*BuildTask{
		{BID: 123, Op: "foo", UnitType: "t", Unit: "u"},
		{BID: 123, Op: "bar", UnitType: "t", Unit: "u"},
	}

	var called bool
	mux.HandleFunc(urlPath(t, router.BuildTasksCreate, map[string]string{"BID": "123"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `[{"BID":123,"UnitType":"t","Unit":"u","Op":"foo","CreatedAt":null,"StartedAt":null,"EndedAt":null,"Queue":false},{"BID":123,"UnitType":"t","Unit":"u","Op":"bar","CreatedAt":null,"StartedAt":null,"EndedAt":null,"Queue":false}]`+"\n")
		writeJSON(w, create)
	})

	tasks, _, err := client.Builds.CreateTasks(BuildSpec{BID: 123}, create)
	if err != nil {
		t.Errorf("Builds.CreateTasks returned error: %v", err)
	}
	if len(tasks) != len(create) {
		t.Error("len(tasks) != len(create)")
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestBuildsService_GetLog(t *testing.T) {
	setup()
	defer teardown()

	want := &LogEntries{MaxID: "1"}

	var called bool
	mux.HandleFunc(urlPath(t, router.BuildLog, map[string]string{"BID": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	entries, _, err := client.Builds.GetLog(BuildSpec{BID: 1}, nil)
	if err != nil {
		t.Errorf("Builds.GetLog returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(entries, want) {
		t.Errorf("Builds.GetLog returned %+v, want %+v", entries, want)
	}
}

func TestBuildsService_GetTaskLog(t *testing.T) {
	setup()
	defer teardown()

	want := &LogEntries{MaxID: "1"}

	var called bool
	mux.HandleFunc(urlPath(t, router.BuildTaskLog, map[string]string{"BID": "1", "TaskID": "2"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	entries, _, err := client.Builds.GetTaskLog(TaskSpec{BuildSpec: BuildSpec{BID: 1}, TaskID: 2}, nil)
	if err != nil {
		t.Errorf("Builds.GetTaskLog returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(entries, want) {
		t.Errorf("Builds.GetTaskLog returned %+v, want %+v", entries, want)
	}
}

func TestBuildsService_DequeueNext(t *testing.T) {
	setup()
	defer teardown()

	want := &Build{BID: 1}

	var called bool
	mux.HandleFunc(urlPath(t, router.BuildDequeueNext, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")

		writeJSON(w, want)
	})

	build, _, err := client.Builds.DequeueNext()
	if err != nil {
		t.Errorf("Builds.DequeueNext returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeBuildTime(build, want)
	if !reflect.DeepEqual(build, want) {
		t.Errorf("Builds.DequeueNext returned %+v, want %+v", build, want)
	}
}

func TestBuildsService_DequeueNext_emptyQueue(t *testing.T) {
	setup()
	defer teardown()

	var called bool
	mux.HandleFunc(urlPath(t, router.BuildDequeueNext, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusNotFound)
	})

	build, _, err := client.Builds.DequeueNext()
	if err != nil {
		t.Errorf("Builds.DequeueNext returned error: %v", err)
	}
	if build != nil {
		t.Errorf("got build %v, want nil (no builds in queue)", build)
	}

	if !called {
		t.Fatal("!called")
	}

}

func normalizeBuildTime(bs ...*Build) {
	for _, b := range bs {
		if b != nil {
			normalizeTime(&b.CreatedAt)
			normalizeTime(&b.StartedAt.Time)
			normalizeTime(&b.EndedAt.Time)
		}
	}
}
