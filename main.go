package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
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

        // Create a new file in the public directory
        destFile, err := os.Create(filepath.Join("./public", filepath.Base(handler.Filename)))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer destFile.Close()

        // Copy the file content to the new file
        _, err = io.Copy(destFile, file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // 返回上传成功的消息
        fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>Upload Successful</title></head><body><p>Successfully Uploaded File</p><button id=\"goToNewPage\">Back to homepage</button><script>const goToNewPageButton=document.getElementById('goToNewPage');goToNewPageButton.addEventListener('click',function(){window.location.href='/';});</script></body></html>")
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
