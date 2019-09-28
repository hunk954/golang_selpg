package main

import(
	"fmt"
	"os"
	"math"
	"bufio"
	"os/exec"
	"github.com/spf13/pflag"
)
var (
	spg int
	epg int
	length int
	f bool //判断是否分页
	d string //printer destination
)

type selpgArgs struct{
	startPage int
	endPage int
	inFilename string
	pageLen int
	pageType int
	printDest string
}

type spArgs selpgArgs

//INBUFSIZE 表示缓冲区大小
const INBUFSIZE int = 16*1024

var progname string

func usage(){
	fmt.Fprintf(os.Stderr, "\nUSAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n", progname)
	pflag.PrintDefaults()
}

func processArgs(psa *spArgs){
	// fmt.Printf("len of Args: %d\n", len(os.Args))
	var argno int // 记录被process的参数个数
	if len(os.Args) < 3{
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
		pflag.Usage()
		os.Exit(1)
	}
	argno = 1
	if spg == -1 {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstart_page\n", progname)
		pflag.Usage()
		os.Exit(2)
	}
	if spg < 1 || math.IsInf(float64(spg),0){
		fmt.Fprintf(os.Stderr, "%s: invalid start page %d\n", progname, spg);
		pflag.Usage()
		os.Exit(3)
	}
	psa.startPage = spg
	argno++

	if epg == -1{
		fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -eend_page\n", progname);
		pflag.Usage()
		os.Exit(4);
	}
	if epg < 1 || math.IsInf(float64(epg),0) || epg < psa.startPage{
		fmt.Fprintf(os.Stderr, "%s: invalid end page %d\n", progname, epg);
		pflag.Usage()
		os.Exit(5)
	}
	psa.endPage = epg
	argno++

	if f && length != -1{
		fmt.Fprintf(os.Stderr, "%s: Fix length and fragment according to \"f\" are mutex",progname)
		pflag.Usage()
		os.Exit(6)
	}
	if f{
		psa.pageType = 'f'
		argno++
	}
	if length != -1{
		if length < 1 || math.IsInf(float64(length), 0){
			fmt.Fprintf(os.Stderr, "%s: invalid page length %d\n", progname, length)
			pflag.Usage()
			os.Exit(7)
		}
		argno++
		psa.pageLen = length
	}
	if d ==  ""{
		fmt.Fprintf(os.Stderr,"%s: -d option requires a printer destination\n", progname)
		pflag.Usage()
		os.Exit(8)
	}else{
		if d == "-"{
			psa.printDest = ""
		}else{
			psa.printDest = d
			argno++
		}
	}
	if argno <= len(os.Args)-1{
		psa.inFilename = os.Args[argno]
		_, err := os.Stat(psa.inFilename)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: input file \"%s\" does not exist\n",progname,psa.inFilename)
			os.Exit(10)
		}
		//此处不知道怎么判断是否可读
	}
	fmt.Fprintf(os.Stderr,"DEBUG: psa.start_page = %d\n", psa.startPage);
	fmt.Fprintf(os.Stderr,"DEBUG: psa.end_page = %d\n", psa.endPage);
	fmt.Fprintf(os.Stderr, "DEBUG: psa.page_len = %d\n", psa.pageLen);
	fmt.Fprintf(os.Stderr, "DEBUG: psa.page_type = %c\n", psa.pageType);
	fmt.Fprintf(os.Stderr, "DEBUG: psa.print_dest = %s\n", psa.printDest);
	fmt.Fprintf(os.Stderr, "DEBUG: psa.in_filename = %s\n", psa.inFilename);
}

func processInput(sa spArgs){
	var stdout *os.File
	var inputFile *os.File
	var pageCtr int
	//如果sa.inFilename为空，那么就为默认的输入流
	if sa.inFilename != ""{
		inputFile, _ = os.Open(sa.inFilename)
	}else{
		inputFile = os.Stdin
	}
	defer inputFile.Close()  
	inputReader := bufio.NewReader(inputFile) //读取器

	//如果sa.printDest为空，那么就为默认的输出流
	if sa.printDest != ""{
		s1 := "lp" 
		s2 := "-d"+ sa.printDest
		cmd := exec.Command(s1,s2)
		stdout, err := cmd.StdoutPipe()
		if err != nil {     //获取输出对象，可以从该对象中读取输出结果
			fmt.Fprintf(os.Stderr, "%s: could not open pipe to \"%s\"\n", progname, s1)
			os.Exit(13)
		}
		cmd.Start()
		defer stdout.Close()   	
	}else{
		stdout = os.Stdout
	}
	if sa.pageType == 'l'{
		lineCtr := 1 //line counter
		pageCtr = 1 // page counter
		for {
			inputString, readerError := inputReader.ReadString('\n')
			if readerError != nil{
				break
			}
			if pageCtr >= sa.startPage && pageCtr <= sa.endPage{
					stdout.WriteString(inputString)
			}
			lineCtr++
			if lineCtr > sa.pageLen{
				pageCtr++
				lineCtr = 1
			}
		}
	}else{
		pageCtr = 1
		for{
			inputString, readerError := inputReader.ReadString('\f')
			if readerError != nil{
				break
			}
			if pageCtr >= sa.startPage && pageCtr <= sa.endPage{
				stdout.WriteString(inputString)
			}
			pageCtr++
		}
	}
	if pageCtr < sa.startPage{
		fmt.Fprintf(os.Stderr,"%s: start_page (%d) greater than total pages (%d), no output written\n", progname, sa.startPage, pageCtr)
	}else if pageCtr < sa.endPage{
		fmt.Fprintf(os.Stderr,"%s: end_page (%d) greater than total pages (%d), less output than expected\n", progname, sa.endPage, pageCtr)
	}

}

func init(){
	pflag.IntVarP(&spg, "startpage","s", -1, "Index of the start page" )
	pflag.IntVarP(&epg, "endpage","e", -1, "Index of the end page")
	pflag.IntVarP(&length, "pageLength", "l", -1, "page length")
	pflag.BoolVarP(&f, "pageType", "f", false, "fragment  according to \"f\" ")
	pflag.StringVarP(&d, "pageDestination", "d", "-", "printer destination")
	pflag.Usage = usage
}


func main(){
	var sa spArgs
	progname = os.Args[0]
	sa.startPage = -1
	sa.endPage = -1
	sa.inFilename = ""
	sa.pageLen = 72
	sa.pageType = 'l'
	sa.printDest = ""
	pflag.Parse()
	processArgs(&sa)
	processInput(sa)
}

