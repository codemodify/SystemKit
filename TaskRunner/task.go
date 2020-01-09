package TaskRunner

// Task -
//   IF `run` is not empty then exec
//   ELSE IF `seqTasks` is not empty then exec sequentially
//   ELSE exec `tasks` in parallel
type Task struct {
	ID              string      `json:"id"`        // Used by the UI renderers to update state
	Label           string      `json:"label"`     // Display label
	Tasks           []*Task     `json:"tasks"`     // Parallel tasks
	SeqTasks        []*Task     `json:"seqTasks"`  // Sequential tasks
	Run             string      `json:"run"`       // A runnable unit
	RunParams       interface{} `json:"runParams"` //
	RunTaskInstance ITask       `json:"-"`         // Used internally only, not imported or exported to JSON
}

// TaskState -
type TaskState struct {
	TaskID   string `json:"id"`
	Error    bool   `json:"error"`
	Message  string `json:"message"`
	Progress int    `json:"progress"`
}

// ITask -
type ITask interface {
	Run(tag string, callStackLevel int)
}
