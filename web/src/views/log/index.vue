<template>
  <div>
    <a-form
      ref="formRef"
      name="advanced_search"
      class="ant-advanced-search-form"
      :model="formState"
      @finish="onFinish"
    >
      <a-row :gutter="24">
        <a-col :span="6">
          <a-form-item name="Type" label="类型">
            <a-input v-model:value="formState.Type" placeholder="请输入类型" />
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item name="Title" label="标题">
            <a-input v-model:value="formState.Title" placeholder="请输入标题" />
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item name="Content" label="内容">
            <a-input v-model:value="formState.Content" placeholder="请输入内容" />
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item name="State" label="状态">
            <a-select v-model:value="formState.State" placeholder="请选择状态" allowClear>
              <a-select-option value="true">成功</a-select-option>
              <a-select-option value="false">失败</a-select-option>

            </a-select>
          </a-form-item>
        </a-col>
      </a-row>
      <a-row>
        <a-col :span="24" style="text-align: right">
          <a-button type="primary" html-type="submit">搜索</a-button>
          <a-button style="margin: 0 8px" @click="resetForm">重置</a-button>
        </a-col>
      </a-row>
    </a-form>

    <a-table
      class="logTable"
      :columns="columns"
      :row-key="(record: LogI) => record.ID"
      :data-source="logData"
      style="margin-top: 16px;"
      :pagination="{ pageSize: 10 }"
      :loading="loading"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'State'">
          <CheckCircleOutlined v-if="record.State == true" style="color: green;" />
          <CloseCircleOutlined v-else style="color: red;" />
        </template>
        <template v-if="column.dataIndex === 'CreatedAt'">
          {{ formatTime(record.CreatedAt) }}
        </template>
        <template v-if="column.dataIndex === 'Type'">
          <a-tag :bordered="false" :color="getColor(record.Type)">{{ record.Type }}</a-tag>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import type { FormInstance } from 'ant-design-vue';
import { getLog } from '@/request/api';
import { CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons-vue';
import { formatTime,getColor } from '@/func';


interface LogI {
  ID: number;
  CreatedAt: string;
  Type: string;
  Title: string;
  Content: string;
  State: boolean;
}

const columns = [
  {
    title: 'ID',
    dataIndex: 'ID',
    align: "center",
    width: 50
  },
  {
    title: '创建时间',
    dataIndex: 'CreatedAt',
    align: "center",
    sorter: (a: LogI, b: LogI) => new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime(),
    defaultSortOrder: 'ascend',
  },
  {
    title: '类型',
    dataIndex: 'Type',
    align: "center",
  },
  {
    title: '标题',
    dataIndex: 'Title',
    align: "center",
  },
  {
    title: '内容',
    dataIndex: 'Content',
    align: "center",
    ellipsis: true,
    width: 500
  },
  {
    title: '状态',
    dataIndex: 'State',
    align: "center",
  }
];

const formRef = ref<FormInstance>();
const formState = reactive({
  Type: '',
  Title: '',
  Content: '',
  State: '',
});
const logData = ref<LogI[]>([]);
const loading = ref(false);

const onFinish = (values: any) => {
  fetchLogData(values);
};

const resetForm = () => {
  formRef.value?.resetFields();
  fetchLogData();
};

const fetchLogData = async (params = {}) => {
  console.log(params)
  loading.value = true;
  getLog(params).then(res => {
    logData.value = res.data.data;
  }).catch(err => {
    message.error('获取日志数据失败，'+err);
  }).finally(() => {
    loading.value = false;
  });
};

onMounted(() => {
  fetchLogData();
});
</script>

<style scoped>
.ant-advanced-search-form {
  padding: 24px;
  background: #fbfbfb;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
}

.logTable {
  margin-top: 16px;
}
</style>