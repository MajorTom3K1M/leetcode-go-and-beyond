package taskqueue

import (
	"testing"
	"time"
)

func TestAddAndGetWorker(t *testing.T) {
	tq := NewTaskQueue()

	// Add worker successfully
	if !tq.AddWorker("w1", "Worker One", 5) {
		t.Error("should add worker successfully")
	}

	// Get worker
	w := tq.GetWorker("w1")
	if w == nil {
		t.Fatal("worker should not be nil")
	}
	if w.Name != "Worker One" || w.MaxTasks != 5 {
		t.Error("worker data mismatch")
	}
	if !w.IsActive {
		t.Error("new worker should be active by default")
	}
	if w.TaskCount != 0 {
		t.Error("new worker should have 0 tasks")
	}

	// Duplicate ID should fail
	if tq.AddWorker("w1", "Another", 3) {
		t.Error("duplicate worker ID should fail")
	}

	// Invalid maxTasks
	if tq.AddWorker("w2", "Bad Worker", 0) {
		t.Error("zero maxTasks should fail")
	}
	if tq.AddWorker("w3", "Bad Worker", -1) {
		t.Error("negative maxTasks should fail")
	}

	// Non-existent worker
	if tq.GetWorker("w999") != nil {
		t.Error("non-existent worker should return nil")
	}
}

func TestSetWorkerActive(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker One", 5)

	// Deactivate
	if !tq.SetWorkerActive("w1", false) {
		t.Error("should deactivate worker")
	}
	if tq.GetWorker("w1").IsActive {
		t.Error("worker should be inactive")
	}

	// Reactivate
	if !tq.SetWorkerActive("w1", true) {
		t.Error("should reactivate worker")
	}
	if !tq.GetWorker("w1").IsActive {
		t.Error("worker should be active")
	}

	// Non-existent worker
	if tq.SetWorkerActive("w999", true) {
		t.Error("non-existent worker should return false")
	}
}

func TestSubmitAndGetTask(t *testing.T) {
	tq := NewTaskQueue()

	// Submit with explicit maxRetries
	id1, ok := tq.SubmitTask("Task 1", PriorityHigh, 5)
	if !ok || id1 != "TASK-1" {
		t.Errorf("expected TASK-1 and true, got %s and %v", id1, ok)
	}

	task := tq.GetTask(id1)
	if task == nil {
		t.Fatal("task should exist")
	}
	if task.Name != "Task 1" || task.Priority != PriorityHigh {
		t.Error("task data mismatch")
	}
	if task.Status != StatusPending {
		t.Error("new task should be pending")
	}
	if task.MaxRetries != 5 {
		t.Errorf("expected maxRetries 5, got %d", task.MaxRetries)
	}
	if task.RetryCount != 0 {
		t.Error("new task should have 0 retry count")
	}

	// Submit with default maxRetries (-1)
	id2, ok := tq.SubmitTask("Task 2", PriorityLow, -1)
	if !ok {
		t.Error("should submit with default maxRetries")
	}
	if tq.GetTask(id2).MaxRetries != 3 {
		t.Error("default maxRetries should be 3")
	}

	// Submit with 0 maxRetries (valid - means no retries)
	id3, ok := tq.SubmitTask("Task 3", PriorityMedium, 0)
	if !ok {
		t.Error("should submit with 0 maxRetries")
	}
	if tq.GetTask(id3).MaxRetries != 0 {
		t.Error("maxRetries should be 0")
	}

	// Invalid maxRetries (less than -1)
	_, ok = tq.SubmitTask("Bad Task", PriorityLow, -2)
	if ok {
		t.Error("maxRetries < -1 should fail")
	}

	// Sequential IDs
	id4, _ := tq.SubmitTask("Task 4", PriorityLow, 1)
	if id4 != "TASK-4" {
		t.Errorf("expected TASK-4, got %s", id4)
	}

	// Non-existent task
	if tq.GetTask("TASK-999") != nil {
		t.Error("non-existent task should return nil")
	}
}

