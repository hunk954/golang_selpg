# 基于golang开发 Linux 命令行实用程序selpg
## 一、selpg简要概述
 selpg是从文本输入选择页范围的实用程序。该文本输入可以来自作为最后一个命令行参数指定的文件，在没有给出文件名参数时也可以来自标准输入

## 二、selpg使用及参数说明
```
USAGE: selpg -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]
  -e, --endpage int              Index of the end page (default -1)
  -d, --pageDestination string   printer destination (default "-")
  -l, --pageLength int           page length (default -1)
  -f, --pageType                 fragment  according to "f" 
  -s, --startpage int            Index of the start page (default -1)
```
- 通过命令行参数，selpg读取初始页与结束页，所以selpg包括文件名在内**至少有三个参数**。`selpg -s[startpage] -e[endpage]`。这时的功能表示从标准输入中读取初始页到结束页的页数范围内的内容。如果在参数最后加上文件名，则表示从输入文件中读取对应的内容。下面讲述其他可选参数的功能
- `lNumber`：表示文本读取对页的定义是一页Number行，程序默认为按72行为一页读取。且该参数与`-f`是互斥的
- `-f`：文本读取默认是按页中行数为72读取，加入该参数表示页的定义由分页符`'\f'`决定。加入该参数表示页定义类型覆盖原来用页行数定义的类型
- `-dDestination`：表示文本选定的页将发送到打印机，Destination”应该是 lp 命令“-d”选项可接受的打印目的地名称。该目的地应该存在 ― selpg 不检查这一点。在运行了带“-d”选项的 selpg 命令后，若要验证该选项是否已生效，请运行命令“lpstat -t”。该命令应该显示添加到“Destination”打印队列的一项打印作业。如果当前有打印机连接至该目的地并且是启用的，则打印机应打印该输出


### 三、selpg程序说明