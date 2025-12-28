package monitor

import (
	"pocketbase-ruoyi/tools"

	"github.com/pocketbase/pocketbase/core"
)

// OperLogInput 操作日志入库参数
type OperLogInput struct {
	TenantID      string `json:"tenant_id" form:"tenant_id"`
	Title         string `json:"title" form:"title"`
	BusinessType  string `json:"business_type" form:"business_type"` // 0/1/2/3
	OperatorType  string `json:"operator_type" form:"operator_type"` // 0/1/2
	Status        string `json:"status" form:"status"`               // 0/1
	Method        string `json:"method" form:"method"`
	RequestMethod string `json:"request_method" form:"request_method"`
	OperName      string `json:"oper_name" form:"oper_name"`
	DeptName      string `json:"dept_name" form:"dept_name"`
	OperURL       string `json:"oper_url" form:"oper_url"`
	OperIP        string `json:"oper_ip" form:"oper_ip"`
	OperLocation  string `json:"oper_location" form:"oper_location"`
	OperParam     string `json:"oper_param" form:"oper_param"`
	JSONResult    string `json:"json_result" form:"json_result"`
	ErrorMsg      string `json:"error_msg" form:"error_msg"`
	CostTime      int64  `json:"cost_time" form:"cost_time"`
}

// RecordOperLog 记录一条操作日志（无路由，纯函数）
func RecordOperLog(e *core.RequestEvent, in OperLogInput) error {
	if in.TenantID == "" {
		if tid := tools.GetUserTenant(e); tid != "" {
			in.TenantID = tid
		} else {
			in.TenantID = "000000"
		}
	}
	if in.RequestMethod == "" && e.Request != nil {
		in.RequestMethod = e.Request.Method
	}
	if in.OperURL == "" && e.Request != nil {
		in.OperURL = e.Request.URL.Path
	}
	if in.OperIP == "" {
		in.OperIP = tools.GetIPAddr(e.Request)
	}
	if in.OperLocation == "" {
		in.OperLocation = tools.GetLocationByIP(in.OperIP)
	}
	if in.OperName == "" && e.Auth != nil {
		in.OperName = e.Auth.GetString("user_name")
	}
	if in.DeptName == "" && e.Auth != nil {
		in.DeptName = e.Auth.GetString("dept_name")
	}

	col, err := e.App.FindCollectionByNameOrId("oper_log")
	if err != nil {
		return err
	}

	rec := core.NewRecord(col)
	rec.Set("tenant_id", in.TenantID)
	rec.Set("title", in.Title)
	if in.BusinessType != "" {
		rec.Set("business_type", in.BusinessType)
	}
	if in.OperatorType != "" {
		rec.Set("operator_type", in.OperatorType)
	}
	if in.Status != "" {
		rec.Set("status", in.Status)
	}
	rec.Set("method", in.Method)
	rec.Set("request_method", in.RequestMethod)
	rec.Set("oper_name", in.OperName)
	rec.Set("dept_name", in.DeptName)
	rec.Set("oper_url", in.OperURL)
	rec.Set("oper_ip", in.OperIP)
	rec.Set("oper_location", in.OperLocation)
	rec.Set("oper_param", in.OperParam)
	rec.Set("json_result", in.JSONResult)
	rec.Set("error_msg", in.ErrorMsg)
	if in.CostTime != 0 {
		rec.Set("cost_time", in.CostTime)
	}

	return e.App.Save(rec)
}
