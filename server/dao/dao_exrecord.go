package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type ExRecordDao interface {
	AddExchangeRecord(e *models.ExchangeRecord) error
	FindExchangeRecordByPhone(phoneNumber int) ([]*models.ExchangeRecord, error)
	FindAllExchangeRecord(pageNum, pageCount int) ([]*models.ExchangeRecord, error)
	FindAllExchangeCount() (int, error)
}
