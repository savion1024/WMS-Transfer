package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

// JsonResult 返回结构
type JsonResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// IndexHandler 计数器接口
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	data, err := Login()
	if err != nil {
		return
	}
	userTime := time.Since(startTime)
	ret := map[string]interface{}{
		"data": data,
		"cost": userTime,
	}
	res := &JsonResult{
		Data:     ret,
		Code:     0,
		ErrorMsg: "",
	}
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Write(msg)

}

// CounterHandler 计数器接口
func CounterHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}

	if r.Method == http.MethodGet {
		counter, err := getCurrentCounter()
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = counter.Count
		}
	} else if r.Method == http.MethodPost {
		count, err := modifyCounter(r)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = count
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// modifyCounter 更新计数，自增或者清零
func modifyCounter(r *http.Request) (int32, error) {
	action, err := getAction(r)
	if err != nil {
		return 0, err
	}

	var count int32
	if action == "inc" {
		count, err = upsertCounter(r)
		if err != nil {
			return 0, err
		}
	} else if action == "clear" {
		err = clearCounter()
		if err != nil {
			return 0, err
		}
		count = 0
	} else {
		err = fmt.Errorf("参数 action : %s 错误", action)
	}

	return count, err
}

// upsertCounter 更新或修改计数器
func upsertCounter(r *http.Request) (int32, error) {
	currentCounter, err := getCurrentCounter()
	var count int32
	createdAt := time.Now()
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		count = 1
		createdAt = time.Now()
	} else {
		count = currentCounter.Count + 1
		createdAt = currentCounter.CreatedAt
	}

	counter := &model.CounterModel{
		Id:        1,
		Count:     count,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}
	err = dao.Imp.UpsertCounter(counter)
	if err != nil {
		return 0, err
	}
	return counter.Count, nil
}

func clearCounter() error {
	return dao.Imp.ClearCounter(1)
}

// getCurrentCounter 查询当前计数器
func getCurrentCounter() (*model.CounterModel, error) {
	counter, err := dao.Imp.GetCounter(1)
	if err != nil {
		return nil, err
	}

	return counter, nil
}

// getAction 获取action
func getAction(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return "", err
	}
	defer r.Body.Close()

	action, ok := body["action"]
	if !ok {
		return "", fmt.Errorf("缺少 action 参数")
	}

	return action.(string), nil
}

// getIndex 获取主页
func getIndex() (string, error) {
	b, err := ioutil.ReadFile("./index.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Login() (string, error) {
	url := "https://wl.yuanxing-wms.com/wms/auth/login" // 替换为实际的第三方接口地址

	// 准备请求体数据
	payload := map[string]interface{}{
		"email":    "wl@gmail.com",
		"password": "password123",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "请求体序列化失败", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "创建请求失败", err
	}

	req.Header.Set("Content-Type", "application/json") // 设置请求头为 JSON 类型

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "发送请求失败", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "读取响应失败", err
	}
	return string(body), nil

}
