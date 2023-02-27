package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Cpu() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0
	} else {
		return percent[0]
	}
}

func (a *App) Disk() [2]uint64 {

	UsageStats, err := disk.Usage("/")
	var values [2]uint64

	if err != nil {
		values[0] = 0
		values[1] = 0
		return values
	}
	values[0] = UsageStats.Free
	values[1] = UsageStats.Used

	return values
}
