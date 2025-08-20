package utils

import (
	"time"

	"casino/entities"
)

type SaveDataManager interface {
	Save(data entities.SaveData)
	Read() entities.SaveData
}

func NewInMemorySaveDataManager() SaveDataManager {
	return &inMemorySaveDataManager{
		saveData: defaultSaveData(),
	}
}

type inMemorySaveDataManager struct {
	saveData entities.SaveData
}

func (b *inMemorySaveDataManager) Save(data entities.SaveData) {
	b.saveData = data
}

func (b *inMemorySaveDataManager) Read() entities.SaveData {
	return b.saveData
}

func getPreviousMidnight() time.Time {
	year, month, day := time.Now().UTC().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func defaultSaveData() entities.SaveData {
	return entities.SaveData{
		RemainingChips: 1000,
		LastResetAt:    getPreviousMidnight(),
	}
}
