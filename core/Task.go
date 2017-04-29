package core

import (
	"context"
	//"fmt"
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
	//fmt.Printf("[%v] NEW\n", id)
	return &Task{id: id}
}

// Stop aborts all current actions of this task.
func (task *Task) Stop() {
	//fmt.Printf("[%v] STOP\n", task.id)
	task.stopPlay()
}

// Update keeps the actions of this task running.
func (task *Task) Update(path string) {
	if task.path != path {
		//fmt.Printf("[%v] UPDATE: <%v>\n", task.id, path)

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
		//fmt.Printf("[%v] CANCEL\n", task.id)
		task.playCancel()
		task.playCommand = nil
		task.playCancel = nil
	}
}

func (task *Task) updatePlay() {
	if task.playCommand == nil {
		audioFile := task.nextAudioFile()
		if len(audioFile) > 0 {
			//fmt.Printf("[%v] PLAY: <%v>\n", task.id, audioFile)
			ctx, cancel := context.WithCancel(context.Background())
			task.playCommand = exec.CommandContext(ctx, "play", audioFile)
			task.playCancel = cancel
			task.startPlay()
		}
	}
}

func (task *Task) startPlay() {
	cmd := task.playCommand
	go func() {
		cmd.Run()
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
	result := []string{}
	file, _ := os.Open(path)

	if file != nil {
		defer file.Close()
		files, _ := file.Readdir(0)

		for _, entry := range files {
			fullPath := filepath.Join(path, entry.Name())
			if entry.IsDir() {
				result = append(result, task.allAudioFiles(fullPath)...)
			} else if task.isAudioFile(filepath.Ext(entry.Name())) {
				result = append(result, fullPath)
			}
		}
	}

	return result
}

func (task *Task) isAudioFile(extension string) bool {
	return (extension == ".flac") || (extension == ".mp3")
}
