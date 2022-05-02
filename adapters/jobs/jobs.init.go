package jobs

import (
	"reflect"
	"sync"
	"time"

	"github.com/deestarks/infiniti/application/services"
	"github.com/deestarks/infiniti/utils"
	"github.com/go-co-op/gocron"
)

// The use of "JobAdapter" in "JobInitializer" is a way to avoid the Init() method from calling itself if the same struct is used to initialize and run the jobs, which will cause a stack overflow.
type JobInitializer struct {
	adapter	*JobAdapter
}

type JobAdapter struct {
	appPort		services.AppServicesPort
	scheduler	*gocron.Scheduler
}

func (i *JobInitializer) NewJobAdapter() *JobAdapter {
	return i.adapter
}

func NewJobInitializer(port services.AppServicesPort) *JobInitializer {
	return &JobInitializer{
		adapter: &JobAdapter{
			appPort: port,
			scheduler: gocron.NewScheduler(time.UTC),
		},
	}
}

/*
	For each job, create a new struct with the name of the job as the struct name that has a bool field "Add" to indicate if the job should be added to the job queue, a "appPort" field of type "services.AppServicesPort" which will be used to access the application services, and a "scheduler" field of type "*gocron.Scheduler" which will be used to schedule the job.

	Return a new instance of the job using a factory method that have the "JobAdapter" as a receiver. Then implement a "Run" method on the job struct, which is where the job logic is implemented. The job will be automatically added to the job queue, if the "Add" field is set to true.

	The "Add" id used as a way to determine if the job should be added to the job queue and if the implementation is still being worked on, which in this case can be set to false. And when implementation of the job is completed to be added, the "Add" field should be set to true.

	Example:
			type YourJobName struct {
				Add			bool
				appPort		services.AppServicesPort
				scheduler	*gocron.Scheduler
			}

			func (adapter *JobAdapter) NewYourJobName() YourJobName {
				return YourJobName{
					Add: true, // Set to false if the implementation is still being worked on, or to dequeue
					appPort: adapter.appPort,
					scheduler: adapter.scheduler,
				}
			}

			func (job YourJobName) Run() {
				// Your job logic here
			}

	Init will auto run all jobs that have the "Add" field set to true
*/
func (inz *JobInitializer) Init() {
	var wg sync.WaitGroup
	adapter := inz.adapter
    jobs := reflect.ValueOf(adapter)

	if jobs.NumMethod() == 0 {
		utils.LogJobMessage("No jobs found")
		return
	}

	utils.LogJobMessage("Running %d job(s)...", jobs.NumMethod())
    for i := 0; i < jobs.NumMethod(); i++ {
		job := jobs.Method(i).Call(nil)

		// Check if job should be added to the job queue
		if job[0].FieldByName("Add").Bool() {
			utils.LogJobMessage("- Job started: \"%s\"", job[0].Type().Name())
			wg.Add(1) // Add to wait group
			go func() { // Run job in a goroutine
				defer wg.Done()
				job[0].MethodByName("Run").Call(nil)
			}()
		} else {
			utils.LogJobMessage("- Job skipped: \"%s\"", job[0].Type().Name())
		}
	}
	utils.LogJobMessage("All jobs initialized")
}