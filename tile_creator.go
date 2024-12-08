package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	TileURLTemplate = "https://services.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/%d/%d/%d.png"
	MaxRetries      = 3
	WorkerPoolSize  = 64
	AvgTileSizeKB   = 50 // Average tile size in KB for estimation
)

type Task struct {
	Zoom int
	X    int
	Y    int
}

func degToNum(lat, lon float64, zoom int) (int, int) {
	n := math.Pow(2, float64(zoom))
	x := int((lon + 180.0) / 360.0 * n)
	y := int((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * n)
	return x, y
}

func downloadTile(task Task, outputDir string, progress *progressbar.ProgressBar, totalSize *uint64) error {
	url := fmt.Sprintf(TileURLTemplate, task.Zoom, task.Y, task.X)
	savePath := filepath.Join(outputDir, strconv.Itoa(task.Zoom), strconv.Itoa(task.X), fmt.Sprintf("%d.png", task.Y))

	if _, err := os.Stat(savePath); err == nil {
		progress.Add(1)
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		return err
	}

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			file, err := os.Create(savePath)
			if err != nil {
				resp.Body.Close()
				return err
			}
			defer file.Close()

			n, err := io.Copy(file, resp.Body)
			resp.Body.Close()
			if err == nil {
				atomic.AddUint64(totalSize, uint64(n))
				progress.Add(1)
				return nil
			}
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to download tile: Zoom %d, X %d, Y %d", task.Zoom, task.X, task.Y)
}

func main() {
	maxZoom := flag.Int("zoom", 5, "Highest zoom level (e.g., 5 means zoom levels 0 to 5)")
	minLat := flag.Float64("minlat", 16.5, "Minimum latitude")
	maxLat := flag.Float64("maxlat", 26.5, "Maximum latitude")
	minLon := flag.Float64("minlon", 51.8, "Minimum longitude")
	maxLon := flag.Float64("maxlon", 60.0, "Maximum longitude")
	outputDir := flag.String("output", "./tiles", "Output directory for downloaded tiles")

	flag.Parse()

	zoomLevels := make([]int, *maxZoom+1)
	for i := 0; i <= *maxZoom; i++ {
		zoomLevels[i] = i
	}

	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	var tasks []Task
	for _, zoom := range zoomLevels {
		xMin, yMin := degToNum(*maxLat, *minLon, zoom)
		xMax, yMax := degToNum(*minLat, *maxLon, zoom)

		maxTile := 1 << zoom
		if xMin < 0 {
			xMin = 0
		}
		if xMax >= maxTile {
			xMax = maxTile - 1
		}
		if yMin < 0 {
			yMin = 0
		}
		if yMax >= maxTile {
			yMax = maxTile - 1
		}

		for x := xMin; x <= xMax; x++ {
			for y := yMin; y <= yMax; y++ {
				tasks = append(tasks, Task{Zoom: zoom, X: x, Y: y})
			}
		}
	}

	progress := progressbar.NewOptions(len(tasks))
	var totalSize uint64

	jobs := make(chan Task, len(tasks))
	var wg sync.WaitGroup

	for i := 0; i < WorkerPoolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range jobs {
				_ = downloadTile(task, *outputDir, progress, &totalSize)
			}
		}()
	}

	for _, task := range tasks {
		jobs <- task
	}
	close(jobs)
	wg.Wait()

	log.Printf("Download complete. Total size: %.2f MB", float64(totalSize)/1024/1024)
}
