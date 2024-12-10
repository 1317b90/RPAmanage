<template>
    <a-form :model="data" name="basic" :label-col="{ span: 5 }" @finish="onFinish">
        <a-form-item label="ID" name="ID">
            <a-input v-model:value="data.ID" :bordered="false" disabled/>
        </a-form-item>
        <a-form-item label="RPA组" name="Group" :rules="[{ required: true, message: '请输入对应的RPA组' }]">
            <a-select ref="select" v-model:value="data.Group">
                <a-select-option v-for="data in prop.RPAGroupDict" :value="data.value">{{ data.text }}</a-select-option>
            </a-select>
        </a-form-item>

        <a-form-item label="RPA命名" name="Name" :rules="[{ required: true, message: '请输入命名' }]">
            <a-input v-model:value="data.Name" />
        </a-form-item>

        <a-form-item label="RPA描述" name="Remark">
            <a-input v-model:value="data.Remark" />
        </a-form-item>


        <a-form-item label="即刻执行" name="Now">
            <a-switch v-model:checked="data.Now" />
        </a-form-item>
        
        <a-form-item label="自发执行" name="Spont">
            <a-switch v-model:checked="data.Spont" />
        </a-form-item>

        <a-form-item style="text-align: right;">
            <a-button type="primary" html-type="submit">提交</a-button>
        </a-form-item>
    </a-form>
</template>
<script lang="ts" setup>
import { putRPA } from "@/request/api"
import { message } from 'ant-design-vue';

const prop = defineProps(['data','RPAGroupDict', 'getData'])

const onFinish = (values: any) => {
    putRPA(values).then(res => {
        prop.getData()
        message.success('提交成功！');
        }).catch(err => {
            message.error(err)
    })
};

</script>
