<template>
    <a-form :model="formData" :label-col="{ span: 9 }"  @finish="onFinish" labelAlign="right"  >
        <a-form-item 
            v-for="item in allVar" 
            :key="item.VarName" 
            :name="item.VarName" 
            :label="`${item.VarRemark} ${item.VarName}`"
            :rules="[{ required: item.Required, message: '不能为空!' }]"
            >
            <a-input v-model:value="formData[item.VarName]" />

        </a-form-item>

        <a-form-item style="text-align: right;">
            <a-button type="primary" html-type="submit">提交</a-button>
        </a-form-item>
    </a-form>

 

</template>
<script lang="ts" setup>

const prop = defineProps(['allVar','formData','onSave'])


const onFinish = (values: any) => {

    prop.onSave(values)
};

</script>
