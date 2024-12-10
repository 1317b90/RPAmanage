<template>
  <a-table
    class="varTable"
    :columns="columns"
    :rowKey="(record: TableRecord) => record.ID"
    :data-source="data"
    style="height: 80vh;"
    :pagination="{ pageSize: 12 }"
  >
    <template #headerCell="{ column }">
      <template v-if="column.key === 'ID'">
        任务ID
      </template>
      <template v-if="column.dataIndex === 'operation'">
        <a-button class="editable-add-btn" @click="onAddCron" type="link">添加定时任务</a-button>
      </template>
    </template>
    <template #bodyCell="{ column, record }: { column: any; record: TableRecord }">


      <template v-if="column.dataIndex === 'operation'">
        <a-button type="primary" size="small" :icon="h(FormOutlined)" @click="onEdit(record)" />
        <a-popconfirm
          title="确定删除?"
          placement="leftTop"
          ok-text="确定"
          cancel-text="取消"
          @confirm="onDel(record)"
        >
          <a-button size="small" :icon="h(DeleteOutlined)" style="margin-left: 10px;" />
        </a-popconfirm>
      </template>
    </template>
  </a-table>

  <a-modal v-if="isFrom" v-model:open="isFrom" title="新建定时任务" style="width: 600px;">
    <a-form :model="formData" name="basic" @finish="onSubmitCron" :label-col="{ span: 5 }"> 
      <a-form-item
        label="RPA名称"
        name="RPAName"
        :rules="[{ required: true, message: '请选择要执行的RPA' }]"
      >
        <a-select ref="select" v-model:value="formData.RPAName">
          <a-select-option v-for="data in RPANameDict" :value="data.value">{{ data.text }}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item
        label="任务描述"
        name="Name"
        :rules="[{ required: true, message: '请输入任务描述' }]"
      >
        <a-input v-model:value="formData.Name" />
      </a-form-item>

      <a-form-item
        label="执行方式"
        name="Type"
      >
        <a-select ref="select" v-model:value="formData.Type">
          <a-select-option v-for="data in typeDict" :value="data.value">{{ data.text }}</a-select-option>
        </a-select>
      </a-form-item>


      <a-form-item
        label="执行规则"
        name="CronData"
      >
        <a-input v-model:value="formData.CronData" />

        <a v-if="formData.Type === 'cron'">* * * * *  分、时、日、月、周</a>
        <a v-else-if="formData.Type === 'month'">请输入在每个月几号执行</a>
        <a v-else>留空即可</a>


      </a-form-item>

      <a-form-item
        label="任务输入"
        name="Input"
      >
        <a-input v-model:value="formData.Input" />
      </a-form-item>



      <a-form-item style="text-align: right;">
        <a-button type="primary" html-type="submit">提交</a-button>
      </a-form-item>
    </a-form>
    <!-- 去除默认的底部按钮 -->
    <template #footer>
    </template>
  </a-modal>
</template>

<script lang="ts" setup>
import { addCron, delCron, getCron, getRPADict } from '@/request/api'
import { onMounted, ref, type Ref } from 'vue';
import { h } from 'vue';
import { DeleteOutlined, FormOutlined, SmileOutlined } from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import type { filterI } from '@/interface';

const typeDict = [
  {
    text: "cron",
    value: "cron"
  },
  {
    text: "每月执行",
    value: "month"
  },
  {
    text: "只执行一次",
    value: "once"
  }
]


interface TableRecord {
  ID: string;
  Name: string;
  RPAName: string;
  Type: string;
  CronData: string;
  NextRun: string;
}

const columns = [
  {
    title: '任务ID',
    dataIndex: 'ID',
    align: "center",
    key: 'ID'
  },
  {
    title: 'RPA名称',
    dataIndex: 'RPAName',
    key: 'RPAName',
    align: "center"
  },
  {
    title: '任务描述',
    dataIndex: 'Name',
    key: 'Name',
    align: "center"
  },
  {
    title: '执行方式',
    dataIndex: 'Type',
    key: 'Type',
    align: "center"
  },
  {
    title: '执行规则',
    dataIndex: 'CronData',
    key: 'CronData',
    align: "center"
  },
  {
    title: '下次执行时间',
    dataIndex: 'NextRun',
    key: 'NextRun',
    align: "center"
  },
  {
    title: '操作',
    dataIndex: 'operation',
    fixed: 'right',
    align: "center"
  },
];

const data = ref([])

// 获取定时任务
function getData() {
  getCron().then(res => {
    data.value = res.data.data
    console.log(data.value)

  }).catch(err => {
    message.error("获取定时任务失败！", err)
  })
}

onMounted(async () => {
  getData()
})

function onEdit(record: TableRecord) {
  console.log(record)
}

function onDel(record: TableRecord) {
  delCron(record.ID).then(res => {
    getData()
    message.success("删除定时任务成功！")
  }).catch(err => {
    message.error("删除定时任务失败！" + err)
  })
}

function onAddCron() {
  isFrom.value = true
}

const isFrom = ref(false)
const formData = ref({
  RPAName: '',
  Name: '',
  Type: 'cron',
  CronData: '',
  Input:"",
})

// 用于筛选的列表
const RPANameDict: Ref<filterI[]> = ref([])

// 获取任务名称字典
getRPADict("").then(res => {
  for (let key in res.data.data) {
    RPANameDict.value.push({
      text: res.data.data[key],
      value: key
    })
  }
})

function onSubmitCron(values: any) {
  addCron(values).then(res => {
    getData()
    message.success("添加定时任务成功！")
    isFrom.value = false
  }).catch(err => {
    message.error("添加定时任务失败！" + err)
    isFrom.value = false
  })
}
</script>