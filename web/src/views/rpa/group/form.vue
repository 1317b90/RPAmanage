<template>
    <a-form :model="formData" name="basic" :label-col="{ span: 5 }" @finish="onFinish" >

        <a-form-item label="命名" name="Name" :rules="[{ required: true, message: '请输入命名' }]">
            <a-input v-model:value="formData.Name" />
        </a-form-item>

        <a-form-item label="描述" name="Remark">
            <a-input v-model:value="formData.Remark" />
        </a-form-item>
        <a-form-item label="虚拟机地址" name="IP">
            <a-input v-model:value="formData.IP" />
        </a-form-item>
        <a-form-item style="text-align: right;">
            <a-button type="primary" html-type="submit">提交</a-button>
        </a-form-item>
    </a-form>
</template>
<script lang="ts" setup>
import { reactive, ref } from 'vue';
import type { RPAGroupI } from "@/interface"
import { putRPAGroup } from "@/request/api"
import { message } from 'ant-design-vue';

const prop = defineProps(['data', 'type', 'getData'])

const formData = reactive<RPAGroupI>({
    ID:-1,
    Name: "",
    Remark: "",
    IP: "127.0.0.1",
});

// 如果存在传入值
if (prop.data) {
    if (prop.type == "set") {
        formData.ID = prop.data.ID
        formData.Name = prop.data.Name
        formData.Remark = prop.data.Remark
        formData.IP = prop.data.IP
    }
}

const onFinish = (values: any) => {

    putRPAGroup(formData).then(res => {
        prop.getData()
        message.success('提交成功！');
    }).catch(err => {
        message.error(err)
        })
    }


</script>
