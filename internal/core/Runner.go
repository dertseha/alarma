package core

import (
	"math/rand"
	"time"

	"github.com/dertseha/alarma/internal/config"
)

// Runner is the main logic of the alarma application.
type Runner struct {
	activeTasks map[string]*Task
}

// NewRunner starts a new runner instance
func NewRunner() *Runner {
	rand.Seed(time.Now().UnixNano())

	runner := &Runner{
		activeTasks: make(map[string]*Task),
	}

	return runner
}

// Update iterates over all active tasks and updates their state.
func (runner *Runner) Update(configuration config.Instance) {
	currentTasks := make(map[string]config.TimeSpan)
	now := time.Now()

	if configuration.TimeSpansActive {
		for _, timeSpan := range configuration.TimeSpans {
			if runner.shouldBeCurrent(timeSpan, now) {
				currentTasks[timeSpan.ID] = timeSpan
			}
		}
	}

	var tasksToStop []string
	for id := range runner.activeTasks {
		if _, isCurrent := currentTasks[id]; !isCurrent {
			tasksToStop = append(tasksToStop, id)
		}
	}
	runner.stopTasks(tasksToStop)

	for id, timeSpan := range currentTasks {
		if _, isActive := runner.activeTasks[id]; !isActive {
			runner.startTask(id)
		}
		runner.updateTask(id, timeSpan)
	}
}

func (runner *Runner) shouldBeCurrent(timeSpan config.TimeSpan, now time.Time) bool {
	reference := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	from := runner.parseTime(reference, timeSpan.From)
	to := runner.parseTime(reference, timeSpan.To)

	return timeSpan.Enabled && from.Before(now) && now.Before(to)
}

func (runner *Runner) parseTime(reference time.Time, timeText string) time.Time {
	t, _ := time.ParseInLocation("15:04", timeText, time.UTC)
	return time.Date(reference.Year(), reference.Month(), reference.Day(), t.Hour(), t.Minute(), 0, 0, reference.Location())
}

func (runner *Runner) stopTasks(ids []string) {
	for _, id := range ids {
		runner.stopTask(id)
	}
}

func (runner *Runner) stopTask(id string) {
	task := runner.activeTasks[id]
	delete(runner.activeTasks, id)
	task.Stop()
}

func (runner *Runner) startTask(id string) {
	runner.activeTasks[id] = NewTask(id)
}

func (runner *Runner) updateTask(id string, configuration config.TimeSpan) {
	runner.activeTasks[id].Update(configuration.Path)
}
