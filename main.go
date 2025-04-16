package main

import (
	"context"
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
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/toddlerya/fakerfactory/faker"
)

var mode string
var port string
var dbPath string
var Conn *sql.DB

func init() {
	if len(os.Args) > 1 {
		mode = os.Args[1]
	} else {
		mode = "server"
	}
	if len(os.Args) > 2 {
		port = os.Args[2]
	} else {
		port = "8001"
	}
	if len(os.Args) > 3 {
		dbPath = os.Args[3]
	} else {
		dbPath = `./data/data.db`
	}
	fmt.Println("args==>", os.Args)
	fmt.Println("args[1]==>", os.Args[1])
	fmt.Println("args[2]==>", os.Args[2])
	fmt.Println("args[3]==>", os.Args[3])
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
	}
	err := router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("在%s端口启动服务失败！", port)
	}
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
	case "deviceid": //采集设备ID、固定21位、前9位为安全厂商ID(如FIBERHOME)，后12位为采集设备MAC，规则同MAC、所有字母大写
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

func StartMCPServer() {
	// TODO MCP服务
	fmt.Println("MCP服务启动...")
	// Create a new MCP server
	s := server.NewMCPServer(
		"fakerfactory",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Add a faker data tool
	fakerTool := mcp.NewTool("fakerfactory",
		mcp.WithDescription("生成仿真测试数据"),
		mcp.WithString("columns",
			mcp.Required(),
			mcp.Description(`需要生成的仿真数据字段参数, 用英文逗号分隔
目前已经支持的数据类型(即columns字段的可选参数)
| 序号   |      参数      | 说明                                    |
| :--- | :----------: | ------------------------------------- |
| 1    |    color     | 颜色                                    |
| 2    |     job      | 职业                                    |
| 3    |     name     | 中文名字                                  |
| 4    |     sex      | 性别                                    |
| 5    |   address    | 地址信息(地区编号、邮编、固话区号、省市信息、社区名称、社区简称、经纬度) |
| 6    |    idcard    | 大陆居民身份证号码                             |
| 7    |     age      | 年龄                                    |
| 8    | mobilephone  | 移动电话号码                                |
| 9    |    email     | 电子邮箱                                  |
| 10   |     imid     | IM类型的用户ID                             |
| 11   |   nickname   | 用户昵称                                  |
| 12   |   username   | 用户名                                   |
| 13   |   password   | 用户密码                                  |
| 14   |   website    | 网站地址                                  |
| 15   |     url      | 网址URL(随机http或https)                   |
| 16   |   airport    | 国内机场信息(IATA编码、城市名称、ICAO编码、机场名称、城市拼音)  |
| 17   |    voyage    | 国内航班号                                 |
| 18   | airlineinfo  | 国内航空公司信息(代号、中文名称)                     |
| 19   |  traintrips  | 火车班次(覆盖高铁、动车、特快、普快、城际、旅游专线)           |
| 20   |  trainseat   | 火车座号                                  |
| 22   |  flightseat  | 飞机座号                                  |
| 23   |     ipv4     | ipv4的点分型IP地址                          |
| 24   |     ipv6     | ipv6的点分型IP地址                          |
| 25   |     mac      | mac地址(随机大小写，分隔符)                      |
| 26   |  useragent   | 浏览器请求头                                |
| 27   |     imsi     | IMSI(目前只支持国内460开头的)                   |
| 28   |     imei     | IMEI(目前支持中国、英国、美国)                    |
| 29   |     meid     | MEID(随机大小写)                           |
| 30   |   deviceid   | DEVICEID(设备编号)                        |
| 31   |   telphone   | 固定电话(暂时只支持国内号码)                       |
| 32   |   citycode   | 国内长途区号                                |
| 33   | specialphone | 特殊电话号码(比如10086、110)                   |
| 34   | capturetime  | 当前时间绝对秒(10位数字)                        |
| 35   |     date     | 当前时间，数据库日期格式{YYYYMMDD,hh:mm:ss}       |
| 36   |   carbrand   | 汽车品牌(中文)       |`),
			mcp.Enum("color", "job", "name", "sex", "address", "idcard", "age", "mobilephone", "email", "imid", "nickname", "username", "password", "website", "url", "airport", "voyage", "airlineinfo", "traintrips", "trainseat", "flightseat", "ipv4", "ipv6", "useragent", "mac", "imsi", "imei", "meid", "deviceid", "telphone", "citycode", "specialphone", "capturetime", "date", "carbrand"),
		),
		mcp.WithNumber("number",
			mcp.Required(),
			mcp.Description("需要生成的数据条数"),
		),
	)

	// Add the fakerfactory handler
	s.AddTool(fakerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		columns := request.Params.Arguments["columns"].(string)
		number := request.Params.Arguments["number"].(string)

		result, count := fakerData(columns, number)

		return mcp.NewToolResultText(fmt.Sprintf("生成数据条数: %d\n生成数据内容: %v", count, result)), nil
	})

	// Static resource example - exposing a README file
	resource := mcp.NewResource(
		"docs://readme",
		"Project README",
		mcp.WithResourceDescription("The project's README file"),
		mcp.WithMIMEType("text/markdown"),
	)

	// Add resource with its handler
	s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content, err := os.ReadFile("MCP_DOCS.md")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",
				MIMEType: "text/markdown",
				Text:     string(content),
			},
		}, nil
	})

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func main() {
	if mode == "server" {
		fmt.Println("API服务模式")
		StartServer()
	} else if mode == "mcp" {
		fmt.Println("MCP服务模式")
		StartMCPServer()
	}
	defer Conn.Close()
}
