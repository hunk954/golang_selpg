# 基于golang开发 Linux 命令行实用程序selpg
## 一、selpg简要概述
 selpg是从文本输入选择页范围的实用程序。该文本输入可以来自作为最后一个命令行参数指定的文件，在没有给出文件名参数时也可以来自标准输入

## 二、selpg使用及参数说明
- clone到本地后执行`go install github.com/user/selpg`之后可以使用
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
### 1、与C语言的对比
- 参数的读取：C语言的处理方式是通过main()中的argv[] 进行读取分析。golang虽然在os包中也有类似的os.Args[]为我们实现同样的功能，但是同时也提供了一种更好的分析参数的方法：flag。而pflag与flag大同小异，所以这里主要讲flag，实际代码用的是pflag。 [参考：golang中使用flag与pflag](https://o-my-chenjian.com/2017/09/20/Using-Flag-And-Pflag-With-Golang/)  
  - `flag.IntVar()`可以让程序自动读取`'-flag'`的参数值，其值会存入变量`name`中，如果我们并没有输入该参数，那么变量的值缺省为`defaultValue`，在Usage中会提供提示信息`“info”`提示我们该参数的用处。对于string、bool的设置只要将方法名中的Int改成对应的类型名即可（注：导出的函数首字母大写）。    
  - `flag.Usage`主要是设置我们的使用手册调用的函数   
  - `flag.Parse()`用于析取参数，运行该命令之后，我们参数值对应的变量就会被刷新为我们输入的参数值
    ```golang
    flag.IntVar(&name,"flag",defaultValue,"info")
    flag.Usage = usage
    flag.Parse()
    ```  
- fprintf()与exit()方法：主要是调用`fmt.Fprintf()`和`os.Exit()`方法，效果和C语言大同小异
- 文件指针：C语言中的`FILE*`对应`*os.FILE`,`stderr/stdin/stdout`对应`os.Stderr/os.Stdin/os.Stdout`
- 文件读取：过程与C语言相仿，首先获取指针后声明缓存空间，读取文本，输出文本，使用的函数接口如下
    ```golang
    //----1.获得输入文件指针----
    inputFile, err := os.Open(filename)
    // inputfile := os.Stdin 
    if err != nil{ 
      //error process
    }
    //----2.声明缓存器-----
    inputReader := bufio.newReader(inputFile)

    //----3.读取文本-----
    inputString, err := inputReader.ReadString('\n')  //表示以'\n'为分界

    //----4.获取输出文件指针并输出----
    stdout := os.Stdout
    stdout.WriteString(inputString)
    ```


### 2、基本框架
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

### 3、主要函数分析
#### processArgs(psa *spArgs)
- 由于我们已经使用了`pflag`,所以参数的读取非常的简单。只要一次检查每个参数变量的值即可。首先使用`len(os.Args)`判断参数是否达到至少三个的要求，如果不符合则打印错误信息并退出程序。
- 检查`startPage(-s)`和`endPage(-e)`参数，我们将缺省值设置为-1.若没有输入（值仍为缺省）或者输入不符合需求（1<page<&infin;）则打印错误信息并退出。其中无穷的判断使用`math.IsInf(float64(num),0)`
- 检查`pagelength(-l)`参数，缺省值为-1。此处的-1只是标志，若-1表明未使用-l参数，若使用了则需要判断参数变量是否符合需求。由于缺省的文本读取类型就是'l'，按行读取，所以该参数检查没有什么特别的注意事项
- 检查`-f`参数，是bool类型的变量，缺省值为false。如果用户使用了该参数，则值为True，将文本读取类型改为'f'即可。由于-f和-l是互斥的，所以需要判断`f && length != -1`，如果满足该条件则打印错误信息，因为-f和-l参数被同时输入了。
- 检查`-d`参数，该参数我们缺省为'-'，因为如果缺省为空，我们无法判断用户是否输入该参数（可能用户只输入了-d），所以这样设置。如果值为空，证明用户只输入了-d，打印错误信息。如果为'-'，那么用户未输入该参数，将psa.printDest设置为空，否则设置为对应参数即可。
- 此时还需要声明已记录的参数数量argno，主要是用于读取上述参数后，判断argno与len(os.Args)的关系，如果小于，则证明还需要读取可能存在的、非特定参数的（即没有前缀‘-’）的输入文本的文件名称。
    ```golang
    //伪代码
    if 参数 < 1 || math.IsInf(float64(参数),0) || 不符合需求{
      fmt.Fprintf(os.Stderr,"错误信息")
    }else{
      psa.参数名 = 读入参数
    }
    ```
#### processInput(sa spArgs)
主要是文件读取代码的书写，过程同*与C语言的对比-文件读取*中的代码。

## 四、selpg使用测试
