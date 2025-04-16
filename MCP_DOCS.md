## 使用方式

参数1：columns: 需要生成数据的字段内容，见下一节清单
参数2：number: 需要生成的数据条数, number参数的取值区间为[1,10000]

使用示例：
- columns=color,job,name
- number=100

## 目前已经支持的数据类型（即columns字段的可选参数）

| 序号   |      参数      | 说明                                    |
| :--- | :----------: | ------------------------------------- |
| 1    |    color     | 颜色                                    |
| 2    |     job      | 职业                                    |
| 3    |     name     | 中文名字                                  |
| 4    |     sex      | 性别                                    |
| 5    |   address    | 地址信息（地区编号、邮编、固话区号、省市信息、社区名称、社区简称、经纬度） |
| 6    |    idcard    | 大陆居民身份证号码                             |
| 7    |     age      | 年龄                                    |
| 8    | mobilephone  | 移动电话号码                                |
| 9    |    email     | 电子邮箱                                  |
| 10   |     imid     | IM类型的用户ID                             |
| 11   |   nickname   | 用户昵称                                  |
| 12   |   username   | 用户名                                   |
| 13   |   password   | 用户密码                                  |
| 14   |   website    | 网站地址                                  |
| 15   |     url      | 网址URL（随机http或https）                   |
| 16   |   airport    | 国内机场信息（IATA编码、城市名称、ICAO编码、机场名称、城市拼音）  |
| 17   |    voyage    | 国内航班号                                 |
| 18   | airlineinfo  | 国内航空公司信息（代号、中文名称）                     |
| 19   |  traintrips  | 火车班次（覆盖高铁、动车、特快、普快、城际、旅游专线）           |
| 20   |  trainseat   | 火车座号                                  |
| 22   |  flightseat  | 飞机座号                                  |
| 23   |     ipv4     | ipv4的点分型IP地址                          |
| 24   |     ipv6     | ipv6的点分型IP地址                          |
| 25   |     mac      | mac地址（随机大小写，分隔符）                      |
| 26   |  useragent   | 浏览器请求头                                |
| 27   |     imsi     | IMSI（目前只支持国内460开头的）                   |
| 28   |     imei     | IMEI（目前支持中国、英国、美国）                    |
| 29   |     meid     | MEID（随机大小写）                           |
| 30   |   deviceid   | DEVICEID（设备编号）                        |
| 31   |   telphone   | 固定电话（暂时只支持国内号码）                       |
| 32   |   citycode   | 国内长途区号                                |
| 33   | specialphone | 特殊电话号码（比如10086、110）                   |
| 34   | capturetime  | 当前时间绝对秒（10位数字）                        |
| 35   |     date     | 当前时间，数据库日期格式{YYYYMMDD,hh:mm:ss}       |
| 36   |   carbrand   | 汽车品牌（中文）       |
