package taskqueue

import (
	"fmt"
	"slices"
	"time"
)

type Priority int

const (
	PriorityLow    Priority = 1
	PriorityMedium Priority = 2
	PriorityHigh   Priority = 3
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID          string
	Name        string
	Priority    Priority
	Status      TaskStatus
	WorkerID    string
	CreatedAt   time.Time
	StartedAt   time.Time
	CompletedAt time.Time
	RetryCount  int
	MaxRetries  int
}

type Worker struct {
	ID        string
	Name      string
	IsActive  bool
	TaskCount int
	MaxTasks  int
}

type TaskQueue struct {
	tasks       map[string]*Task
	workers     map[string]*Worker
	lastTaskSeq int
	// Add fields as needed
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		tasks:       make(map[string]*Task),
		workers:     make(map[string]*Worker),
		lastTaskSeq: 0,
	}
}

// AddWorker adds a new worker to the system.
// Returns false if worker ID already exists or maxTasks <= 0.
func (tq *TaskQueue) AddWorker(id, name string, maxTasks int) bool {
	_, ok := tq.workers[id]
	if ok || maxTasks <= 0 {
		return false
	}
	tq.workers[id] = &Worker{
		ID:        id,
		Name:      name,
		IsActive:  true,
		TaskCount: 0,
		MaxTasks:  maxTasks,
	}
	return true
}

// SetWorkerActive sets a worker's active status.
// Inactive workers cannot be assigned new tasks.
// Returns false if worker doesn't exist.
func (tq *TaskQueue) SetWorkerActive(id string, active bool) bool {
	worker, ok := tq.workers[id]
	if !ok {
		return false
	}
	worker.IsActive = active
	return true
}

// GetWorker returns a worker by ID, or nil if not found.
func (tq *TaskQueue) GetWorker(id string) *Worker {
	return tq.workers[id]
}

// SubmitTask adds a new task to the queue with pending status.
// Returns task ID (format: "TASK-{sequential number}") and true if successful.
// Returns empty string and false if maxRetries < 0.
// Default maxRetries is 3 if not specified (pass -1 to use default, 0 means no retries).
func (tq *TaskQueue) SubmitTask(name string, priority Priority, maxRetries int) (string, bool) {
	if maxRetries == -1 {
		maxRetries = 3
	}

	if maxRetries < 0 {
		return "", false
	}

	tq.lastTaskSeq++
	newTaskId := fmt.Sprintf("TASK-%d", tq.lastTaskSeq)
	newTask := &Task{
		ID:         newTaskId,
		Name:       name,
		Priority:   priority,
		Status:     StatusPending,
		MaxRetries: maxRetries,
		CreatedAt:  time.Now(),
	}
	tq.tasks[newTaskId] = newTask

	return newTaskId, true
}

// GetTask returns a task by ID, or nil if not found.
func (tq *TaskQueue) GetTask(id string) *Task {
	return tq.tasks[id]
}

// AssignTask assigns the highest priority pending task to an available worker.
// Priority order: High > Medium > Low. For same priority, use FIFO (earliest CreatedAt first).
// A worker is available if: active, TaskCount < MaxTasks.
// Returns (taskID, workerID, true) if assignment made.
// Returns ("", "", false) if no pending tasks or no available workers.
// Updates task status to Running and sets StartedAt.
func (tq *TaskQueue) AssignTask() (string, string, bool) {
	if len(tq.tasks) == 0 {
		return "", "", false
	}

	tasks := make([]*Task, 0, len(tq.tasks))
	for _, t := range tq.tasks {
		if t.Status == StatusPending {
			tasks = append(tasks, t)
		}
	}

	if len(tasks) == 0 {
		return "", "", false
	}

	slices.SortFunc(tasks, func(a, b *Task) int {
		if a.Priority != b.Priority {
			return int(b.Priority) - int(a.Priority)
		}
		return a.CreatedAt.Compare(b.CreatedAt)
	})

	assigningTask := tasks[0]

	for _, worker := range tq.workers {
		if worker.IsActive && worker.TaskCount < worker.MaxTasks {
			assigningTask.Status = StatusRunning
			assigningTask.WorkerID = worker.ID
			assigningTask.StartedAt = time.Now()
			worker.TaskCount++
			return assigningTask.ID, worker.ID, true
		}
	}

	return "", "", false
}

