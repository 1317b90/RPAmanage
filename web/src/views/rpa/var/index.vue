<template>
  <a-radio-group v-model:value="groupValue" @change="changeGroup" button-style="solid">
    <a-radio-button value="">所有组</a-radio-button>
    <a-radio-button v-for="item in  RPAGroupDict" :value="item.value">{{ item.text }}</a-radio-button>
    
  </a-radio-group>
  <br><br>
  <a-radio-group v-model:value="rpaValue" @change="getData" button-style="solid" >
    <a-radio-button v-for="item in  RPANameDict" :value="item.value">{{ item.text }}</a-radio-button>
  </a-radio-group>

  <br>
  <a-table class="varTable" :columns="columns"  :rowKey="(record: varI) => record.ID" :data-source="data" style="height: 70vh;" :pagination="{ pageSize: 10 }">
    <template #headerCell="{ column }">
      <template v-if="column.dataIndex === 'operation'">
        <a-button class="editable-add-btn" @click="onAddVar" type="link">添加</a-button>
      </template>
    </template>

    <template #bodyCell="{ column, record}">
      <template v-if="column.dataIndex === 'Required'">
        <a-switch v-model:checked="record.Required"  size="small" disabled/>
      </template>
 
      <template v-if="column.dataIndex === 'operation'">
        <a-button type="primary" size="small" :icon="h(FormOutlined)" @click="onEdit(record)" />
        <a-popconfirm title="确定删除?" placement="leftTop" ok-text="确定" cancel-text="取消" @confirm="onDel(record)">

          <a-button size="small" :icon="h(DeleteOutlined)" class="delButton" />
        </a-popconfirm>
      </template>
    </template>
  </a-table>

  <!-- 添加v-if 让每一次打开都重新初始化 -->

  <a-modal v-if="isModal" v-model:open="isModal" title="Title">
    <VarForm :data="modalData"  :RPANameDict="RPANameDict" :RPAGroupDict="RPAGroupDict" :getData="getData" />
    <!-- 去除默认的地步按钮 -->
    <template #footer>
    </template>
  </a-modal>

</template>
<script lang="ts" setup>
import { message, type TableColumnType } from 'ant-design-vue';
import { onMounted, ref, type Ref } from 'vue';
import { h } from 'vue';
import { DeleteOutlined, FormOutlined } from '@ant-design/icons-vue';
import { getVars, getRPADict, getRPAGroupDict, delVar } from "@/request/api";
import VarForm from './form.vue';
import type { filterI, varI } from '@/interface';

// 状态和引用
const data: Ref<varI[]> = ref([]);
const RPANameDict: Ref<filterI[]> = ref([]);
const RPAGroupDict: Ref<filterI[]> = ref([]);
const isModal = ref(false);
const modalData = ref<varI>();
const groupValue=ref("")
const rpaValue=ref("")

// 表格列定义
const columns: TableColumnType<varI>[] = [
  { title: 'ID', dataIndex: 'ID', width: 50, align: "center" },
  { title: '变量名称', dataIndex: 'VarName', align: "center" },
  { title: '变量描述', dataIndex: 'VarRemark', align: "center" },
  { title: '变量别名', dataIndex: 'AsName', align: "center" },
  { title: '默认值', dataIndex: 'Default', align: "center" },
  { title: '变量数据类型', dataIndex: 'VarType', align: "center" },
  { title: '变量验证类型', dataIndex: 'VeriftType', align: "center" },
  { title: '是否必填', dataIndex: 'Required', align: "center" },
  { title: '操作', dataIndex: 'operation', fixed: 'right', width: 100, align: "center" },
];

// 获取所有的变量
const getData = () => {
  getVars(groupValue.value,rpaValue.value).then(res => {
    data.value = res.data.data;
    isModal.value = false;
  });
};

// 点击新增变量
const onAddVar = () => {
  isModal.value = true;
  modalData.value = {
    ID: -1,
    VarName: "",
    VarRemark: "",
    RPAGroup: groupValue.value,
    RPAName: rpaValue.value,
    AsName: "",
    VarType: "string",
    VerifyType: "",
    Default: "",
    Required: false,
  }
};

// 点击编辑变量
const onEdit = (rowData: varI) => {
  modalData.value = rowData;
  isModal.value = true;
};

// 点击删除变量
const onDel = (rowData: varI) => {
  delVar(rowData.ID || -1).then(() => {
    getData();
    message.success("删除成功！");
  }).catch(err => {
    message.error(err);
  });
};

// 改变组
const changeGroup = ()=>{
  getRPADictFunc(groupValue.value)
 
  
}

// 根据RPA组获取RPA的字典
function getRPADictFunc(group:string){

  RPANameDict.value=[]
  getRPADict(group).then(res => {
    for (let key in res.data.data) {
      RPANameDict.value.push({
        text: res.data.data[key],
        value: key
      }
      )
    }
    rpaValue.value=RPANameDict.value[0].value
    getData()
  }
)
}

// 初始化
onMounted(() => {
  getRPADictFunc("")

  getRPAGroupDict('').then(res => {
    for (let key in res.data.data) {
      RPAGroupDict.value.push({
        text: res.data.data[key],
        value: key
      }
      )
    }
  })
})


</script>

<style scoped>
.delButton {
  margin-left: 10px;
}
</style>