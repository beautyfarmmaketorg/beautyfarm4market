package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var dbconnection = &sql.DB{}

const (
	dbhostsip  = "116.62.194.207" //IP地址
	dbusername = "root"           //用户名
	port       = 3306
	dbpassword = "Beautyfarm886633"     //密码
	dbname     = "db_beautyfarm_market" //db
)

func init() {
	//root:123456@tcp(127.0.0.1:3306)/Test?charset=utf8
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", dbusername, dbpassword, dbhostsip, port, dbname)
	db, err := sql.Open("mysql", dataSourceName)
	checkErr(err)
	dbconnection = db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type TempOrder struct {
	MappingOrderNo string
	ProductCode    string
	ProductName    string
	MobileNo       string
	UserName       string
	AccountNo      string
	TotalPrice     float64
	OrderStatus    int
	PayStatus      int
	OrderNo        string
	CardNo         string
	Channel        int
	CreateDate     string
	ModifyDate     string
}

//mappingOrder_no, product_code, mobile_no, user_name, account_no,
// total_price, order_status, pay_status, order_no, card_no, channel, create_date, modify_date
func AddTempOrder(t TempOrder) bool {
	//插入数据
	stmt, err := dbconnection.Prepare("INSERT temp_order SET mappingOrder_no=?,product_code=?,mobile_no=?," +
		"user_name=?,account_no=?,total_price=?,order_status=?,pay_status=?,channel=?,create_date=?,modify_date=?,product_name=?")
	checkErr(err)

	res, err := stmt.Exec(t.MappingOrderNo, t.ProductCode, t.MobileNo, t.UserName, t.AccountNo,
		t.TotalPrice, t.OrderStatus, t.PayStatus, t.Channel, t.CreateDate, t.ModifyDate, t.ProductName)
	checkErr(err)
	rows, _ := res.RowsAffected()
	return rows > 0
}

func GetAllOrders() []TempOrder {
	//查询数据
	rows, err := dbconnection.Query("SELECT * FROM temp_order")
	checkErr(err)

	var tempOrders []TempOrder

	for rows.Next() {
		var mappingOrder_no string
		var product_code string
		var product_name string
		var mobile_no string
		var user_name string
		var account_no string
		var total_price float64
		var order_status int
		var pay_status int
		var order_no string
		var card_no string
		var channel int
		var create_date string
		var modify_date string

		err = rows.Scan(&mappingOrder_no, &product_code, &product_name, &mobile_no, &user_name,
			&account_no, &total_price, &order_status,
			&pay_status, &order_no, &card_no, &channel, &create_date, &modify_date)

		t := TempOrder{
			MappingOrderNo: mappingOrder_no,
			UserName:       user_name,
			MobileNo:       mobile_no,
			ProductCode:    product_code,
			AccountNo:      account_no,
			Channel:        channel,
			CreateDate:     create_date,
			ModifyDate:     modify_date,
			TotalPrice:     total_price,
			ProductName:    product_name,
			OrderNo:        order_no,
			CardNo:         card_no,
		}
		tempOrders = append(tempOrders, t)
	}
	return tempOrders
}

func GetOrdersByMobile(mobile string, productCode string) []TempOrder {
	//查询数据
	stmt, err := dbconnection.Prepare("select *  FROM db_beautyfarm_market.temp_order where mobile_no=? and product_code=?")
	checkErr(err)
	rows, err := stmt.Query(mobile, productCode)
	var tempOrders []TempOrder

	for rows.Next() {
		var mappingOrder_no string
		var product_code string
		var product_name string
		var mobile_no string
		var user_name string
		var account_no string
		var total_price float64
		var order_status int
		var pay_status int
		var order_no string
		var card_no string
		var channel int
		var create_date string
		var modify_date string

		err = rows.Scan(&mappingOrder_no, &product_code, &product_name, &mobile_no, &user_name,
			&account_no, &total_price, &order_status,
			&pay_status, &order_no, &card_no, &channel, &create_date, &modify_date)

		t := TempOrder{
			MappingOrderNo: mappingOrder_no,
			UserName:       user_name,
			MobileNo:       mobile_no,
			ProductCode:    product_code,
			AccountNo:      account_no,
			Channel:        channel,
			CreateDate:     create_date,
			ModifyDate:     modify_date,
			TotalPrice:     total_price,
			ProductName:    product_name,
			OrderNo:        order_no,
			CardNo:         card_no,
		}
		tempOrders = append(tempOrders, t)
	}
	return tempOrders

}
