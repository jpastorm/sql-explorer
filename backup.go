package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

var (
	tableMutex     sync.Mutex
	tableBuffer    table.Model
	tableFilePath  = ".tableBackup"
	tableFlushChan = make(chan bool, 1)
)

const (
	flushInterval = 2 * time.Second
)

func InitBackupSystems() {
	go tableBackupWriter()

	loadInitialBackups()
}

func loadInitialBackups() {
	go func() {
		if tbl, err := loadTableBackup(); err == nil {
			tableMutex.Lock()
			tableBuffer = tbl
			tableMutex.Unlock()
		}
	}()
}

func CloseBackupSystems() {
	flushTableBackup()
	close(tableFlushChan)
}

func tableBackupWriter() {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			flushTableBackup()
		case <-tableFlushChan:
			flushTableBackup()
			return
		}
	}
}

func flushTableBackup() {
	tableMutex.Lock()
	defer tableMutex.Unlock()

	if len(tableBuffer.Rows()) == 0 {
		return
	}

	file, err := os.OpenFile(tableFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	data := struct {
		Columns []table.Column `json:"columns"`
		Rows    []table.Row    `json:"rows"`
	}{
		Columns: tableBuffer.Columns(),
		Rows:    tableBuffer.Rows(),
	}

	encoder := json.NewEncoder(file)
	_ = encoder.Encode(data)
}

func loadTableBackup() (table.Model, error) {
	file, err := os.Open(tableFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return table.Model{}, nil
		}
		return table.Model{}, err
	}
	defer file.Close()

	var data struct {
		Columns []table.Column `json:"columns"`
		Rows    []table.Row    `json:"rows"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return table.Model{}, err
	}

	tbl := table.New(
		table.WithColumns(data.Columns),
		table.WithRows(data.Rows),
	)

	return tbl, nil
}

func SaveTableState(tbl table.Model) {
	tableMutex.Lock()
	tableBuffer = tbl
	tableMutex.Unlock()

	go flushTableBackup()
}

func RestoreTable() (table.Model, error) {
	return loadTableBackup()
}

func saveEditorBackup(content string) error {
	err := os.WriteFile("./.editorBackup", []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("Err writefile backup: %v", err)
	}

	return nil
}

func loadEditorBackup() (string, error) {
	content, err := os.ReadFile("./.editorBackup")
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("error reading file backup: %v", err)
	}

	return string(content), nil
}
