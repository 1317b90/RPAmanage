<template>
    <a-form :model="data" name="basic" :label-col="{ span: 5 }" @finish="onFinish">
        <a-form-item label="ID" name="ID">
            <a-input v-model:value="data.ID" :bordered="false" disabled/>
        </a-form-item>
        <a-form-item label="RPA组" name="RPAGroup" :rules="[{ required: true, message: '请输入对应的RPA组' }]">
            <a-select ref="select" v-model:value="data.RPAGroup">
                <a-select-option v-for="data in prop.RPAGroupDict" :value="data.value">{{ data.text }}</a-select-option>
            </a-select>
        </a-form-item>

        <a-form-item label="RPA名称" name="RPAName" :rules="[{ required: true, message: '请输入对应RPA名称' }]">
            <a-select ref="select" v-model:value="data.RPAName">
                <a-select-option v-for="data in prop.RPANameDict" :value="data.value">{{ data.text }}</a-select-option>
            </a-select>
        </a-form-item>

        <a-form-item label="变量命名" name="VarName" :rules="[{ required: true, message: '请输入变量名称' }]">
            <a-input v-model:value="data.VarName" />
        </a-form-item>

        <a-form-item label="变量描述" name="VarRemark" :rules="[{ required: true, message: '请输入变量描述' }]">
            <a-input v-model:value="data.VarRemark" />
        </a-form-item>

        <a-form-item label="变量别名" name="AsName">
            <a-input v-model:value="data.AsName" />
        </a-form-item>

        <a-form-item label="默认值" name="Default">
            <a-input v-model:value="data.Default" />
        </a-form-item>

        <a-form-item label="数据类型" name="VarType">
            <a-select ref="select" v-model:value="data.VarType">
                <a-select-option value="string">string</a-select-option>
                <a-select-option value="int">int</a-select-option>
                <a-select-option value="float">float</a-select-option>
                <a-select-option value="json">json</a-select-option>
                <a-select-option value="data">data</a-select-option>
            </a-select>
        </a-form-item>

        <a-form-item label="验证类型" name="VeriftType">
            <a-select ref="select" v-model:value="data.VerifyType">
                <a-select-option value="">无</a-select-option>
                <a-select-option value="手机号">手机号</a-select-option>
                <a-select-option value="身份证号">身份证号</a-select-option>
            </a-select>

        </a-form-item>

        <a-form-item label="是否必填" name="Required">
            <a-switch v-model:checked="data.Required" />
        </a-form-item>

        <a-form-item style="text-align: right;">
            <a-button type="primary" html-type="submit">提交</a-button>
        </a-form-item>
    </a-form>
</template>
<script lang="ts" setup>
import { putVar } from "@/request/api"
import { message } from 'ant-design-vue';

const prop = defineProps(['data', 'RPANameDict', 'RPAGroupDict', 'getData'])


const onFinish = (values: any) => {
    console.log(values)
    putVar(values).then(res => {
        prop.getData()
        message.success('提交成功！');
    }).catch(err => {
        message.error(err)
    })

};

</script>
