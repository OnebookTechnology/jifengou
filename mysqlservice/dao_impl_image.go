package mysql

import "github.com/OnebookTechnology/Engineer/server/models"

func (m *MysqlService) AddImage(image *models.Image) error {
	tx, err := m.Db.Begin()
	_, err = tx.Exec("INSERT INTO image(product_id,image_name,image_type) VALUES(?,?,?)",
		image.ISBN, image.ImageURL, image.ImageType)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