// CompleteTask marks a task as completed.
// Returns false if task doesn't exist, not running, or workerID doesn't match.
// Decrements worker's TaskCount.
func (tq *TaskQueue) CompleteTask(taskID, workerID string) bool {
	task, ok := tq.tasks[taskID]
	if !ok || task.Status != StatusRunning || task.WorkerID != workerID {
		return false
	}

	worker, ok := tq.workers[workerID]
	if !ok {
		return false
	}

	worker.TaskCount--
	task.CompletedAt = time.Now()
	task.Status = StatusCompleted
	return true
}

// FailTask marks a task as failed and handles retry logic.
// If RetryCount < MaxRetries: increment RetryCount, set status back to Pending, clear WorkerID.
// If RetryCount >= MaxRetries: set status to Failed.
// Returns false if task doesn't exist, not running, or workerID doesn't match.
// Decrements worker's TaskCount in both cases.
func (tq *TaskQueue) FailTask(taskID, workerID string) bool {
	task, ok := tq.tasks[taskID]
	if !ok || task.Status != StatusRunning || task.WorkerID != workerID {
		return false
	}

	worker, ok := tq.workers[workerID]
	if !ok {
		return false
	}

	if task.RetryCount < task.MaxRetries {
		task.Status = StatusPending
		task.WorkerID = ""
		task.RetryCount++
		worker.TaskCount--
		return true
	}

	worker.TaskCount--
	task.Status = StatusFailed
	return true
}

// GetPendingTasks returns all pending tasks sorted by priority (high to low),
// then by CreatedAt (earliest first) for same priority.
func (tq *TaskQueue) GetPendingTasks() []*Task {
	tasks := make([]*Task, 0, len(tq.tasks))
	for _, t := range tq.tasks {
		if t.Status == StatusPending {
			tasks = append(tasks, t)
		}
	}

	slices.SortFunc(tasks, func(a, b *Task) int {
		if a.Priority != b.Priority {
			return int(b.Priority) - int(a.Priority)
		}
		return a.CreatedAt.Compare(b.CreatedAt)
	})

	return tasks
}

// GetWorkerTasks returns all tasks currently assigned to a worker (status = Running).
// Returns empty slice if worker doesn't exist or has no tasks.
func (tq *TaskQueue) GetWorkerTasks(workerID string) []*Task {
	tasks := make([]*Task, 0, len(tq.tasks))
	for _, t := range tq.tasks {
		if t.WorkerID != "" && t.WorkerID == workerID && t.Status == StatusRunning {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

// GetQueueStats returns statistics about the queue.
// Returns: (pendingCount, runningCount, completedCount, failedCount)
func (tq *TaskQueue) GetQueueStats() (int, int, int, int) {
	pendingCount := 0
	runningCount := 0
	completedCount := 0
	failedCount := 0

	for _, t := range tq.tasks {
		if t.Status == StatusPending {
			pendingCount++
		}

		if t.Status == StatusRunning {
			runningCount++
		}

		if t.Status == StatusCompleted {
			completedCount++
		}

		if t.Status == StatusFailed {
			failedCount++
		}
	}

	return pendingCount, runningCount, completedCount, failedCount
}

// ReassignAbandonedTasks finds all running tasks that have been running longer than
// the given duration and resets them to pending (simulating worker timeout).
// Decrements the original worker's TaskCount.
// Returns the count of tasks reassigned.
func (tq *TaskQueue) ReassignAbandonedTasks(timeout time.Duration) int {
	reassignCount := 0
	for _, t := range tq.tasks {
		if t.Status == StatusRunning && time.Since(t.StartedAt) > timeout {
			reassignCount++
			t.Status = StatusPending
			t.StartedAt = time.Time{}
			t.CompletedAt = time.Time{}
			if w, ok := tq.workers[t.WorkerID]; ok {
				w.TaskCount--
			}
			t.WorkerID = ""
		}
	}

	return reassignCount
}
