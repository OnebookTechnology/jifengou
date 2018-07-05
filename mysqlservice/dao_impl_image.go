package mysql

import (
	"errors"
	"github.com/OnebookTechnology/jifengou/server/models"
)

func (m *MysqlService) AddImage(image *models.Image) (int64, error) {
	tx, err := m.Db.Begin()
	r, err := tx.Exec("INSERT INTO image(product_id,image_url,image_type,is_main,upload_time) VALUES(?,?,?,?,?)",
		image.ProductId, image.ImageUrl, image.ImageType, image.IsMain, image.UploadTime)
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return 0, rollBackErr
		}
		return 0, errors.New("Add Coupon err:" + err.Error())
	}
	id, err := r.LastInsertId()
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return 0, rollBackErr
		}
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return 0, rollBackErr
		}
		return 0, err
	}
	return id, nil
}
