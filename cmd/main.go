package main

import (
	"HotelParse/internal/ui"
	"HotelParse/internal/xlsxutils"
	"log"
	"sync"
)

func main() {
	app := ui.NewApp()
	app.Run()
	var (
		infoArray []map[string]string
		mu        sync.Mutex
	)

	go func() {
		app.ProcessData(&infoArray, &mu)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
			}
			err := xlsxutils.ExportToExcel(infoArray)
			if err != nil {
				return
			}
		}()
	}()
}