func TestAssignTask(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker One", 2)

	tq.SubmitTask("Low Priority", PriorityLow, 3)
	time.Sleep(1 * time.Millisecond)
	tq.SubmitTask("High Priority", PriorityHigh, 3)
	time.Sleep(1 * time.Millisecond)
	tq.SubmitTask("Medium Priority", PriorityMedium, 3)

	// Should assign high priority first
	taskID, workerID, ok := tq.AssignTask()
	if !ok {
		t.Fatal("should assign task")
	}
	task := tq.GetTask(taskID)
	if task.Priority != PriorityHigh {
		t.Error("should assign highest priority first")
	}
	if workerID != "w1" {
		t.Errorf("expected w1, got %s", workerID)
	}
	if task.Status != StatusRunning {
		t.Error("assigned task should be running")
	}
	if task.WorkerID != "w1" {
		t.Error("task should have worker assigned")
	}
	if task.StartedAt.IsZero() {
		t.Error("StartedAt should be set")
	}
	if tq.GetWorker("w1").TaskCount != 1 {
		t.Error("worker TaskCount should be 1")
	}

	// Should assign medium priority next
	taskID, _, ok = tq.AssignTask()
	if !ok {
		t.Fatal("should assign second task")
	}
	if tq.GetTask(taskID).Priority != PriorityMedium {
		t.Error("should assign medium priority second")
	}
	if tq.GetWorker("w1").TaskCount != 2 {
		t.Error("worker TaskCount should be 2")
	}

	// Worker at capacity - should not assign
	_, _, ok = tq.AssignTask()
	if ok {
		t.Error("should not assign when worker at capacity")
	}
}

func TestAssignTaskFIFO(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker One", 10)

	// Submit tasks with same priority
	tq.SubmitTask("First", PriorityMedium, 3)
	time.Sleep(1 * time.Millisecond)
	tq.SubmitTask("Second", PriorityMedium, 3)
	time.Sleep(1 * time.Millisecond)
	tq.SubmitTask("Third", PriorityMedium, 3)

	// Should assign in FIFO order
	taskID, _, _ := tq.AssignTask()
	if tq.GetTask(taskID).Name != "First" {
		t.Error("should assign first task first (FIFO)")
	}

	taskID, _, _ = tq.AssignTask()
	if tq.GetTask(taskID).Name != "Second" {
		t.Error("should assign second task second (FIFO)")
	}
}

func TestAssignTaskWorkerSelection(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker One", 1)
	tq.AddWorker("w2", "Worker Two", 1)

	tq.SubmitTask("Task 1", PriorityHigh, 3)
	tq.SubmitTask("Task 2", PriorityHigh, 3)
	tq.SubmitTask("Task 3", PriorityHigh, 3)

	// First two should be assigned
	tq.AssignTask()
	tq.AssignTask()

	// Both workers at capacity
	_, _, ok := tq.AssignTask()
	if ok {
		t.Error("should not assign when all workers at capacity")
	}

	// Inactive worker should not receive tasks
	tq.AddWorker("w3", "Worker Three", 5)
	tq.SetWorkerActive("w3", false)

	_, _, ok = tq.AssignTask()
	if ok {
		t.Error("should not assign to inactive worker")
	}
}

func TestAssignTaskNoWork(t *testing.T) {
	tq := NewTaskQueue()

	// No workers
	tq.SubmitTask("Task", PriorityHigh, 3)
	_, _, ok := tq.AssignTask()
	if ok {
		t.Error("should not assign with no workers")
	}

	// No tasks
	tq2 := NewTaskQueue()
	tq2.AddWorker("w1", "Worker", 5)
	_, _, ok = tq2.AssignTask()
	if ok {
		t.Error("should not assign with no tasks")
	}
}

func TestCompleteTask(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker One", 5)

	tq.SubmitTask("Task", PriorityHigh, 3)
	taskID, workerID, _ := tq.AssignTask()

	// Complete successfully
	if !tq.CompleteTask(taskID, workerID) {
		t.Error("should complete task")
	}

	task := tq.GetTask(taskID)
	if task.Status != StatusCompleted {
		t.Error("task should be completed")
	}
	if task.CompletedAt.IsZero() {
		t.Error("CompletedAt should be set")
	}
	if tq.GetWorker("w1").TaskCount != 0 {
		t.Error("worker TaskCount should be decremented")
	}

	// Complete non-existent task
	if tq.CompleteTask("TASK-999", "w1") {
		t.Error("non-existent task should fail")
	}

	// Complete with wrong worker
	tq.SubmitTask("Task 2", PriorityHigh, 3)
	tq.AddWorker("w2", "Worker Two", 5)
	taskID2, _, _ := tq.AssignTask()
	if tq.CompleteTask(taskID2, "w2") {
		t.Error("wrong worker should fail")
	}

	// Complete already completed task
	if tq.CompleteTask(taskID, workerID) {
		t.Error("completing completed task should fail")
	}
}

