package mysql

import "github.com/OnebookTechnology/jifengou/server/models"

// 查找所有商户
func (m *MysqlService) FindAllCategory() ([]models.Category, error) {
	rows, err := m.Db.Query("SELECT category_id,category_name FROM category ")
	if err != nil {
		return nil, nil
	}
	var cs []models.Category
	for rows.Next() {
		c := new(models.Category)
		err := rows.Scan(&b.BusinessId, &b.BusinessNo, &b.BusinessName, &b.BusinessRegisterTime, &b.BusinessAuth, &avail)
		if err != nil {
			return nil, err
		}
		cs = append(bs, *b)
	}
	return bs, nil
}