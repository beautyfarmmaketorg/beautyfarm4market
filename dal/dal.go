package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	dbhostsip  = "116.62.194.207" //IP地址
	dbusername = "root"           //用户名
	port       = 3306
	dbpassword = "Beautyfarm886633"     //密码
	dbname     = "db_beautyfarm_market" //db
)

var dbconnection *sql.DB

func init() {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", dbusername, dbpassword, dbhostsip, port, dbname)
	dbconnection, err =  sql.Open("mysql", dataSourceName)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/****************** 订单************************/

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
	WechatorderNo  string
	PayTime        string
	ClientIp string
}

//mappingOrder_no, product_code, mobile_no, user_name, account_no,
// total_price, order_status, pay_status, order_no, card_no, channel, create_date, modify_date
func AddTempOrder(t TempOrder) bool {
	//插入数据
	stmt, err := dbconnection.Prepare("INSERT temp_order SET mappingOrder_no=?,product_code=?,mobile_no=?," +
		"user_name=?,account_no=?,total_price=?,order_status=?,pay_status=?,channel=?,create_date=?,modify_date=?,product_name=?,client_ip=?")
	checkErr(err)

	res, err := stmt.Exec(t.MappingOrderNo, t.ProductCode, t.MobileNo, t.UserName, t.AccountNo,
		t.TotalPrice, t.OrderStatus, t.PayStatus, t.Channel, t.CreateDate, t.ModifyDate, t.ProductName,t.ClientIp)
	checkErr(err)
	rows, _ := res.RowsAffected()
	return rows > 0
}

func UpdateTempOrder(cardNo string, orderNo string, mappingOrderNo string, wechatOrderNo string, timeEnd string) bool {
	//插入数据
	stmt, err := dbconnection.Prepare("UPDATE temp_order SET card_no=?,order_no=?,wechatorder_no=?,pay_time=?,modify_date=NOW(),order_status=2,pay_status=2 where mappingOrder_no=?")
	checkErr(err)
	res, err := stmt.Exec(cardNo, orderNo, wechatOrderNo, timeEnd, mappingOrderNo)
	checkErr(err)
	rows, _ := res.RowsAffected()
	return rows > 0
}

//更新支付中状态
func UpdateTempOrderPayStatus(mappingOrderNo string, payStatus int) bool {
	//插入数据
	stmt, err := dbconnection.Prepare("UPDATE temp_order SET modify_date=NOW(),pay_status=? where mappingOrder_no=?")
	checkErr(err)
	res, err := stmt.Exec(payStatus, mappingOrderNo)
	checkErr(err)
	rows, _ := res.RowsAffected()
	return rows > 0
}

func GetAllOrders() []TempOrder {
	//查询数据
	rows, err := dbconnection.Query("SELECT * FROM temp_order")
	checkErr(err)
	var tempOrders []TempOrder = toTempOrder(rows)
	return tempOrders
}

func GetOrdersByMobile(mobile string, productCode string) []TempOrder {
	//查询数据
	stmt, err := dbconnection.Prepare("select *  FROM db_beautyfarm_market.temp_order where mobile_no=? and product_code=?")
	checkErr(err)
	rows, err := stmt.Query(mobile, productCode)
	var tempOrders []TempOrder = toTempOrder(rows)
	return tempOrders
}

func GetOrdersByMappingOrderNo(mappingOrderNo string) TempOrder {
	var tempOrder TempOrder = TempOrder{MappingOrderNo: "",}
	//查询数据
	stmt, err := dbconnection.Prepare("select *  FROM db_beautyfarm_market.temp_order where mappingOrder_no=? ")
	checkErr(err)
	rows, err := stmt.Query(mappingOrderNo)
	var tempOrders []TempOrder = toTempOrder(rows)
	if len(tempOrders) > 0 {
		tempOrder = tempOrders[0]
	}
	return tempOrder
}

func toTempOrder(rows *sql.Rows) []TempOrder {
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
		var wechatorder_no string
		var pay_time string
		var channel int
		var create_date string
		var modify_date string
		var client_ip string
		errScan := rows.Scan(&mappingOrder_no, &product_code, &product_name, &mobile_no, &user_name,
			&account_no, &total_price, &order_status,
			&pay_status, &order_no, &card_no, &wechatorder_no, &pay_time, &channel, &create_date, &modify_date,&client_ip)
		checkErr(errScan)
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
			PayStatus:      pay_status,
			OrderStatus:order_status,
			WechatorderNo:wechatorder_no,
			PayTime:pay_time,
			ClientIp:client_ip,
		}
		tempOrders = append(tempOrders, t)
	}
	return tempOrders
}

/****************** 订单END************************/

/****************** 日志************************/
type LogInfo struct {
	LogId       int64
	Title       string
	Description string
	Type        int
	CreateDate  time.Time
}

func AddLog(log LogInfo) bool {
	//插入数据
	stmt, err := dbconnection.Prepare("INSERT log SET title=?,description=?,logType=?")
	checkErr(err)
	res, err := stmt.Exec(log.Title, log.Description, log.Type)
	checkErr(err)
	rows, _ := res.RowsAffected()
	return rows > 0
}

/******************************************/