func TestFailTaskWithRetry(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker One", 5)

	tq.SubmitTask("Retry Task", PriorityHigh, 2) // max 2 retries
	taskID, workerID, _ := tq.AssignTask()

	// First failure - should retry
	if !tq.FailTask(taskID, workerID) {
		t.Error("should fail task")
	}

	task := tq.GetTask(taskID)
	if task.Status != StatusPending {
		t.Error("task should be back to pending for retry")
	}
	if task.RetryCount != 1 {
		t.Errorf("retry count should be 1, got %d", task.RetryCount)
	}
	if task.WorkerID != "" {
		t.Error("worker should be cleared for retry")
	}
	if tq.GetWorker("w1").TaskCount != 0 {
		t.Error("worker TaskCount should be decremented")
	}

	// Assign again
	taskID, workerID, _ = tq.AssignTask()
	tq.FailTask(taskID, workerID) // retry count = 2

	// Assign again
	taskID, workerID, _ = tq.AssignTask()
	tq.FailTask(taskID, workerID) // retry count = 2, maxRetries = 2, should fail permanently

	task = tq.GetTask(taskID)
	if task.Status != StatusFailed {
		t.Error("task should be permanently failed after max retries")
	}
	if task.RetryCount != 2 {
		t.Errorf("retry count should be 2, got %d", task.RetryCount)
	}
}

func TestFailTaskNoRetry(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker One", 5)

	tq.SubmitTask("No Retry Task", PriorityHigh, 0) // no retries
	taskID, workerID, _ := tq.AssignTask()

	tq.FailTask(taskID, workerID)

	task := tq.GetTask(taskID)
	if task.Status != StatusFailed {
		t.Error("task with 0 maxRetries should fail immediately")
	}
}

func TestFailTaskValidation(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker One", 5)

	// Non-existent task
	if tq.FailTask("TASK-999", "w1") {
		t.Error("non-existent task should fail")
	}

	// Wrong worker
	tq.SubmitTask("Task", PriorityHigh, 3)
	tq.AddWorker("w2", "Worker Two", 5)
	taskID, _, _ := tq.AssignTask()
	if tq.FailTask(taskID, "w2") {
		t.Error("wrong worker should fail")
	}

	// Pending task (not running)
	tq.SubmitTask("Pending Task", PriorityLow, 3)
	if tq.FailTask("TASK-2", "w1") {
		t.Error("failing pending task should fail")
	}
}

func TestGetPendingTasks(t *testing.T) {
	tq := NewTaskQueue()

	tq.SubmitTask("Low 1", PriorityLow, 3)
	time.Sleep(1 * time.Millisecond)
	tq.SubmitTask("High 1", PriorityHigh, 3)
	time.Sleep(1 * time.Millisecond)
	tq.SubmitTask("Medium 1", PriorityMedium, 3)
	time.Sleep(1 * time.Millisecond)
	tq.SubmitTask("High 2", PriorityHigh, 3)
	time.Sleep(1 * time.Millisecond)
	tq.SubmitTask("Low 2", PriorityLow, 3)

	pending := tq.GetPendingTasks()
	if len(pending) != 5 {
		t.Fatalf("expected 5 pending tasks, got %d", len(pending))
	}

	// Check order: High1, High2, Medium1, Low1, Low2
	expectedOrder := []string{"High 1", "High 2", "Medium 1", "Low 1", "Low 2"}
	for i, task := range pending {
		if task.Name != expectedOrder[i] {
			t.Errorf("position %d: expected %s, got %s", i, expectedOrder[i], task.Name)
		}
	}

	// Assign one task - should not appear in pending
	tq.AddWorker("w1", "Worker", 5)
	tq.AssignTask()

	pending = tq.GetPendingTasks()
	if len(pending) != 4 {
		t.Errorf("expected 4 pending after assignment, got %d", len(pending))
	}
}

