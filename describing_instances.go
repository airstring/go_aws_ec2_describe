package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//checking if your credentials have been found is fairly easy.

// sess, err := session.NewSession(&aws.Config{
//     Region:      aws.String("us-west-2"),
//     Credentials: credentials.NewSharedCredentials("", "test-account"),
// })

//csv文件写入
func WriterCSV(path string, line []string) {
	//OpenFile读取文件，不存在时则创建，使用追加模式
	File, err1 := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err1 != nil {
		log.Println("文件打开失败!")
	}
	defer File.Close()
	//创建写入接口
	wcsv := csv.NewWriter(File)
	err2 := wcsv.Write(line) //写入一条数据，传入数据为切片(追加模式)
	if err2 != nil {
		log.Println("WriterCsv写入文件失败")
	}
	wcsv.Flush() //刷新，不刷新是无法写入的
	log.Println("数据写入成功...")
}

func main() {
	// Load session from shared config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	_, err := sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println("get session err")
	}
	// Create new EC2 client
	ec2Svc := ec2.New(sess)
	// Call to get detailed information on each instance
	result, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		//fmt.Println(reflect.TypeOf(result))
		//result.Reservations是一个结构体类型的切片
		instance_list := result.Reservations

		for _, v1 := range instance_list {
			//instance是一个结构体类型的切片
			instance := v1.Instances
			for _, v2 := range instance {
				//fmt.Println(*v2.InstanceType)
				//fmt.Println(v2.Tags)
				for _, v3 := range v2.Tags {
					if *v3.Key == "Name" {
						//fmt.Println(*v3.Value)
						line := []string{*v3.Value, *v2.InstanceType}
						WriterCSV("/tmp/instance.csv", line)
					}
				}
			}
		}
	}

}
