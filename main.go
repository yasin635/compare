package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"os"
	"io/ioutil"
	"crypto/md5"
	"encoding/json"
	"strings"
)
//存储文件以及对应md5码
var md5files = make(map[string]string)
//校验文件信息
var md5fileInfo = make(map[string]string)
//临时路径
var fpath string
func main()  {
	//fmt.Println(os.Args)
	//需要对比的目录
	dir := flag.String("from", "./", "对比文件夹目录,默认当前目录")
	//目标文件
	md5file := flag.String("md5file", "", "md5目标文件")
	//输出md5文件目录
	tofile := flag.String("todfile", "./md5file.json", "md5文件目录输出目录")

	flag.Parse()

	if  strings.Compare(*md5file,"") < 1 {
		//读取目录 from
		readDir(*dir)
		//map转换转换字符串
		json,err := json.Marshal(md5files)
		checkError(err)
		jsonstr := fmt.Sprintf("%s", json)
		//写入文件
		//err = ioutil.WriteFile("./md5file.txt", jsonstr)
		data := []byte(jsonstr)
		writeFile(data, *tofile)
		fmt.Println("目标文件为空,将输出目录下md5 JSON文件",*tofile)
		fmt.Println(jsonstr)
	} else {
		fmt.Println("目标文件为：", *md5file)
		//判断文件是否存在
		isfile := pathExist(*md5file)
		if isfile == false {
			fmt.Println("校验失败：")
			fmt.Println("目标文件：", *md5file, " 不存在")

		}

		//一文件为准，对比目录下所有文件
		readFile(*md5file)
		//遍历文件
		for f,v :=range md5fileInfo {
			//目录文件夹
			fpath := *dir
			//文件名称
			fpath += f //组合文件路径
			//读取文件并获取文件md5值
			tmpmd5, _ := md5SumFile(fpath)
			if strings.Compare(tmpmd5,v) < 1{
				fmt.Println("文件不一致:", fpath)
			}
			//fmt.Println(f,v,fpath,tmpmd5)
		}
	}


	fmt.Println("===========================================")
	fmt.Println("from:", *dir)
	fmt.Println("md5file:", *md5file)
	fmt.Println("todfile:", *tofile)
}

//读取json文件
func readFile(file string) (map[string]string, error)  {
	bytes,err := ioutil.ReadFile(file)
	checkError(err)
	if err := json.Unmarshal(bytes, &md5fileInfo);err != nil {
		checkError(err)
	}

	return md5fileInfo,err

}
//判断文件是否存在
func pathExist(path string) bool {
	_,err := os.Stat(path)
	if err != nil && os.IsNotExist(err){
		return false
	}

	return true
}

//写入文件
func writeFile(data []byte, filename string){
	err := ioutil.WriteFile(filename, data,0777)
	checkError(err)
}

//读取文件夹目录
func readDir(root string){
	fmt.Println(root)
	filepath.Walk(root,
		//循环读取文件
	func (path string, f os.FileInfo, err error) error{
		if f == nil {
			fmt.Println("err")
			return err
		}
		if f.IsDir() {
			//fmt.Println("dir:" + path)
			return nil
		}
		//md5值计算
		//file, _ := ioutil.ReadFile(path)
		//fmt.Printf("%x", md5h.Sum([]byte(""))) //md5
		//md5Value := md5.Sum(file)
		md5Value, _ := md5SumFile(path)
		//md5Map := make(map[string]string)
		md5files[path] = md5Value
		//fmt.Println("file:", path)
		//fmt.Println("MD5:", md5Value)
		return nil
	})
}

//获取文件并获取md5值
func md5SumFile(file string) (filemd5 string, err error){
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	//获取md5
	value := md5.Sum(data)
	//格式化md5
	filemd5 = fmt.Sprintf("%x",value)

	return filemd5, nil
}

//错误处理
func checkError(err error)  {
	if err != nil{
		panic(err)
	}
}