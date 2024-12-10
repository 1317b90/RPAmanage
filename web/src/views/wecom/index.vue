<template>
  <div>
    <div class="table-operations">
      <a-button type="primary" @click="showModal">新增消息</a-button>
      <a-button style="margin-left: 8px" @click="fetchMessageList">刷新</a-button>
    </div>
    <a-table :columns="columns" :data-source="messageList" :rowKey="(record: FormState) => record.ID" >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'operation'">
          <a-button size="small" :icon="h(DeleteOutlined)" @click="delMessage(record.ID)" />
        </template>
      </template>

    </a-table>

    <a-modal
      v-model:open="visible"
      title="发送消息"
      @ok="handleOk"
      @cancel="handleCancel"
      :confirmLoading="confirmLoading"
    >
      <a-form
        :model="formData"
        name="basic"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 16 }"
      >
        <a-form-item
          label="用户ID"
          name="userID"
          :rules="[{ required: true, message: '请输入用户ID!' }]"
        >
          <a-input v-model:value="formData.userID" />
        </a-form-item>

        <a-form-item
          label="任务ID"
          name="taskID"
          :rules="[{ required: true, message: '请输入任务ID!' }]"
        >
          <a-input v-model:value="formData.taskID" />
        </a-form-item>

        <a-form-item
          label="群组"
          name="Group"
          :rules="[{ required: true, message: '请输入群组!' }]"
        >
          <a-input v-model:value="formData.Group" />
        </a-form-item>

        <a-form-item label="类型" name="Type">
          <a-input v-model:value="formData.Type" />
        </a-form-item>

        <a-form-item label="内容" name="Content">
          <a-input v-model:value="formData.Content" />
        </a-form-item>

        <a-form-item label="数据" name="Data">
          <a-input v-model:value="formData.Data" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import { getWecomMessageList, delWecomMessage, sendWecomMessage } from '@/request/api';
import { message } from 'ant-design-vue';
import { DeleteOutlined } from '@ant-design/icons-vue';
import { h } from 'vue';

interface FormState {
  ID?: string;
  userID: string;
  taskID: string;
  Group: string;
  Type: string;
  Content: string;
  Data: string;
}

const visible = ref<boolean>(false);
const confirmLoading = ref<boolean>(false);
const messageList = ref<FormState[]>([]);

const formData = reactive<FormState>({
  userID: '',
  taskID: '',
  Group: '',
  Type: '',
  Content: '',
  Data: '',
});

const columns = [
  {
    title: '用户ID',
    dataIndex: 'userID',
    key: 'userID',
    width: 100,
    align: 'center',
  },
  {
    title: '任务ID',
    dataIndex: 'taskID',
    key: 'taskID',
    width: 100,
    align: 'center',
  },
  {
    title: '群组',
    dataIndex: 'Group',
    key: 'Group',
    width: 100,
    align: 'center',
  },
  {
    title: '类型',
    dataIndex: 'Type',
    key: 'Type',
    width: 100,
    align: 'center',
  },
  {
    title: '内容',
    dataIndex: 'Content',
    key: 'Content',
    align: 'center',
  },
  {
    title: '数据',
    dataIndex: 'Data',
    key: 'Data',
    align: 'center',
  },
  {
    title: '操作',
    key: 'operation',
    dataIndex: 'operation',
    width: 100,
    align: 'center',
  },
];

const fetchMessageList = async () => {
  try {
    const res = await getWecomMessageList();
    messageList.value = res.data;
    console.log(messageList.value);
  } catch (err: any) {
    message.error(err.message);
  }
};

const showModal = () => {
  visible.value = true;
};

const handleOk = () => {
  confirmLoading.value = true;
  sendWecomMessage(formData)
    .then(() => {
      message.success('发送成功');
      visible.value = false;
      fetchMessageList();
      resetForm();
    })
    .catch((err) => {
      message.error(err.message);
    })
    .finally(() => {
      confirmLoading.value = false;
    });
};

const handleCancel = () => {
  visible.value = false;
  resetForm();
};

const resetForm = () => {
  formData.userID = '';
  formData.taskID = '';
  formData.Group = '';
  formData.Content = '';
  formData.Type = '';
  formData.Data = '';
};

const delMessage = (id?: string) => {
  if (!id) return;
  delWecomMessage(id)
    .then(() => {
      fetchMessageList();
      message.success('删除成功');
    })
    .catch((err: any) => {
      message.error(err.message);
    });
};

onMounted(() => {
  fetchMessageList();
});
</script>

<style scoped>
.table-operations {
  margin-bottom: 16px;
}
</style>
