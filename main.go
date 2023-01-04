package main

import (
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	common "github.com/yqhcode/paas-common"
	"log"
)

func main() {
	//1.注册中心
	consulRegistry := consul.NewRegistry(func(option *registry.Options) {
		option.Addrs = []string{
			"192.168.65.145:8500",
		}
	})
	//2.配置中心，存放经常变动的变量
	consulConfig, err := common.GetConsulConfig("192.168.65.145", 8500, "/micro/config")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(consulConfig)
	//3.使用配置中心连接 mysql
	mysqlConfig := common.GetMysqlConfig(consulConfig, "mysql")
	//初始化数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", mysqlConfig.User, mysqlConfig.Pwd, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接mysql 成功")
	defer db.Close()
	//禁止复表
	db.SingularTable(true)

	// 创建服务
	service := micro.NewService(
		micro.Name("base"),
		micro.Version("v0.01"),
		//添加注册中心
		micro.Registry(consulRegistry),
	)
	// 初始化服务
	service.Init()
	// 启动服务
	if err := service.Run(); err != nil {
		//输出启动失败信息
		log.Fatal(err)
	}
}
