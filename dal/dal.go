package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"encoding/json"
	"beautyfarm4market/config"
)

const (
	dbhostsip  = "116.62.194.207" //IP地址
	dbusername = "root"           //用户名
	port       = 3306
	dbpassword = "Beautyfarm886633" //密码
)

var ViewLogList []ViewLog = []ViewLog{}
var dbconnection *sql.DB

func init() {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", dbusername, dbpassword, dbhostsip, port, config.ConfigInfo.Dbname)
	dbconnection, err = sql.Open("mysql", dataSourceName)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		//panic(err)
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
	Channel        string
	CreateDate     string
	ModifyDate     string
	WechatorderNo  string
	PayTime        string
	ClientIp       string
	OrignalPrice   float64
	ProductId      int64
}

//mappingOrder_no, product_code, mobile_no, user_name, account_no,
// total_price, order_status, pay_status, order_no, card_no, channel, create_date, modify_date
func AddTempOrder(t TempOrder) bool {
	//插入数据
	stmt, err := dbconnection.Prepare("INSERT temp_order SET mappingOrder_no=?,product_code=?,mobile_no=?," +
		"user_name=?,account_no=?,total_price=?,order_status=?,pay_status=?,channel=?,create_date=?,modify_date=?,product_name=?,client_ip=?,orignal_price=?,product_id=?")
	checkErr(err)

	res, err := stmt.Exec(t.MappingOrderNo, t.ProductCode, t.MobileNo, t.UserName, t.AccountNo,
		t.TotalPrice, t.OrderStatus, t.PayStatus, t.Channel, t.CreateDate, t.ModifyDate, t.ProductName, t.ClientIp, t.OrignalPrice, t.ProductId)
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
	stmt, err := dbconnection.Prepare("select *  FROM temp_order where mobile_no=? and product_code=? and pay_status=2 ")
	checkErr(err)
	rows, err := stmt.Query(mobile, productCode)
	var tempOrders []TempOrder = toTempOrder(rows)
	return tempOrders
}

func GetOrdersByMappingOrderNo(mappingOrderNo string) TempOrder {
	var tempOrder TempOrder = TempOrder{MappingOrderNo: "",}
	//查询数据
	stmt, err := dbconnection.Prepare("select *  FROM temp_order where mappingOrder_no=? ")
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
		var channel string
		var create_date string
		var modify_date string
		var client_ip string
		var orignal_price float64
		var product_id int64
		errScan := rows.Scan(&mappingOrder_no, &product_code, &product_name, &mobile_no, &user_name,
			&account_no, &total_price, &order_status,
			&pay_status, &order_no, &card_no, &wechatorder_no, &pay_time, &channel, &create_date, &modify_date, &client_ip, &orignal_price, &product_id)
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
			OrderStatus:    order_status,
			WechatorderNo:  wechatorder_no,
			PayTime:        pay_time,
			ClientIp:       client_ip,
			OrignalPrice:   orignal_price,
			ProductId:      product_id,
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
	sql := "INSERT log SET title=?,description=?,logType=?"
	res, err := dbconnection.Exec(sql,log.Title, log.Description, log.Type)
	checkErr(err)
	rows, _ := res.RowsAffected()
	return rows > 0
}

func AddJsonLog(title string, obj interface{}) bool {
	jstr, _ := json.Marshal(obj)
	return AddLog(LogInfo{Title: title, Description: string(jstr), Type: 1})
}

/******************************************/

/*prodcutinfo*/
type ProductInfo struct {
	Product_id         int64
	Prodcut_name       string
	Prodcut_desc       string
	Prodcut_rule       string
	Price              float64
	Orignal_price      float64
	Backgroud_image    string
	Rule_image         string
	PurhchaseBtn_image string
	Isactive           int
	Create_date        string
	Product_code       string
	MaskImage          string
}

func AddProductInfo(p ProductInfo) bool {
	//插入数据
	stmt, err := dbconnection.Prepare("INSERT product SET prodcut_name=?,prodcut_desc=?,prodcut_rule=?,price=?,orignal_price=?,backgroud_image=?,rule_image=?,purhchaseBtn_image=?,product_code=?,mask_image=?")
	checkErr(err)
	res, err := stmt.Exec(p.Prodcut_name, p.Prodcut_desc, p.Prodcut_rule, p.Price, p.Orignal_price, p.Backgroud_image, p.Rule_image, p.PurhchaseBtn_image, p.Product_code, p.MaskImage)
	checkErr(err)
	rows, _ := res.RowsAffected()
	return rows > 0
}

func UpdateProductInfo(p ProductInfo) bool {
	//插入数据
	stmt, err := dbconnection.Prepare("update product SET prodcut_name=?,prodcut_desc=?" +
		",prodcut_rule=?,price=?,orignal_price=?,backgroud_image=?,rule_image=?,purhchaseBtn_image=?,product_code=?,mask_image=?,isactive=? where product_id=?")
	checkErr(err)
	res, err := stmt.Exec(p.Prodcut_name, p.Prodcut_desc, p.Prodcut_rule, p.Price, p.Orignal_price, p.Backgroud_image,
		p.Rule_image, p.PurhchaseBtn_image, p.Product_code, p.MaskImage, p.Isactive, p.Product_id)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)
	return true
}