func TestGetWorkerTasks(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker One", 5)

	tq.SubmitTask("Task 1", PriorityHigh, 3)
	tq.SubmitTask("Task 2", PriorityHigh, 3)
	tq.SubmitTask("Task 3", PriorityHigh, 3)

	tq.AssignTask()
	tq.AssignTask()

	tasks := tq.GetWorkerTasks("w1")
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks for worker, got %d", len(tasks))
	}

	// Non-existent worker
	tasks = tq.GetWorkerTasks("w999")
	if len(tasks) != 0 {
		t.Error("non-existent worker should return empty slice")
	}

	// Worker with no tasks
	tq.AddWorker("w2", "Worker Two", 5)
	tasks = tq.GetWorkerTasks("w2")
	if len(tasks) != 0 {
		t.Error("worker with no tasks should return empty slice")
	}
}

func TestGetQueueStats(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker", 5)

	tq.SubmitTask("Task 1", PriorityHigh, 3)
	tq.SubmitTask("Task 2", PriorityHigh, 3)
	tq.SubmitTask("Task 3", PriorityHigh, 0) // no retry

	pending, running, completed, failed := tq.GetQueueStats()
	if pending != 3 || running != 0 || completed != 0 || failed != 0 {
		t.Errorf("expected 3,0,0,0 got %d,%d,%d,%d", pending, running, completed, failed)
	}

	// Assign one
	taskID1, workerID1, _ := tq.AssignTask()
	pending, running, completed, failed = tq.GetQueueStats()
	if pending != 2 || running != 1 || completed != 0 || failed != 0 {
		t.Errorf("expected 2,1,0,0 got %d,%d,%d,%d", pending, running, completed, failed)
	}

	// Complete one
	tq.CompleteTask(taskID1, workerID1)
	pending, running, completed, failed = tq.GetQueueStats()
	if pending != 2 || running != 0 || completed != 1 || failed != 0 {
		t.Errorf("expected 2,0,1,0 got %d,%d,%d,%d", pending, running, completed, failed)
	}

	// Fail one permanently (no retry)
	taskID2, workerID2, _ := tq.AssignTask()
	taskID3, workerID3, _ := tq.AssignTask()
	tq.FailTask(taskID3, workerID3) // This is task 3 with 0 retries
	tq.CompleteTask(taskID2, workerID2)

	pending, running, completed, failed = tq.GetQueueStats()
	if pending != 0 || running != 0 || completed != 2 || failed != 1 {
		t.Errorf("expected 0,0,2,1 got %d,%d,%d,%d", pending, running, completed, failed)
	}
}

func TestReassignAbandonedTasks(t *testing.T) {
	tq := NewTaskQueue()
	tq.AddWorker("w1", "Worker One", 5)

	tq.SubmitTask("Old Task 1", PriorityHigh, 3)
	tq.SubmitTask("Old Task 2", PriorityHigh, 3)
	tq.SubmitTask("New Task", PriorityHigh, 3)

	// Assign all tasks
	id1, _, _ := tq.AssignTask()
	id2, _, _ := tq.AssignTask()
	id3, _, _ := tq.AssignTask()

	// Manually set old StartedAt for testing
	tq.GetTask(id1).StartedAt = time.Now().Add(-2 * time.Hour)
	tq.GetTask(id2).StartedAt = time.Now().Add(-90 * time.Minute)
	// id3 stays recent

	// Reassign tasks older than 1 hour
	count := tq.ReassignAbandonedTasks(1 * time.Hour)

	if count != 2 {
		t.Errorf("expected 2 reassigned, got %d", count)
	}

	if tq.GetTask(id1).Status != StatusPending {
		t.Error("old task 1 should be pending")
	}
	if tq.GetTask(id1).WorkerID != "" {
		t.Error("old task 1 worker should be cleared")
	}
	if tq.GetTask(id2).Status != StatusPending {
		t.Error("old task 2 should be pending")
	}
	if tq.GetTask(id3).Status != StatusRunning {
		t.Error("new task should still be running")
	}

	// Worker task count should be updated
	if tq.GetWorker("w1").TaskCount != 1 {
		t.Errorf("worker should have 1 task, got %d", tq.GetWorker("w1").TaskCount)
	}
}
