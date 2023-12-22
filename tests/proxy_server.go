/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-22 11:46:04
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-22 14:24:57
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package main

import (
	"fmt"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ip", func(c *gin.Context) {

		// 随机生成一个数
		rand1 := rand.Intn(255)
		rand2 := rand.Intn(255)
		rand3 := rand.Intn(255)
		rand4 := rand.Intn(255)
		ip := fmt.Sprintf("%d.%d.%d.%d", rand1, rand2, rand3, rand4)
		c.String(200, ip)
	})

	r.Run(":8001") // listen and serve on
}
