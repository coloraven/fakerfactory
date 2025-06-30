package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/toddlerya/fakerfactory/faker"
)

var port string
var dbPath string
var Conn *sql.DB


func init() {
	if len(os.Args) > 1 {

		port = os.Args[1]
	} else {
		port = "8001"
	}
	if len(os.Args) > 2 {
		dbPath = os.Args[2]
	} else {
		dbPath = `./data/data.db`
	}
	fmt.Println("args==>", os.Args)
	fmt.Println("args[1]==>", os.Args[1])
	fmt.Println("args[2]==>", os.Args[2])
	Conn = faker.CreateConn(dbPath)
}

func StartServer() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	// TODO 后续投入生产要考虑日志分割，日志大小等问题
	f, _ := os.Create("serve.log")

	// Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()
	router.Use(cors.Default()) // 允许任何服务ajax跨域调用
	v1 := router.Group("api/v1")
	{
		v1.GET("/fakerfactory", GetFaker)
		v1.POST("/fakerfactory", PostFaker) // 添加这一行
	}
	err := router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("在%s端口启动服务失败！", port)
	}
}

func PostFaker(c *gin.Context) {
	// 定义接收结构体
	type RequestBody struct {
		Columns string `json:"columns" binding:"required"`
		Number  string `json:"number" binding:"required"`
	}

	var req RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": gin.H{
				"status": "error",
				"code":   "100"},
			"data": gin.H{
				"count":   nil,
				"records": "请输入有效的参数"},
		})
		return
	}

	records, count := fakerData(req.Columns, req.Number)

	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"status": "ok",
			"code":   "0"},
		"data": gin.H{
			"count":   count,
			"records": records},
	})
}

func GetFaker(c *gin.Context) {
	// todo: 需要对Query参数进行bind，先粗暴的判断下长度
	columns := c.Query("columns")
	number := c.Query("number")
	//	fmt.Println("columns==>", columns, len(columns))
	//	fmt.Println("number==>", number, len(number))
	if len(columns) == 0 || len(number) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status": "error",
				"code":   "100"},
			"data": gin.H{
				"count":   nil,
				"records": "请输入有效的参数"},
		})
	} else {
		//		startTime := time.Now()
		records, count := fakerData(columns, number)
		//		costTime := time.Now().Sub(startTime)
		//		fmt.Println("构造数据耗时", costTime)

		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status": "ok",
				"code":   "0"},
			"data": gin.H{
				"count":   count,
				"records": records},
		})
	}

}

func fakerData(columns, number string) ([]map[string]interface{}, int) {
	itemCols := strings.Split(columns, ",")
	fakerNumber, err := strconv.Atoi(number)
	if err != nil {
		panic(err)
	}
	if fakerNumber >= 10000 {
		fakerNumber = 10000
	}
	var results []map[string]interface{}
	for i := 0; i < fakerNumber; i++ {
		resultMap := make(map[string]interface{})
		for _, col := range itemCols {
			resultMap[col] = matchFaker(strings.ToLower(col), Conn)
		}
		results = append(results, resultMap)
	}
	count := len(results)
	return results, count
}

func matchFaker(col string, c *sql.DB) interface{} {
	switch col {
	case "color":
		return faker.Color("zh_CN")
	case "job":
		return faker.Job("zh_CN")
	case "name":
		return faker.Name("zh_CN")
	case "sex":
		return faker.Gender("zh_CN")
	case "address":
		return faker.Address(c)
	case "citycode": // 中国长途区号
		return faker.CityCode()
	case "idcard":
		return faker.IdCard()
	case "age":
		return faker.Age()
	case "mobilephone": // 移动电话
		return faker.MobilePhone("zh_CN")
	case "telphone": // 固定电话
		return faker.TelPhone("zh_CN")
	case "specialphone": // 特殊号码，比如95555招商银行,10086中国移动
		return faker.SpecialTellPhone()
	case "email":
		return faker.Email()
	case "imid":
		return faker.IMID()
	case "nickname":
		return faker.NickName()
	case "username":
		return faker.UserName()
	case "password":
		return faker.PassWord(true, true, true, true, true, 10)
	case "website":
		return faker.WebSite()
	case "url":
		return faker.URL()
	case "airport":
		return faker.AirPortInfo()
	case "voyage": // 航班号
		return faker.Voyage()
	case "airlineinfo": // 航空公司信息(代号+名称)
		return faker.AirlineInfo()
	case "traintrips":
		return faker.TrainTripis()
	case "trainseat":
		return faker.SeatOfTrain()
	case "flightseat":
		return faker.SeatOfFlight()
	case "ipv4":
		return faker.IPv4Address()
	case "ipv6":
		return faker.IPv6Address()
	case "mac": // 暂时随机返回各种类型的MAC地址
		return faker.RandMacAddress()
	case "imsi": // 暂时只支持国内imsi
		return faker.Imsi()
	case "imei": //
		return faker.Imei()
	case "meid": // 随机大小写
		return faker.RandMeid()
	case "deviceid": //采集设备ID、固定21位、前9位为安全厂商ID（如FIBERHOME），后12位为采集设备MAC，规则同MAC、所有字母大写
		return faker.DeviceID()
	case "date": // 数据库日期格式{YYYYMMDD,hh:mm:ss}  (当前时间)
		return faker.NowDate()
	case "capturetime": // 10位绝对秒(当前时间)
		return faker.NowTimeStamp()
	case "useragent":
		return faker.UserAgent()
	case "carbrand":
		return faker.CarBrand("zh_CN")
	case "gapassport":
		return "暂未支持"
	case "twpassport":
		return "暂未支持"

	default:
		return "暂未支持的字段"
	}
}

func main() {
	StartServer()
	defer Conn.Close()
}
