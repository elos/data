package mongo

import (
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/elos/autonomous"
)

var (
	mongod exec.Cmd
)

type runner struct {
	autonomous.Life
	autonomous.Stopper
	autonomous.Managed

	mongod     *exec.Cmd
	ConfigFile string
	sync.Mutex
}

var Runner = &runner{
	Life:    autonomous.NewLife(),
	Stopper: make(autonomous.Stopper),
	Managed: *new(autonomous.Managed),
}

func (r *runner) SetConfigFile(s string) {
	r.Lock()
	defer r.Unlock()

	r.ConfigFile = s
}

func (r *runner) Start() {
	// Lock the runner
	r.Lock()
	defer r.Unlock()

	if r.ConfigFile != "" {
		r.mongod = exec.Command("mongod", "--config", r.ConfigFile)
	} else {
		r.mongod = exec.Command("mongod")
	}

	r.mongod.Stdout = os.Stdout
	r.mongod.Stderr = os.Stderr

	if err := r.mongod.Start(); err != nil {
		log.Fatal(err)
		return
	}

	r.Life.Begin()
	log.Print("Mongo successfully started")

	<-r.Stopper

	if err := r.mongod.Process.Signal(os.Interrupt); err != nil {
		log.Fatal(err)
		return
	}

	r.Life.End()
	log.Print("Mongo succesfully stopped")
}

var lock = sync.Mutex{}

func Testify(r *runner) {
	lock.Lock()
	defer lock.Unlock()

	Runner.SetConfigFile("./test.conf")
	go Runner.Start()
	Runner.WaitStart()
}

func Detestify(r *runner) {
	lock.Lock()
	defer lock.Unlock()

	go r.Stop()
	r.WaitStop()
	r.SetConfigFile("")
}
