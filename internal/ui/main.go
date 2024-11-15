package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/SamoylikV/HotelParse/internal/parse"
	"github.com/SamoylikV/HotelParse/internal/xlsxutils"
	"sync"
)

type MyApp struct {
	startID     int
	N           int
	logs        string
	progress    float64
	startChan   chan bool
	logText     *widget.Label
	progressBar *widget.ProgressBar
	sliderLabel *widget.Label
	slider      *widget.Slider
}

func NewApp() *MyApp {
	return &MyApp{
		startID:   550000000,
		N:         100,
		startChan: make(chan bool),
	}
}

func (app_ *MyApp) Run() {
	myApp := app.New()
	app_.logText = widget.NewLabel("")
	app_.sliderLabel = widget.NewLabel(fmt.Sprintf("Будет проверено %d номеров", app_.N))
	app_.slider = widget.NewSlider(100, 100000)
	app_.slider.SetValue(float64(app_.N))
	app_.slider.OnChanged = func(value float64) {
		app_.N = int(value)
		app_.sliderLabel.SetText(fmt.Sprintf("Будет проверено %d номеров", app_.N))
	}

	startButton := widget.NewButton("Start", func() {
		app_.slider.Disable()
		go app_.startProcessing()
	})
	app_.progressBar = widget.NewProgressBar()
	content := container.NewVBox(app_.slider, app_.sliderLabel, startButton, app_.progressBar, app_.logText)
	window := myApp.NewWindow("Hotel parse")
	window.Resize(fyne.NewSize(400, 200))
	window.SetContent(content)
	window.ShowAndRun()
}

func (app_ *MyApp) UpdateProgress(completed, total int) {
	app_.progress = float64(completed) / float64(total)
	app_.logs = fmt.Sprintf("Progress: %.2f%%\n", app_.progress*100)
	app_.progressBar.SetValue(app_.progress)
}

func (app_ *MyApp) startProcessing() {
	app_.startChan = make(chan bool)
	var infoArray []map[string]string
	var mu sync.Mutex
	go func() {
		app_.ProcessData(&infoArray, &mu)
		go func() {
			err := xlsxutils.ExportToExcel(infoArray)
			if err != nil {
				app_.logText.SetText("Не удалось экспортировать в Excel.")
			} else {
				app_.logText.SetText("Данные успешно экспортированы в Excel.")
			}
			app_.slider.Enable()
		}()
	}()
}

func (app_ *MyApp) ProcessData(infoArray *[]map[string]string, mu *sync.Mutex) {
	var wg sync.WaitGroup
	progressChan := make(chan int, app_.N)
	dataChan := make(chan int, 100)
	numWorkers := 50

	go func() {
		completed := 0
		for range progressChan {
			completed++
			app_.UpdateProgress(completed, app_.N)
		}
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for num := range dataChan {
				id, err := parse.Id(num)
				if err != nil {
					progressChan <- 1
					continue
				}
				dict, err := parse.Info(id)
				if err != nil {
					progressChan <- 1
					continue
				}
				mu.Lock()
				*infoArray = append(*infoArray, dict)
				mu.Unlock()
				progressChan <- 1
			}
		}()
	}

	for i := app_.startID; i < app_.startID+app_.N; i++ {
		dataChan <- i
	}
	close(dataChan)

	wg.Wait()
	close(progressChan)
}
