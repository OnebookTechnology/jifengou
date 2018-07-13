package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type ExRecordDao interface {
	AddExchangeRecord(e *models.ExchangeRecord) error
	FindExchangeRecordByPhone(phoneNumber int) ([]*models.ExchangeRecord, error)
}
