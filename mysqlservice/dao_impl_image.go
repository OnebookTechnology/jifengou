package mysql

import "github.com/OnebookTechnology/jifengou/server/models"

func (m *MysqlService) AddImage(image *models.Image) error {
	tx, err := m.Db.Begin()
	_, err = tx.Exec("INSERT INTO image(product_id,image_name,image_type) VALUES(?,?,?)",
		image.ProductId, image.ImageName, image.ImageType)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
