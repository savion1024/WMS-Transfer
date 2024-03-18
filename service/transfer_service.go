package service

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	// 从r.header里面拿参数
	url := r.Header.Get("url")
	// 创建要转发的请求
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, "无法创建请求", http.StatusInternalServerError)
		return
	}

	// 复制原始请求的请求头到转发请求
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 创建 HTTP 客户端并发送转发请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "请求转发失败", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 读取转发响应的内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "读取响应失败", http.StatusInternalServerError)
		return
	}
	// 打印耗时时间  格式为毫秒
	costTime := time.Since(startTime)
	log.Printf("耗时: %d ms, url: %s", costTime.Milliseconds(), url)

	// 将转发响应的内容写入原始请求的响应
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

}
