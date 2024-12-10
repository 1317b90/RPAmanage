<template>
  <!-- 选择要执行的RPA -->
  <div id="selectFuncDiv">
    <h3>选择要执行的RPA</h3>

    <a-select
      ref="select"
      v-model:value="groupValue"
      style="width: 200px;margin-right: 20px;"
      @change="getRPADictFunc()"
    >
      <a-select-option v-for="item in RPAGroupDict" :value="item.value">{{item.text}}</a-select-option>
    </a-select>

    <a-checkbox-group v-model:value="RPAList" :options="RPADict"  />

    <div id ="selectButtonDiv">
      <a-button type="primary" @click="onSelect">确定选择</a-button>
      <a-button @click="onReset">重置</a-button>
    </div>
  </div>

  <div class="content" v-if="isSelectRPA">
    <!-- 上传数据 -->
    <div class="file-upload">
      <a-button :icon="h(PlusOutlined)" class="addButton" type="primary" @click="onAdd"/>

      <a-button class="addButton" @click="isWedaFrom=true" >从优账通中获取数据</a-button>
      <a-upload
      accept=".xlsx"
      :beforeUpload="handleBeforeUpload"
      :showUploadList="false"
    >
      <a-button class="action-button">
        从本地上传文件
      </a-button>
    </a-upload>
    <a-button type="link" @click="downloadTemplate">下载模板文件</a-button>
    </div>

    <!-- 数据表格 -->
    <a-table 
    :row-selection="{ onChange: onSelectChange, selectedRowKeys: selectedRow }" 
    :row-key="(record: DataType, index: number) => index" 
    :columns="columns" 
    :data-source="tableData" 
    :expandedRowKeys="errorIndex"
    :pagination="{ pageSize: 9 }"
    :scroll="{ x: 850 }"
  >
    <template #expandedRowRender="{ index }">
      <p v-for="errorMsg in errorList[index]">
        <a-tag :bordered="false" color="processing">{{ errorMsg.name }}</a-tag>
        <a-tag :bordered="false" color="error">{{ errorMsg.msg }}</a-tag>
      </p>
    </template>
    <template #bodyCell="{ column, text }">
      <template v-if="column.dataIndex === 'name'">
        <a>{{ text }}</a>
      </template>
    </template>
  </a-table>

  <!-- 操作按钮 -->
  <div class="doButtonDiv">
    <a-button type="primary" :icon="h(CaretRightOutlined)"  @click="onExec"  :loading="isLoading">执行</a-button>
    <a-button  :icon="h(EditOutlined)" @click="onEdit"/>
    <a-button :icon="h(DeleteOutlined)" @click="onDel"/>
    <a-button :icon="h(ClearOutlined)"  @click="onClear"/>
  </div>

  </div>

  <a-modal v-if="isFrom" v-model:open="isFrom" title="Title" style="width: 650px;">
    <VarForm :allVar="allVar" :formData="formData" :onSave="onSave" />

    <!-- 去除默认的底部按钮 -->
    <template #footer>
    </template>
  </a-modal>

  <a-modal v-if="isWedaFrom" v-model:open="isWedaFrom" title="微搭数据库" style="width: 90vw;height: 60vh;">
    <wedaTable :returnWeda="returnWeda" />
    <!-- 去除默认的底部按钮 -->
    <template #footer>
    </template>
  </a-modal>

</template>

<script lang="ts" setup>
import { onMounted, ref, type Ref } from 'vue';
import { h } from 'vue';
import { message } from 'ant-design-vue';
import { DeleteOutlined, EditOutlined, PlusOutlined, ClearOutlined, CaretRightOutlined } from '@ant-design/icons-vue';
import { upFileBatch, createTask, getRPADict, getRPAGroupDict, getVarByRPANameS, downFileTemplate } from '@/request/api'
import VarForm from './form.vue'
import WedaTable from './weda.vue'
import type { filterI, varI } from '@/interface';

// 接口定义
interface DataType {
  [key: string]: string;
}

