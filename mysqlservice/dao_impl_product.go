package mysql

import (
	"database/sql"
	"errors"
	"github.com/OnebookTechnology/jifengou/server/models"
	"time"
)

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

// 根据ItemStatement查询商品
func (m *MysqlService) FindProductByItemStatement(itemStatement string) (p *models.Product, err error) {
	row := m.Db.QueryRow("SELECT product_name, product_info FROM product WHERE product_item_statement = ?", itemStatement)
	p = new(models.Product)
	err = row.Scan(&p.ProductName, &p.ProductInfo)
	return
}

// 查找所有商品
func (m *MysqlService) FindAllProducts() ([]*models.Product, error) {
	rows, err := m.Db.Query("SELECT product_id, product_item_statement,product_name, product_status, product_item_statement, product_price FROM product")
	if err != nil {
		return nil, nil
	}
	var products []*models.Product
	for rows.Next() {
		p := new(models.Product)
		err = rows.Scan(&p.ProductId, &p.ProductItemStatement, &p.ProductName, &p.ProductStatus, &p.ProductItemStatement, &p.ProductPrice)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// 查找商家的所有商品
func (m *MysqlService) FindAllProductByBusinessIdAndStatus(businessId int, status int, pageNum, pageCount int) ([]*models.Product, error) {
	var rows *sql.Rows
	var err error
	sqlStr := "SELECT product_id, product_item_statement, product_name, product_status FROM product WHERE business_id=? "
	if status != -1 {
		//All
		sqlStr += " AND product_status=?"
		sqlStr += " ORDER BY product_online_time DESC LIMIT ?,?"
		rows, err = m.Db.Query(sqlStr, businessId, status, (pageNum-1)*pageCount, pageCount)
	} else {
		sqlStr += " ORDER BY product_online_time DESC LIMIT ?,?"
		rows, err = m.Db.Query(sqlStr, businessId, (pageNum-1)*pageCount, pageCount)
	}

	if err != nil {
		return nil, err
	}
	var products []*models.Product
	for rows.Next() {
		p := new(models.Product)
		err = rows.Scan(&p.ProductId, &p.ProductItemStatement, &p.ProductName, &p.ProductInfo)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// 添加商品
func (m *MysqlService) AddProduct(p *models.Product) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	// s1. update online book's last_op_time、last_op_phone_number、online_status
	_, err = tx.Exec("INSERT INTO product(product_item_statement, product_name, product_info,product_status,business_id,product_category,"+
		"product_subtitle,product_price,product_start_time,product_end_time,product_alert_count,product_online_time) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
		p.ProductItemStatement, p.ProductName, p.ProductInfo, models.ProductReviewing, p.BusinessId, p.ProductCategory,
		p.ProductSubtitle, p.ProductPrice, p.ProductStartTime, p.ProductEndTime, p.ProductAlertCount, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return errors.New("AddProduct err:" + err.Error())
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

// 确认商品状态
