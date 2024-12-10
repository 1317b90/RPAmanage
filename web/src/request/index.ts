import axios from "axios";
import urlJson from '@/request/url.json';


//创建axios实例
const service = axios.create({
    baseURL: urlJson.url,
    // baseURL: urlJson.local,
    timeout: 10000,//超时时间 60秒
})

//请求拦截
//就是你请求接口的时候，我会先拦截下来，对你的数据做一个判断，或者携带个token给你
service.interceptors.request.use((config) => {
    return config
}, error => {
    //处理错误请求
    return Promise.reject(error)
})

//响应拦截：后端返回来的结果
service.interceptors.response.use((res) => {
    if (res.status != 200) {
        //请求失败（包括token失效，302，404...根据和后端约定好的状态码做出不同的处理）
        return Promise.reject(res)
    } else {
        //请求成功
        return Promise.resolve(res)

    }
}, (err) => {

    // 如果是事先定义好的报错情况
    try {
        return Promise.reject(err.response.data.message)
    }
    // 如果是意外代码报错
    catch {
        return Promise.reject("请重试！")
    }
})
//因为别的地方要用，所以就把实例暴露出去，导出
export default service