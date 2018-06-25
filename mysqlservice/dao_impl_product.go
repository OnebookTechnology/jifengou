package mysql

import "github.com/OnebookTechnology/jifengou/server/models"

// 根据Id查找商品
func (m *MysqlService) FindProductById(productId string) (*models.Product, error) {
	row := m.Db.QueryRow("SELECT product_name, product_info FROM product WHERE product_id=?",
		productId)
	p := new(models.Product)
	err := row.Scan(&p.ProductName, &p.ProductInfo)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// 查找所有商品
func (m *MysqlService) FindAllProducts() ([]*models.Product, error) {
	rows, err := m.Db.Query("SELECT product_name, product_info, product_item_statement, product_price FROM product")
	if err != nil {
		return nil, nil
	}
	var products []*models.Product
	for rows.Next() {
		p := new(models.Product)
		err = rows.Scan(&p.ProductName, &p.ProductInfo, &p.ProductItemStatement, &p.ProductPrice)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// 查找商家的所有商品
func (m *MysqlService) FindAllProductByBusinessId(businessId uint) ([]*models.Product, error) {
	rows, err := m.Db.Query("SELECT product_name, product_info FROM product WHERE business_id=? ", businessId)
	if err != nil {
		return nil, nil
	}
	var products []*models.Product
	for rows.Next() {
		p := new(models.Product)
		err = rows.Scan(&p.ProductName, &p.ProductInfo)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// 添加商品

// 确认商品状态
