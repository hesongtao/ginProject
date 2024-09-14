package user

import (
	"fmt"
	"ginProject/app/model/product"
	"ginProject/common/cache"
	"ginProject/common/queue"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	// 获取get方式参数
	// username := c.Param("username")
	data := make(map[string]interface{}, 0)
	err := hystrix.Do("aaa", func() error {
		username := c.Query("username")
		// 获取post方式参数
		// username := c.PostForm("username")
		// 获取header参数
		// c.Request.Header.Get("username")
		err := queue.AMQP.Publish("", "my_queue", []byte("Hi,  Mbit!"))
		if err != nil {
			fmt.Println("Failed to publish message:", err)
		}
		//err = queue.AMQP.Consume("my_queue", func(delivery amqp.Delivery) {
		//	fmt.Println(string(delivery.Body))
		//})
		//if err == nil {
		//	fmt.Println("Failed to consume message:", err)
		//}
		product.CreateProduct(username)
		cache.RC.Set("abc", 100, 0)

		//_ = models.GetProductsById(1)
		myname := product.GetProductsById(1)
		data["username"] = username
		data["myname"] = myname
		data["abc"], _ = cache.RC.Get("abc")
		//time.Sleep(2 * time.Second)
		fmt.Println("do services")
		return nil
	}, func(err error) error {
		fmt.Println("短暂出现了错误，请让服务器休息一下")
		return err
	})

	//time.Sleep(10 * time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"status":  http.StatusOK,
		"data":    data,
	})
	c.Abort()
	return
}

func Info(c *gin.Context) {
	limit := ratelimit.NewBucket(3*time.Second, 5)
	//print("--------limit token left:")
	//print(limit.TakeAvailable(1))
	//print(limit.TakeAvailable(1))
	for i := 0; i < 20; i++ {
		// 获取当前时间
		startTime := time.Now()
		limit.Wait(1)
		//if limit.TakeAvailable(1) < 1 {
		//	print("--------limit token less----------")
		//	//c.JSON(http.StatusOK, gin.H{
		//	//	"message": "rate limit...",
		//	//	"status":  http.StatusOK,
		//	//})
		//	//c.String(http.StatusOK, "rate limit...")
		//	//c.Abort()
		//	//return
		//}

		// 计算处理时间
		elapsedTime := time.Since(startTime)
		fmt.Printf("Processing time: %dms", elapsedTime.Milliseconds())
	}

	// 获取get方式参数
	// username := c.Param("username")
	username := c.Query("username")
	// 获取post方式参数
	// username := c.PostForm("username")
	// 获取header参数
	// c.Request.Header.Get("username")

	data := make(map[string]interface{}, 0)

	data["username"] = username

	c.JSON(http.StatusOK, gin.H{
		"message": "success_ok",
		"status":  http.StatusOK,
		"data":    data,
	})
	c.Abort()
	return
}
