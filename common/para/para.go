package para

// Workers configures a task executor with the most number of tasks to be executed in parallel.
type Workers struct {
  sem chan struct{}
}
