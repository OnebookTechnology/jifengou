package mysql

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"strings"
)

// 添加记录
func (m *MysqlService) AddExchangeRecord(e *models.ExchangeRecord) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	// s1. update online book's last_op_time、last_op_phone_number、online_status
	_, err = tx.Exec("INSERT INTO ex_record(phone_number,b_codes,p_code,ex_time,p_id) VALUES (?,?,?,?,?)",
		e.PhoneNumber, e.BCodes, e.PCode, e.ExTime, e.PId)
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return err
	}
	return nil
}

// 查询记录
func (m *MysqlService) FindExchangeRecordByPhone(phoneNumber int) ([]*models.ExchangeRecord, error) {
	rows, err := m.Db.Query("SELECT p.product_name, e.b_codes, e.p_code, e.ex_time, e.p_id FROM ex_record e "+
		"LEFT JOIN product p ON p.product_id=e.p_id where e.phone_number=?", phoneNumber)
	if err != nil {
		return nil, err
	}
	var es []*models.ExchangeRecord
	for rows.Next() {
		e := new(models.ExchangeRecord)
		err := rows.Scan(&e.Name, &e.BCodes, &e.PCode, &e.ExTime, &e.PId)
		if err != nil {
			return nil, err
		}
		e.BCodeArray = strings.Split(e.BCodes, ",")
		es = append(es, e)
	}
	return es, nil
}
