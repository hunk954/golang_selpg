# 基于golang开发 Linux 命令行实用程序selpg
## 一、selpg简要概述
 selpg是从文本输入选择页范围的实用程序。该文本输入可以来自作为最后一个命令行参数指定的文件，在没有给出文件名参数时也可以来自标准输入

## 二、selpg使用及参数说明
```
USAGE: selpg -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]
  -s, --startpage int            Index of the start page (default -1)
  -e, --endpage int              Index of the end page (default -1)
  -d, --pageDestination string   printer destination (default "-")
  -l, --pageLength int           page length (default -1)
  -f, --pageType                 fragment  according to "f" 
```
- 通过命令行参数，selpg读取初始页与结束页，所以selpg包括文件名在内**至少有三个参数**。  
    ```selpg -s[startpage] -e[endpage]```  
    这时的功能表示从**标准输入**中读取初始页到结束页的页数范围内的内容。  
    如果在参数最后加上文件名，则表示从**输入文件**中读取对应的内容。
- 下面讲述其他可选参数的功能
    - `lNumber`：表示文本读取对页的定义是一页Number行，程序默认为按72行为一页读取。且该参数与`-f`是互斥的
    - `-f`：文本读取默认是按页中行数为72读取，加入该参数表示页的定义由分页符`'\f'`决定。加入该参数表示页定义类型覆盖原来用页行数定义的类型
    - `-dDestination`：表示文本选定的页将发送到打印机，Destination”应该是 lp 命令“-d”选项可接受的打印目的地名称。该目的地应该存在 ― selpg 不检查这一点。在运行了带“-d”选项的 selpg 命令后，若要验证该选项是否已生效，请运行命令“lpstat -t”。该命令应该显示添加到“Destination”打印队列的一项打印作业。如果当前有打印机连接至该目的地并且是启用的，则打印机应打印该输出


## 三、selpg程序代码实现
### 1.与C语言的对比
- 参数的读取：C语言的处理方式是通过main()中的argv[] 进行读取分析。golang虽然在os包中也有类似的os.Args[]为我们实现同样的功能，但是同时也提供了一种更好的分析参数的方法：flag。  
```golang
flag.IntVar(&name,"flag",defaultValue,"info")
```  
该代码可以让程序自动读取`'-flag'`的参数值，其值会存入变量`name`中，如果我们并没有输入该参数，那么变量的值缺省为`defaultValue`，在Usage中会提供提示信息`“info”`提示我们该参数的用处。对于string、bool的设置只要将方法名中的Int改成对应的类型名即可（注：导出的函数首字母大写）。


### 2.基本框架
```golang
//用来存储读取页面的参数结构体
type selpgArgs struct{
	startPage int
	endPage int
	inFilename string
	pageLen int
	pageType int
	printDest string
}
type spArgs selpgArgs //简写，方便代码书写

func usage()//程序出错时打印使用说明
func processArgs(psa *spArgs)//处理参数，参数传指针储存参数
func processInput(sa spArgs)//根据参数要求执行文本输入
func init()//golang特性与main一样会被默认执行，这里放设置pflag参数的代码
func main()
```
## 四、selpg使用测试