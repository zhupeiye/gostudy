package main

 import (
 	"io"
 	"log"
 	"net/http"
	"fmt"
    "io/ioutil"
    "os"
    "path"
    "errors"
    "strings"
 )

 func main() {
	 //https://github.com/k8sp/tls/blob/master/openssl.md#https-server
 	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
 		io.WriteString(w, "hello, world!\n")
 	})

    http.HandleFunc("/upload/", uploadHandle) // 上传
    http.HandleFunc("/uploaded/", showPicHandle)  //显示图片
 	if e := http.ListenAndServeTLS(":9000", "server.crt", "server.key", nil); e != nil {
		fmt.Println(e)
 		log.Fatal("ListenAndServe: ", e)
 	}
 }

// 显示图片接口
func showPicHandle( w http.ResponseWriter, req *http.Request ) {
    file, err := os.Open("." + req.URL.Path)
    errorHandle(err, w);

    defer file.Close()
    buff, err := ioutil.ReadAll(file)
    errorHandle(err, w);
    w.Write(buff)
}

// 统一错误输出接口
func errorHandle(err error, w http.ResponseWriter) {
    if  err != nil {
        w.Write([]byte(err.Error()))
    }
}

// 上传图像接口
func uploadHandle (w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    req.ParseForm()
    if req.Method != "POST" {
        w.Write([]byte(html))
    } else {
        // 接收图片
        uploadFile, handle, err := req.FormFile("image")
        errorHandle(err, w)

        // 检查图片后缀
        ext := strings.ToLower(path.Ext(handle.Filename))
        if ext != ".jpg" && ext != ".png" {
            errorHandle(errors.New("只支持jpg/png图片上传"), w);
            return
            //defer os.Exit(2)
        }

        // 保存图片
        os.Mkdir("./uploaded/", 0777)
        saveFile, err := os.OpenFile("./uploaded/" + handle.Filename, os.O_WRONLY|os.O_CREATE, 0666);
        errorHandle(err, w)
        io.Copy(saveFile, uploadFile);

        defer uploadFile.Close()
        defer saveFile.Close()
        // 上传图片成功
        w.Write([]byte("查看上传图片: <a target='_blank' href='/uploaded/" + handle.Filename + "'>" + handle.Filename + "</a>"));
    }
}
const html = `<html>
    <head></head>
    <body>
        <form method="post" enctype="multipart/form-data">
            <input type="file" name="image" />
            <input type="submit" />
        </form>
    </body>
</html>`
 

	
