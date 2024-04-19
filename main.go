package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

func main() {
    // 设置静态文件服务的目录
    fs := http.FileServer(http.Dir("./public"))
    http.Handle("/", fs)

    // 设置上传文件的路由
    http.HandleFunc("/upload", uploadFileHandler)

    // 监听端口
    log.Println("Listening on :11801...")
    err := http.ListenAndServe(":11801", nil)
    if err != nil {
        log.Fatal(err)
    }
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        // 解析上传的文件
        r.ParseMultipartForm(10 << 20) // 限制上传的文件大小为10MB

        // 从表单中获取文件
        file, handler, err := r.FormFile("myFile")
        if err != nil {
            fmt.Println("Error Retrieving the File")
            fmt.Println(err)
            return
        }
        defer file.Close()
        fmt.Printf("Uploaded File: %+v\n", handler.Filename)
        fmt.Printf("File Size: %+v\n", handler.Size)
        fmt.Printf("MIME Header: %+v\n", handler.Header)

        // 创建文件
        tempFile, err := ioutil.TempFile("public", handler.Filename)
        if err != nil {
            fmt.Println(err)
        }
        defer tempFile.Close()

        // 读取文件内容
        fileBytes, err := ioutil.ReadAll(file)
        if err != nil {
            fmt.Println(err)
        }
        // 写入文件
        tempFile.Write(fileBytes)

        // 返回上传成功的消息
        fmt.Fprintf(w, "Successfully Uploaded File\n")
    } else {
        // 如果不是POST请求，返回上传表单
        w.Write([]byte(`
            <html>
            <head>
            <title>Upload file</title>
            </head>
            <body>
            <form enctype="multipart/form-data" action="/upload" method="post">
            <input type="file" name="myFile" />
            <input type="submit" value="Upload" />
            </form>
            </body>
            </html>
        `))
    }
}