// ref 变量
const file = ref<File | null>(null);
const isLoading = ref(false);
const columns = ref<{ title: string; dataIndex: string }[]>([]);
const tableData = ref<DataType[]>([]);
const selectedRow = ref<number[]>([]);
const formData = ref();
const isFrom = ref(false);
const fromType = ref("set");
const RPAGroupDict: Ref<filterI[]> = ref([]);
const groupValue = ref();
const RPADict = ref<{value:string,label:string}[]>([]);
const RPAList = ref([]);
const errorIndex = ref<number[]>([]);
const errorList = ref<Record<number, { name: string; msg: string }[]>>({});
const isWedaFrom = ref(false);
const isSelectRPA = ref(false);
const allVar = ref<varI[]>([]);
const asNameDict = ref<Record<string, string>>({});


// RPA相关函数
function getRPADictFunc() {
  RPADict.value = [];
  getRPADict(groupValue.value).then(res => {
    for (let key in res.data.data) {
      RPADict.value.push({
        label: res.data.data[key],
        value: key
      });
    }
  }).catch(res => {
    message.error("获取RPA字典失败");
  });
}

// 点击确认选择后，获取表头
function onSelect() {
  if (groupValue.value == "" || RPAList.value.length == 0) {
    message.warning('选择不能为空！');
    return;
  }
  getVarByRPANameS(RPAList.value.toString()).then(res => {
    columns.value = [];
    allVar.value = res.data.data as varI[];
    for (let item of allVar.value) {
      if (item.AsName !== undefined && item.AsName !== '') {
        asNameDict.value[item.AsName] = item.VarName;
      }
      if (item.Required) {
        const existingColumn = columns.value.find(col => col.title === item.VarRemark);
        if (!existingColumn) {
          columns.value.push({
            title: item.VarRemark + " " + item.VarName,
            dataIndex: item.VarName,
          });
        }
      }
    }
    isSelectRPA.value = true;
  }).catch(err => {
    message.error("表头获取失败，" + err);
  });
}

// 点击重置后
function onReset() {
  groupValue.value = "";
  RPAList.value = [];
  RPADict.value = [];
  isSelectRPA.value = false;
  errorIndex.value = [];
  errorList.value = {};
}

// 点击手动添加数据
const onAdd = () => {
  formData.value = allVar.value.reduce((obj:any, column:varI) => {
    obj[column.VarName] = '';
    return obj;
  }, {});
  isFrom.value = true;
  fromType.value = "new";
};

// 点击下载模板
function downloadTemplate() {
  console.log(allVar.value); 
  downFileTemplate(allVar.value).then(res => {
    window.open(res.data.url);
  }).catch(err => {
    message.error("下载模板失败，" + err);
  });
}

// 上传文件之前
const handleBeforeUpload = (fileObj: File) => {
  file.value = fileObj;
  uploadFile();
  return false;
};

// 上传文件函数
const uploadFile = async () => {
  if (!file.value) {
    message.error('请选择文件');
    return;
  }
  const formData = new FormData();
  formData.append('file', file.value);
  try {
    const res = await upFileBatch(formData);
    message.success('上传成功');
    tableData.value.push(...res.data.data);
  } catch (err) {
    message.error(`上传失败: ${err}`);
    console.error('上传失败', err);
  }
};

// 选中表格数据时
const onSelectChange = (selectedRowKeys: number[]) => {
  selectedRow.value = selectedRowKeys;
};

// 表格数据新增或修改时
const onSave = (newData:DataType) => {
  if (fromType.value == "set") {
    tableData.value[selectedRow.value[0]] = newData;
    message.success('保存成功！');
  } else {
    tableData.value.push(newData);
    message.success('新增成功！');
  }

  isFrom.value = false;
};


