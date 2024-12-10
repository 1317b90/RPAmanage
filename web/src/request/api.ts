import service from "@/request/index";


// --------- 任务 --------- 任务 --------- 任务 --------- 任务 --------- 任务 --------- 任务 --------- 任务

// 从内存获取任务
export async function getTaskMemory() {
    return service({
        url: "task/memory",
        method: "GET",
    })
}

// 从数据库获取任务
export async function getTaskDB(state:string) {
    if (state=="all"){
        return service({
            url: "task/db",
            method: "GET",
        })
    }else{
        return service({
        url: "task/db?state="+state,
        method: "GET",
    })}

}

// 删除内存任务
export async function delTaskMemory(id:string) {
    return service({
        url: "task/memory?id="+id,
        method: "Delete",
    })
}

// 创建任务
export async function createTask(RPAName: string, Input: object) {

    return service({
        url: "task",
        method: "POST",
        data: {
            "RPAName": RPAName,
            "Input": Input
        }
    })
}

// 统计任务数据
export async function countTask() {
    return service({
        url: "task/count",
        method: "GET",
    })
}
// --------- 变量 --------- 变量 --------- 变量 --------- 变量 --------- 变量 --------- 变量 --------- 变量

// 获取指定任务的变量
export async function getVar(RPAName: string) {
    return service({
        url: "var?RPAName=" + RPAName,
        method: "GET",
    })
}


// 获取所有变量
export async function getVars(RPAGroupName:string,RPAName:string) {
    return service({
        url: "var?RPAGroupName="+RPAGroupName+"&RPAName="+RPAName,
        method: "GET",
    })
}

// 修改或新增变量
export async function putVar(data: any) {
    return service({
        url: "var",
        method: "PUT",
        data: data
    })
}


// 删除变量
export async function delVar(id: number) {
    return service({
        url: "var?id=" + String(id),
        method: "Delete",
    })
}


// --------- RPA组 --------- RPA组 --------- RPA组 --------- RPA组 --------- RPA组 --------- RPA组 --------- RPA组
export async function getRPAGroup(){
    return service({
        url: "rpa/group",
        method: "GET",
    })
}

export async function putRPAGroup(data: any) {
    return service({
        url: "rpa/group",
        method: "PUT",
        data: data
    })
}

export async function delRPAGroup(id: number) {
    return service({
        url: "rpa/group?id=" + String(id),
        method: "DELETE",
    })
}


export async function getRPAGroupDict(name:string) {    
    return service({
        url: "rpa/group/dict?name=" + name,
        method: "GET",
    })
}

// --------- RPA表 --------- RPA表 --------- RPA表 --------- RPA表 --------- RPA表 --------- RPA表 --------- RPA表

// 增加和修改RPA 
export async function putRPA(data: any) {
    return service({
        url: "rpa",
        method: "PUT",
        data: data
    })
}



// 获取RPA名称字典
export async function getRPADict(group:string) {
    return service({
        url: "rpa/dict?&group="+group,
        method: "GET",
    })
}


// 获取所有的RPA
export async function getRPA(group:string) {
    return service({
        url: "rpa?group="+group,
        method: "GET",
    })
}

// 删除RPA
export async function delRPA(id: number) {
    return service({
        url: "rpa?id=" + String(id),
        method: "DELETE",
    })
}

// 根据输入的多个RPAName获取变量
export async function getVarByRPANameS(RPANameList:string) {
    return service({
        url: "var/rpa?RPANameList=" + RPANameList,
        method: "GET",
    })
}

// --------- 文件 --------- 文件 --------- 文件 --------- 文件 --------- 文件 --------- 文件 --------- 文件

// 亿企代账上传文件
export async function upFileBatch(formData:FormData){
    return service({
        url: "upfile/batch",
        method: "POST",
        data: formData,
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
}

// 下载模板文件
export async function downFileTemplate(head:object){
    return service({
        url: "downfile/template",
        method: "POST",
        data: head
    })
}
// 上传通用文件
export async function upFileCommon(formData:FormData){
    return service({
        url: "upfile/common",
        method: "POST",
        data: formData,
    })
}

// --------- 亿企代账 --------- 亿企代账 --------- 亿企代账 --------- 亿企代账 --------- 亿企代账 --------- 亿企代账 --------- 亿企代账

// 获取亿企代账账号列表
export async function getYqdzAccount() {
    return service({
        url: "yqdz/account",
        method: "GET",
    })
}

// 选择亿企代账账号
export async function useYqdzAccount(id: number) {
    return service({
        url: "yqdz/useAccount?id=" + String(id),
        method: "GET",
    })
}


// --------- 微搭 --------- 微搭 --------- 微搭 --------- 微搭 --------- 微搭 --------- 微搭 --------- 微搭


// 获取微搭数据
export async function getWeda(filter:string){
    return service({
        url: "weda?filter="+filter,
        method: "GET",
    })
}

// --------- 定时任务 --------- 定时任务 --------- 定时任务 --------- 定时任务 --------- 定时任务 --------- 定时任务 --------- 定时任务

// 获取所有定时任务
export async function getCron(){
    return service({
        url: "cron",
        method: "GET",
    })
}

// 创建定时任务
export async function addCron(data: object) {
    return service({
        url: "cron",
        method: "POST",
        data: data
    })
}

// 删除定时任务
export async function delCron(id: string) {
    return service({
        url: "cron?id=" + id,
        method: "DELETE",
    })
}

// --------- 日志 --------- 日志 --------- 日志 --------- 日志 --------- 日志 --------- 日志 --------- 日志

// 获取日志
export async function getLog(params:object){
    return service({
        url: "log",
        method: "GET",
        params: params
    })
}


// --------- 企业微信 --------- 企业微信 --------- 企业微信 --------- 企业微信 --------- 企业微信 --------- 企业微信 --------- 企业微信

// 发送消息
export async function sendWecomMessage(data: object) {
    return service({
        url: "wecom",
        method: "POST",
        data: data
    })
}


// 获取消息队列
export async function getWecomMessageList() {
    return service({
        url: "wecom/list",
        method: "GET",
    })
}

// 删除消息
export async function delWecomMessage(id: string) {
    return service({
        url: "wecom/" + id,
        method: "DELETE",
    })
}

