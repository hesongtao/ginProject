package main

import (
	"ginProject/common/cache"
	"ginProject/common/dao"
	"ginProject/common/queue"
	"ginProject/router"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"time"
)

var CircuitBreakerName = "api_circuit_breaker"

func main() {

	router := router.InitRouter()
	dao.Init()
	cache.Init()
	queue.Init()
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
	hystrix.ConfigureCommand("api_circuit_breaker", hystrix.CommandConfig{
		Timeout:                int(1 * time.Second), // 执行command的超时时间为3s
		MaxConcurrentRequests:  5,                    // command的最大并发量
		RequestVolumeThreshold: 5,                    // 统计窗口10s内的请求数量，达到这个请求数量后才去判断是否要开启熔断
		SleepWindow:            int(2 * time.Second), // 当熔断器被打开后，SleepWindow的时间就是控制过多久后去尝试服务是否可用了
		ErrorPercentThreshold:  20,                   // 错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
	})
	log := logrus.New()
	// 设置日志格式为JSON，便于Logstash解析
	log.Formatter = &logrus.JSONFormatter{}

	file, err := os.OpenFile("log/logfile.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	log.Info("这是一条日志信息")

	if httperr := router.Run(":8081"); httperr != nil {
		panic(httperr)
	}

	//r := gin.Default()
	//r.GET("/api/ping/baidu", func(c *gin.Context) {
	//	//_, err := http.Get("https://www.baidu.com")
	//	err := hystrix.Do("aaa", func() error {
	//		//test case 1 并发测试
	//		//if i == 0 {
	//		//	return errors.New("service error")
	//		//}
	//		//_, httperr := http.Get("https://www.baidu.com")
	//
	//		//time.Sleep(1 * time.Second)
	//		//if httperr != nil {
	//		//	fmt.Println("hystrix err:" + httperr.Error())
	//		//	time.Sleep(1 * time.Second)
	//		//	fmt.Println("sleep 1 second")
	//		//}
	//		//test case 2 超时测试
	//		time.Sleep(2 * time.Second)
	//		fmt.Println("do services")
	//		return nil
	//	}, func(err error) error {
	//		fmt.Println("短暂出现了错误，请让服务器休息一下")
	//		return err
	//	})
	//
	//	//time.Sleep(10 * time.Second)
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
	//		return
	//	}
	//	c.JSON(http.StatusOK, gin.H{"msg": "success"})
	//})
	//r.Run()
}

//func main() {
//	router := router.InitRouter()
//	dao.Init()
//	cache.Init()
//	queue.Init()
//
//	if err := router.Run(":8081"); err != nil {
//		panic(err)
//	}
//	//consulReg := consul.NewRegistry(func(options *registry.Options) {
//	//	options.Addrs = []string{
//	//		"localhost:8500", // Consul服务地址
//	//	}
//	//})
//	//
//	//// 创建新的服务，使用consul注册中心
//	//service := micro.NewService(
//	//	micro.Name("my.service"),
//	//	micro.Registry(consulReg),
//	//)
//	//
//	//// 初始化服务，并启动
//	//service.Init()
//	//service.Run()
//}
