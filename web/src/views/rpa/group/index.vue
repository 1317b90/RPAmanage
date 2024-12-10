<template>
  <a-table class="varTable" :columns="columns" :rowKey="(record: RPAGroupI) => record.ID" :data-source="data" style="height: 80vh;" :pagination="{ pageSize: 12 }">
    <template #headerCell="{ column }">
      <template v-if="column.dataIndex === 'operation'">
        <a-button class="editable-add-btn" @click="onAddVar" type="link">添加</a-button>
      </template>
    </template>

    <template #bodyCell="{ column,record,text }">
      <template v-if="column.dataIndex === 'IP'">
        <a class="IP-text">{{ text }}</a>
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
    <VarForm :data="modalData" :type="modalType" :getData="getData" />
    <!-- 去除默认的底部按钮 -->
    <template #footer>
    </template>
  </a-modal>

</template>
<script lang="ts" setup>
import { message } from 'ant-design-vue';
import { ref, type Ref } from 'vue';
import { h } from 'vue';
import { DeleteOutlined, FormOutlined } from '@ant-design/icons-vue';
import { getRPAGroup, delRPAGroup } from '@/request/api';
import VarForm from './form.vue';
import type { RPAGroupI } from '@/interface';

// 表格列定义
const columns = [
    {
        title: 'ID',
        dataIndex: 'ID',
        align: "center",
    },
    {
        title: '名称',
        dataIndex: 'Name',
        align: "center",
    },
    {
        title: '描述',
        dataIndex: 'Remark',
        align: "center",
    },
    { title: '虚拟机地址', dataIndex: 'IP', align: "center" },
    {
        title: '操作',
        dataIndex: 'operation',
        align: "center",
        width: '18%',
    },
];

// 数据和状态
const data: Ref<RPAGroupI[]> = ref([]);
const isModal = ref(false);
const modalData = ref<RPAGroupI>();
const modalType = ref("new");

// 获取所有数据
function getData() {
  getRPAGroup().then(res => {
    data.value = res.data.data.map((redata: any) => redata);
    isModal.value = false;
  });
}

// 初始化数据
getData();

// 点击新增
function onAddVar() {
  modalType.value = "new";
  isModal.value = true;
}

// 点击编辑
function onEdit(rowData: RPAGroupI) {
  modalType.value = "set";
  modalData.value = rowData;
  isModal.value = true;
}

// 点击删除
function onDel(rowData: RPAGroupI) {
  delRPAGroup(rowData.ID).then(() => {
    getData();
    message.success("删除成功！");
  }).catch(err => {
    message.error(err);
  });
}


</script>

<style scoped>
.delButton {
  margin-left: 10px;
}
</style>