func GetProductInfo(productId int64, checkActive bool) ProductInfo {
	p := ProductInfo{Product_id: int64(0),}
	//查询数据
	sql := "select *  FROM product where Product_id=?"
	if checkActive {
		sql += " and isactive=1"
	}
	stmt, err := dbconnection.Prepare(sql)
	checkErr(err)
	rows, err := stmt.Query(productId)
	var products []ProductInfo = toProducts(rows)
	if len(products) > 0 {
		p = products[0]
	}
	return p
}

func GetAllProductInfos() []ProductInfo {
	//查询数据
	stmt, err := dbconnection.Prepare("select *  FROM product")
	checkErr(err)
	rows, err := stmt.Query()
	var products []ProductInfo = toProducts(rows)
	return products
}

func toProducts(rows *sql.Rows) []ProductInfo {
	var products []ProductInfo
	for rows.Next() {
		var product_id int64
		var prodcut_name string
		var prodcut_desc string
		var prodcut_rule string
		var price float64
		var orignal_price float64
		var backgroud_image string
		var rule_image string
		var purhchaseBtn_image string
		var isactive int
		var create_date string
		var product_code string
		var mask_image string
		errScan := rows.Scan(&product_id, &prodcut_name, &prodcut_desc, &prodcut_rule, &price,
			&orignal_price, &backgroud_image, &rule_image,
			&purhchaseBtn_image, &isactive, &create_date, &product_code, &mask_image)
		checkErr(errScan)
		p := ProductInfo{
			Product_id:         product_id,
			Prodcut_name:       prodcut_name,
			Prodcut_desc:       prodcut_desc,
			Prodcut_rule:       prodcut_rule,
			Price:              price,
			Orignal_price:      orignal_price,
			Backgroud_image:    backgroud_image,
			Rule_image:         rule_image,
			PurhchaseBtn_image: purhchaseBtn_image,
			Isactive:           isactive,
			Create_date:        create_date,
			Product_code:       product_code,
			MaskImage:          mask_image,
		}
		products = append(products, p)
	}
	return products
}

/*prodcutinfo end*/

/*view log*/
type ViewLog struct {
	Pageview_id  int64
	Pange_url    string
	Client_ip    string
	Channel_code string
}

func AddViewLog(v ViewLog) bool {
	//插入数据
	query := "INSERT page_view SET pange_url=?,client_ip=?,channel_code=?"
	res, err := dbconnection.Exec(query, v.Pange_url, v.Client_ip, v.Channel_code)
	checkErr(err)
	rows, _ := res.RowsAffected()
	return rows > 0
}

/*view log*/

/*流量统计*/
type ViewLogReportEntity struct {
	ChannelCode string
	TotalView   int
}

func GetViewLogReport() []ViewLogReportEntity {
	//查询数据
	rows, err := dbconnection.Query("SELECT channel_code,COUNT(1) as totalView from page_view GROUP BY channel_code ;")
	checkErr(err)
	var ViewLogReportArr []ViewLogReportEntity = toViewLogReport(rows)
	return ViewLogReportArr
}

func toViewLogReport(rows *sql.Rows) []ViewLogReportEntity {
	var viewLogReportArr []ViewLogReportEntity
	for rows.Next() {
		var channel_code string
		var totalView int
		errScan := rows.Scan(&channel_code, &totalView)
		checkErr(errScan)
		r := ViewLogReportEntity{
			ChannelCode: channel_code,
			TotalView:   totalView,
		}
		viewLogReportArr = append(viewLogReportArr, r)
	}
	return viewLogReportArr
}

/*流量统计*/

/*订单统计*/
type OrderReportEntity struct {
	Channel    string
	ProductId  int64
	Ordercount int
}

func GetOrderReport() []OrderReportEntity {
	//查询数据
	rows, err := dbconnection.Query("SELECT  channel,product_id,count(1) as ordercount from temp_order WHERE order_status=2 GROUP BY channel,product_id ;")
	checkErr(err)
	var orderReportArr []OrderReportEntity = toOrderReport(rows)
	return orderReportArr
}

func toOrderReport(rows *sql.Rows) []OrderReportEntity {
	var orderReportEntityArr []OrderReportEntity
	for rows.Next() {
		var channel string
		var product_id int64
		var ordercount int
		errScan := rows.Scan(&channel, &product_id, &ordercount)
		checkErr(errScan)
		r := OrderReportEntity{
			Channel:    channel,
			ProductId:  product_id,
			Ordercount: ordercount,
		}
		orderReportEntityArr = append(orderReportEntityArr, r)
	}
	return orderReportEntityArr
}

/*订单统计*/
