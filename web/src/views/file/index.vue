<template>
  <div class="upload-container">
    <a-upload-dragger
      v-model:fileList="fileList"
      name="file"
      :multiple="false"
      :maxCount="1"
      :customRequest="customRequest"
      @preview="handlePreview"
      @change="handleChange"
      @drop="handleDrop"
    >
      <p class="ant-upload-drag-icon">
        <inbox-outlined></inbox-outlined>
      </p>
      <p class="ant-upload-text">点击或拖拽文件到此区域上传</p>
    </a-upload-dragger>


  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { InboxOutlined } from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import type { UploadChangeParam } from 'ant-design-vue';
import { upFileCommon } from '@/request/api';

const fileList = ref([]);

const customRequest = async (options: any) => {
  const { file, onSuccess, onError } = options;
  const formData = new FormData();
  formData.append('file', file);
  
  try {
    const res = await upFileCommon(formData);
    file.url = res.data.url;
    onSuccess(res.data.url);
  } catch (err) {
    onError(err);
  }
};

const handleChange = (info: UploadChangeParam) => {
  const status = info.file.status;
  if (status !== 'uploading') {
    console.log(info.file, info.fileList);
  }
  if (status === 'done') {
    message.success(`${info.file.name} 文件上传成功`);
  } else if (status === 'error') {
    message.error(`${info.file.name} 文件上传失败`);
  }
};

function handleDrop(e: DragEvent) {
  console.log(e);
}


const handlePreview = (file: any) => {
  navigator.clipboard.writeText(file.url);
  window.open(file.url, '_blank');
  message.success('复制成功');
}


</script>

<style scoped>
.upload-container {
  width: 500px;
  height: 150px;
  margin: 0 auto;
}

:deep(.ant-upload-drag) {
  height: 150px;
}


</style>
