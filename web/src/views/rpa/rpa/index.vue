<template>

  <a-radio-group v-model:value="groupValue" @change="getData" button-style="solid">
    <a-radio-button value="">所有组</a-radio-button>
    <a-radio-button v-for="item in  RPAGroupDict" :value="item.value">{{ item.text }}</a-radio-button>
    
  </a-radio-group>

  <a-table class="varTable" :columns="columns" :rowKey="(record: RPAI) => record.ID" :data-source="data" style="height: 80vh;" :pagination="{ pageSize: 12 }">
    <template #headerCell="{ column }">
      <template v-if="column.dataIndex === 'operation'">
        <a-button class="editable-add-btn" @click="onAddVar" type="link">添加</a-button>
      </template>
    </template>

    <template #bodyCell="{ column,record,text }">
      <template v-if="column.dataIndex === 'Now'">
        <a-switch v-model:checked="record.Now"  size="small" disabled/>
      </template>
      <template v-if="column.dataIndex === 'Spont'">
        <a-switch v-model:checked="record.Spont"  size="small" disabled/>
      </template>
      <template v-if="column.dataIndex === 'Name'">
        <a-tag :bordered="false" :color="getColor(text)">{{ text }}</a-tag>
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
    <VarForm :data="modalData" :getData="getData" :RPAGroupDict="RPAGroupDict" />
    <!-- 去除默认的地步按钮 -->
    <template #footer>
    </template>
  </a-modal>

</template>
<script lang="ts" setup>
import { message } from 'ant-design-vue';
import { onMounted, ref, type Ref } from 'vue';
import { h } from 'vue';
import { DeleteOutlined, FormOutlined } from '@ant-design/icons-vue';
import { getRPA, delRPA, getRPAGroupDict } from '@/request/api';
import VarForm from './form.vue';
import type { RPAI, filterI } from '@/interface';
import { getColor } from "@/func";

// 状态和引用
const data: Ref<RPAI[]> = ref([]);
const RPAGroupDict: Ref<filterI[]> = ref([]);
const isModal = ref(false);
const modalData = ref<RPAI>();
const groupValue=ref("")

// 表格列定义
const columns = [
    { title: 'ID', dataIndex: 'ID', align: "center" },
    { title: '组', dataIndex: 'Group', align: "center" },
    { title: '名称', dataIndex: 'Name', align: "center" },
    { title: '描述', dataIndex: 'Remark', align: "center" },
    { title: '即刻执行', dataIndex: 'Now', align: "center" },
    { title: '自发执行', dataIndex: 'Spont', align: "center" },
    { title: '操作', dataIndex: 'operation', align: "center", width: '18%' },
];


// 获取所有的RPA数据
const getData = () => {
    getRPA(groupValue.value).then(res => {
        data.value = res.data.data;
        isModal.value = false;
    });
};

// 点击新增RPA
const onAddVar = () => {
    modalData.value = {
        ID: -1,
        Group: groupValue.value,
        Name: "",
        Remark: "",
        Now: false,
        Spont: false,
    };
    isModal.value = true;
};

// 点击编辑RPA
const onEdit = (rowData: RPAI) => {
    modalData.value = rowData;
    isModal.value = true;
};

// 点击删除RPA
const onDel = (rowData: RPAI) => {
    delRPA(rowData.ID).then(() => {
        getData();
        message.success("删除成功！");
    }).catch(err => {
        message.error(err);
    });
};

// 初始化
onMounted(() => {
  getRPAGroupDict('').then(res => {
    for (let key in res.data.data) {
      RPAGroupDict.value.push({
        text: res.data.data[key],
        value: key
      }
      )
    }
  })
  getData();
})

</script>

<style scoped>
.delButton {
  margin-left: 10px;
}
.IP-text{
  text-decoration: underline;
}
</style>