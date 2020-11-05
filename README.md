# lorm
<h1>我不会排版，没有那个审美
<h3>git 将master 分支改名为main 分支了！！！！！！！！！！！！
<h3>go 语言 单元测试 写法。 比如文件名a,里边方法b，那么test 文件 就叫a_test,直接敲Test,会提示Testb方法.针对每个具体的方法都会生成一个测试方法
<h3>可以相当于mian方法直接运行,go语言会自动识别你的文件是xxx_test.试试就知道了。 
<h3>  main文件要单独放，不要和其他go文件放到一个包下，不然会报错：go报错# command-line-arguments undefined: * 。相当于是springboot的application方法
<h3>要引用 并调用这个包下的init方法_ "github.com/mattn/go-sqlite3"，不然会报错。
<h3>传进去一个切片指针，将查询后的数据保存到切片，不用返回这个切片，此处为啥不传切片，而是 要传切片的 指针。因为当吧切片传过去之后，如果发生了扩容，是不会影响外部的切片的。也就是外部打印还是发生扩容之前的切片。