// 点击提交任务
const onExec = async () => {
  
  if (RPAList.value.length === 0) {
    message.error('请选择功能');

  }  else if (selectedRow.value.length === 0 && columns.value.length != 0) {
    message.warning('请选择要执行的数据');
  } else if (tableData.value.length === 0 && columns.value.length != 0) {
    message.error('表格数据为空，请先上传数据');
  }  else {
    isLoading.value = true;
    errorIndex.value = [];
    errorList.value = {};
    let delIndex = [];

    // 如果列头为空，则说明不需要输入，可直接执行
    if (columns.value.length === 0) {
      selectedRow.value = [0];
      tableData.value=[{}]
    }

    for (let i of selectedRow.value) {
      const rowData = tableData.value[i];
      const rowErrorMsg: { name: string; msg: string }[] = [];

      for (const RPAName of RPAList.value) {
        try {
          await createTask(RPAName, rowData);
          delIndex.push(i);
        } catch (err:any) {
          errorIndex.value.push(i);
          rowErrorMsg.push({ name: RPAName, msg: err.toString() });
        }
      }

      errorList.value[i] = rowErrorMsg;
    }

    delIndex = delIndex.sort((a, b) => b - a);
    tableData.value = tableData.value.filter((_, index) => !delIndex.includes(index));

    selectedRow.value = [];
    if (errorIndex.value.length > 0) {
      message.warning('部分任务添加失败，请检查错误信息！');
    } else {
      message.success('任务添加完毕！');
    }
    isLoading.value = false;
  }
};

// 编辑表格数据
const onEdit = () => {
  if (selectedRow.value.length === 0) {
    message.warning('请选择要编辑的数据');
  } else if (selectedRow.value.length > 1) {
    message.warning('只能选择一条数据进行编辑');
  } else {
    const index = selectedRow.value[0];
    const rowData = tableData.value[index];
    formData.value = rowData;
    isFrom.value = true;
    fromType.value = "set";
  }
};

// 删除表格数据
const onDel = () => {
  if (selectedRow.value.length === 0) {
    message.warning('请选择要删除的数据');
    return;
  }
  const indexesToDelete = new Set(selectedRow.value.sort((a, b) => b - a));
  tableData.value = tableData.value.filter((_, index) => !indexesToDelete.has(index));
  selectedRow.value = [];
  message.success(`成功删除 ${indexesToDelete.size} 条数据`);
};

// 清空表格数据
const onClear = () => {
  tableData.value = [];
  selectedRow.value = [];
  message.success('清空成功！');
};

// 获取微搭数据的回调
function returnWeda(data: any) {
  isWedaFrom.value = false;
  errorIndex.value = [];
  errorList.value = {};
  
  const transformedData = data.map((item:any) => {
    const newItem:any = {};
    for (const [oldKey, newKey] of Object.entries(asNameDict.value)) {
      if (oldKey in item) {
        newItem[newKey] = item[oldKey];
      }
    }
    return newItem;
  });
  
  tableData.value.push(...transformedData);
}

// 初始化数据
onMounted(() => {
  getRPAGroupDict('').then(res => {
    for (let key in res.data.data) {
      RPAGroupDict.value.push({
        text: res.data.data[key],
        value: key
      });
    }
  });
});

</script>

<style scoped>
#selectFuncDiv{
  background-color: #f7f7f7;
  padding: 10px 10px 10px 30px;
  border-radius: 4px;
  
}

#selectFuncDiv h3{
  color: #1890ff;
  font-weight: bold;
  margin-bottom: 10px;
}

#selectButtonDiv{
  margin-top: 10px;
  display: flex;
  justify-content: flex-end;
  gap: 15px;
}
.addButton {
  margin-right: 20px;
}

.file-upload {
  margin-top: -30px;
  margin-bottom: 10px;
  padding-bottom: 20px;
  padding-left: 20px;
  display: flex;
  align-items: center;
  background-color: #f5f5f5;
}
.file-label {
  display: inline-block;
  padding: 4px 15px;
  background-color: #1890ff;
  color: white;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}
.file-label:hover {
  background-color: #40a9ff;
}
.file-input {
  display: none;
}
.file-name {
  margin-left: 10px;
  font-size: 14px;
}
.doButtonDiv{
  margin-left: 20px;
  margin-top: 20px;
}

.doButtonDiv > * {
  margin-right: 10px;
}


</style>