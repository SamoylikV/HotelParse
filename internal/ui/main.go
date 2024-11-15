package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/SamoylikV/HotelParse/internal/parse"
	"github.com/SamoylikV/HotelParse/internal/xlsxutils"
	"math"
	"sync"
	"time"
)

type MyApp struct {
	startID         int
	N               int
	logs            string
	progress        float64
	startChan       chan bool
	logText         *widget.Label
	progressBar     *widget.ProgressBar
	sliderLabel     *widget.Label
	slider          *widget.Slider
	startTime       time.Time
	lastUpdatedTime time.Time
	lastCompleted   int
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
		app_.startTime = time.Now()
		app_.lastUpdatedTime = app_.startTime
		app_.lastCompleted = 0
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

	elapsed := time.Since(app_.startTime).Seconds()
	var remainingTime string

	if elapsed > 0 && completed > 0 {
		averageSpeed := float64(completed) / elapsed
		remaining := float64(total-completed) / averageSpeed

		if remaining > 60 {
			minutes := int(remaining) / 60
			seconds := int(remaining) % 60
			remainingTime = fmt.Sprintf("⏳ Осталось: %d мин. %d сек.", minutes, seconds)
		} else {
			remainingTime = fmt.Sprintf("⏳ Осталось: %.0f сек.", remaining)
		}
	} else {
		remainingTime = "⏳ Расчет времени невозможен."
	}

	app_.logs = fmt.Sprintf("✅ Прогресс: %.2f%%\n%s\n", app_.progress*100, remainingTime)
	app_.progressBar.SetValue(app_.progress)
	app_.logText.SetText(app_.logs)
}

func (app_ *MyApp) startProcessing() {
	app_.startChan = make(chan bool)
	var infoArray []map[string]string
	var mu sync.Mutex

	defer func() {
		if r := recover(); r != nil {
			app_.logText.SetText(fmt.Sprintf("🔥 Ошибка: %v. Сохраняем собранные данные.", r))
			err := xlsxutils.ExportToExcel(infoArray)
			if err != nil {
				app_.logText.SetText("❌ Не удалось сохранить частичные данные в Excel.")
			} else {
				app_.logText.SetText("✅ Частичные данные сохранены в Excel.")
			}
			app_.slider.Enable()
		}
	}()

	go func() {
		app_.ProcessData(&infoArray, &mu)
		go func() {
			err := xlsxutils.ExportToExcel(infoArray)
			if err != nil {
				app_.logText.SetText("❌ Не удалось экспортировать в Excel.")
			} else {
				app_.logText.SetText("✅ Данные успешно экспортированы в Excel.")
			}
			app_.slider.Enable()
		}()
	}()
}

func (app_ *MyApp) ProcessData(infoArray *[]map[string]string, mu *sync.Mutex) {
	var wg sync.WaitGroup

	numWorkers := 50
	dataChan := make(chan int, 1000)
	progressChanSize := int(math.Min(float64(app_.N)/50, 50))

	progressChan := make(chan int, progressChanSize)

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

	batchSize := 1000
	for i := app_.startID; i < app_.startID+app_.N; i += batchSize {
		for j := i; j < i+batchSize && j < app_.startID+app_.N; j++ {
			dataChan <- j
		}
	}
	close(dataChan)

	wg.Wait()
	close(progressChan)
}
