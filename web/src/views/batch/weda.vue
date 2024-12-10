<template>
    <div id="headerButtonDiv">
        <a-input-group compact style="width: 600px;">
        <a-select v-model:value="searchCondition" style="width: 100px">
            <a-select-option value="gsmc">公司名称</a-select-option>
            <a-select-option value="xm">法人姓名</a-select-option>
            <a-select-option value="sjhm">手机号码</a-select-option>
            <a-select-option value="_id">ID</a-select-option>
        </a-select>
        <a-input v-model:value="searchValue" style="width: 400px" @pressEnter="onSearch"/>
        <a-button type="primary" :icon="h(SearchOutlined)" @click="onSearch">搜索</a-button>
        </a-input-group>
        
        <a-button :icon="h(UndoOutlined)" @click="getData('')"/>
    </div>
    <a-table 
    :row-selection="{ selectedRowKeys: selectedRow, onChange: onSelectChange }" 
    :columns="columns"  
    :rowKey="(record: wedaI) => record._id" 
    :data-source="data" 

    :pagination="{ pageSize: 9}"
    :loading="isLoading"
    >

    </a-table>
    <a-button type="primary"  :icon="h(CheckOutlined)" style="margin-right: 10px;" @click="onConfirm">确认选中</a-button>
    <a-button  style="margin-right: 10px;" @click="onSelectAll">全选</a-button>
    <a-button @click="onCancel">取消</a-button>

</template>
<script lang="ts" setup>
import { message, type TableColumnType} from 'ant-design-vue';
import { SearchOutlined ,UndoOutlined,CheckOutlined} from '@ant-design/icons-vue';
import { getWeda } from "@/request/api"
import { ref, type Ref } from 'vue';
import { h } from 'vue';


const prop = defineProps(['returnWeda'])

const searchCondition = ref('gsmc')
const searchValue = ref('')
const isLoading = ref(false)
const selectedRow = ref<string[]>([]);

// 选中之后
const onSelectChange = (selectedRowKeys: any) => {
    selectedRow.value = selectedRowKeys;
};

interface wedaI{
    _id: string,
    dq: string,
    gsmc: string,
    sjhm: string,
    xm: string,
    dqsj: string,
}


const columns: TableColumnType<wedaI>[] = [
{
    title: 'ID',
    dataIndex: '_id',
    width: 50,
    align: "center"
},
{
    title: '公司名称',
    dataIndex: 'gsmc',
    align: "center"
},
{
    title: '法人姓名',
    dataIndex: 'xm',
    align: "center"
},
{
    title: '手机号码',
    dataIndex: 'sjhm',
    align: "center"
},
{
    title: '地区',
    dataIndex: 'dq',
    align: "center"
},
{
    title: '到期时间',
    dataIndex: 'dqsj',
    align: "center"
},
];

const data: Ref<wedaI[]> = ref([]);

// 获取数据
function getData(filter:string){
    isLoading.value = true
    getWeda(filter).then(res => {
        data.value=res.data.data
    }).catch(err=>{
        console.log(err)
        message.error("获取数据失败"+err)
    }).finally(()=>{
        isLoading.value = false
    })
}

getData("")

// 点击搜索
function onSearch(){
    const filter = "contains("+searchCondition.value+",'"+searchValue.value+"')"
    getData(filter)
}
// 确认选中,返回选中的数据
function onConfirm(){
    const selectedData = data.value.filter(item => selectedRow.value.includes(item._id));
    console.log(selectedData)
    prop.returnWeda(selectedData);
}

// 全选
function onSelectAll(){
    prop.returnWeda(data.value);
}

// 取消
function onCancel(){
    prop.returnWeda([])
}

</script>

<style scoped>
.delButton {
margin-left: 10px;
}

#headerButtonDiv{
display: flex;
flex-direction: row;
align-items: center;
}
</style>