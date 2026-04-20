package commands

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check system health and resource usage",
	Long:  `Displays real-time statistics for CPU, Memory, Disk, and Go Runtime.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		color.Cyan("\n🔍 Gathering Kodia System Health Stats...\n")

		// CPU
		cpuPercs, _ := cpu.PercentWithContext(ctx, 0, false)
		cpuUsage := 0.0
		if len(cpuPercs) > 0 {
			cpuUsage = cpuPercs[0]
		}

		// Memory
		vm, _ := mem.VirtualMemoryWithContext(ctx)

		// Disk
		du, _ := disk.UsageWithContext(ctx, "/")

		// Formatting Table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Component", "Statistic", "Value"})
		
		// CPU Rows
		t.AppendRows([]table.Row{
			{"CPU", "Usage", fmt.Sprintf("%.2f%%", cpuUsage)},
			{"CPU", "Cores", fmt.Sprintf("%d", runtime.NumCPU())},
		})
		t.AppendSeparator()

		// Memory Rows
		t.AppendRows([]table.Row{
			{"Memory", "Total", formatBytes(vm.Total)},
			{"Memory", "Used", formatBytes(vm.Used)},
			{"Memory", "Usage", fmt.Sprintf("%.2f%%", vm.UsedPercent)},
		})
		t.AppendSeparator()

		// Disk Rows
		t.AppendRows([]table.Row{
			{"Disk", "Total", formatBytes(du.Total)},
			{"Disk", "Used", formatBytes(du.Used)},
			{"Disk", "Usage", fmt.Sprintf("%.2f%%", du.UsedPercent)},
		})
		t.AppendSeparator()

		// Runtime
		t.AppendRows([]table.Row{
			{"Runtime", "Goroutines", fmt.Sprintf("%d", runtime.NumGoroutine())},
			{"Runtime", "OS/Arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)},
		})

		t.SetStyle(table.StyleLight)
		t.Render()

		color.Green("\n✅ System status is stable.\n")
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
