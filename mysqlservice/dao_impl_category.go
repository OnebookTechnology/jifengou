package mysql

import "github.com/OnebookTechnology/jifengou/server/models"

// 查找所有商户
func (m *MysqlService) FindAllCategory() ([]*models.Category, error) {
	rows, err := m.Db.Query("SELECT category_id,category_name FROM category ")
	if err != nil {
		return nil, err
	}
	var cs []*models.Category
	for rows.Next() {
		c := new(models.Category)
		err := rows.Scan(&c.CategoryId, &c.CategoryName)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}
