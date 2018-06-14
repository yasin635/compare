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
	"regexp"
)
//存储文件以及对应md5码
var md5files = make(map[string]string)
//校验文件信息
var md5fileInfo = make(map[string]string)
//临时路径
var fpath string
//总文件数
var total int32
func main()  {
	//fmt.Println(os.Args)
	//需要对比的目录
	dir := flag.String("from", "./", "对比文件夹目录,默认当前目录")
	//目标文件
	md5file := flag.String("md5file", "", "md5目标文件")
	//输出md5文件目录
	tofile := flag.String("tofile", "./md5file.json", "md5文件json目录输出目录")

	flag.Parse()

	//判断是否是json文件格式
	file := *tofile
	isjsonfile,err := regexp.MatchString(".json$", file)
	checkError(err)
	if isjsonfile == false {
		fmt.Println("tofile is not a json file")
		return
	}

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
		//总文件数
		total = int32(len(md5files))
		fmt.Println("目标文件为空,将输出目录下md5 JSON文件",*tofile)
		//fmt.Println(jsonstr)
	} else {
		//判断目录是否有添加/,没有则补上
		fpath := *dir
		reg := regexp.MustCompile("/$")
		isdir := fmt.Sprintf("%q\n", reg.FindAllString(fpath, -1))
		fmt.Println(isdir)
		fmt.Println(len(isdir))

		if len(isdir) <= 3 {
			fpath += "/"
		}
		//判断是否
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
			//fmt.Println("dir:", *dir)
			//目录文件夹
			//fpath := ""
			//文件名称
			fpath += f //组合文件路径
			//读取文件并获取文件md5值
			tmpmd5, _ := md5SumFile(fpath)
			if strings.Compare(string(tmpmd5),string(v)) < 0 {
				msg := "文件不一致:"+ fpath+"=>"+tmpmd5+":"+v
				log := fmt.Sprintf(" %c[%d;%d;%dm%s[%s]%c[0m ",0x1B, 7, 47, 31, "",msg, 0x1B)
				fmt.Println(log)
			}
			fpath = *dir;
		}
		total = int32(len(md5fileInfo))
	}


	fmt.Println("===========================================")
	fmt.Println("from:", *dir)
	fmt.Println("md5file:", *md5file)
	fmt.Println("tofile:", *tofile)
	fmt.Println("total file:", total)
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
	//fmt.Println(root)
	filepath.Walk(root,
		//循环读取文件
	func (path string, f os.FileInfo, err error) error{
		if f == nil {
			fmt.Println("err")
			return err
		}
		//判断是否权限读取
		//file_mode := f.Mode()
		//fmt.Println("file_model:", file_mode)
		//prem := file_mode.Perm()
		//fmt.Println("permission:", uint32(prem))
		//// 73: 000 001 001 001
		//flag := prem & os.FileMode(73)
		//if uint32(flag) == uint32(73){
		//	fmt.Println("没有权限：")
		//	//return nil
		//}
		//ok := strings.HasSuffix(f.Name(), ".git")
		//fmt.Println("is git:", ok," ->file:",f.Name())
		//fmt.Println("is git2:",strings.Contains(f.Name(), ".git"))
		//if strings.Contains(".git/refs/remotes/origin", ".git") == false {
		//	fmt.Println("is git:",strings.Contains(f.Name(), ".git"))
		//}
		//if strings.Contains(path, ".git") == false {
		//	fmt.Println("is git2:",strings.Contains(f.Name(), ".git"),"->file:",f.Name())
		//}
		//return nil
		if f.IsDir() && (path != "./" && path != ".." && path != "." && path != root) && strings.Contains(path, ".git") == false{
			//fmt.Println("dir:" + path)
			readDir(path) //递归读取文件夹目录
			//return nil
		}

		if strings.Contains(path, ".git") == true || path == "./" || path == ".." && path == "." || path == root {
			//fmt.Println("is git2:",strings.Contains(f.Name(), ".git"),"->file:",f.Name())
			return nil
		}
		//md5值计算
		//file, _ := ioutil.ReadFile(path)
		//fmt.Printf("%x", md5h.Sum([]byte(""))) //md5
		//md5Value := md5.Sum(file)
		md5Value, _ := md5SumFile(path)
		//md5Map := make(map[string]string)
		realtivePath := strings.Replace(path,root,"",-1)
		md5files[realtivePath] = md5Value
		//fmt.Println("file:", path)
		//fmt.Println("MD5:", md5Value)
		fmt.Println("pathFile:",path," ",md5Value)
		return nil
	})
}

//获取文件并获取md5值
func md5SumFile(file string) (filemd5 string, err error){
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "err", nil
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