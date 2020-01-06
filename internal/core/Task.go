package core

import (
	"context"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
)

// Task is one currently active action
type Task struct {
	id string

	path string

	playMutex   sync.Mutex
	playCommand *exec.Cmd
	playCancel  func()
}

// NewTask creates a new task.
func NewTask(id string) *Task {
	return &Task{id: id}
}

// Stop aborts all current actions of this task.
func (task *Task) Stop() {
	task.stopPlay()
}

// Update keeps the actions of this task running.
func (task *Task) Update(path string) {
	if task.path != path {
		task.stopPlay()
		task.path = path
	}
	if len(task.path) > 0 {
		task.updatePlay()
	}
}

func (task *Task) stopPlay() {
	task.playMutex.Lock()
	defer task.playMutex.Unlock()

	if task.playCancel != nil {
		task.playCancel()
		task.playCommand = nil
		task.playCancel = nil
	}
}

func (task *Task) updatePlay() {
	if task.playCommand == nil {
		audioFile := task.nextAudioFile()
		if len(audioFile) > 0 {
			ctx, cancel := context.WithCancel(context.Background())
			task.playCommand = exec.CommandContext(ctx, "play", "-q", audioFile)
			task.playCancel = cancel
			task.startPlay()
		}
	}
}

func (task *Task) startPlay() {
	cmd := task.playCommand
	go func() {
		_ = cmd.Run()
		task.onPlayStopped(cmd)
	}()
}

func (task *Task) onPlayStopped(cmd *exec.Cmd) {
	task.playMutex.Lock()
	defer task.playMutex.Unlock()

	if cmd == task.playCommand {
		task.playCommand = nil
		task.playCancel = nil
	}
}

func (task *Task) nextAudioFile() (result string) {
	audioFiles := task.allAudioFiles(task.path)
	count := len(audioFiles)

	if count > 0 {
		sort.Strings(audioFiles)
		result = audioFiles[rand.Intn(count)]
	}

	return
}

func (task *Task) allAudioFiles(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer func() { _ = file.Close() }()
	files, err := file.Readdir(0)
	if err != nil {
		return nil
	}

	var result []string
	for _, entry := range files {
		fullPath := filepath.Join(path, entry.Name())
		resolvedEntry, _ := os.Stat(fullPath)
		if resolvedEntry.IsDir() {
			result = append(result, task.allAudioFiles(fullPath)...)
		} else if task.isAudioFile(filepath.Ext(entry.Name())) {
			result = append(result, fullPath)
		}
	}

	return result
}

func (task *Task) isAudioFile(extension string) bool {
	return (extension == ".flac") || (extension == ".mp3")
}
