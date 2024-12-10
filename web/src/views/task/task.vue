<template>
   <dashView :getData="getData" />
   <div class="cards">
    <template v-for="(value, key) in cardColor" :key="key">
      <div class="card">
        <a-button type="text" class="card-title" @click="getData(key)">{{ value.title }}</a-button>
        <p class="card-number" :style="{ color: value.color }">{{ cardData[key] || 0 }}</p>
      </div>
    </template>
  </div>
   <br>
  <a-table class="taskTable" :rowKey="(record: taskI) => record.ID" :columns="columns" :data-source="data"  :pagination="{ pageSize: 10 }">

    <template
      #customFilterDropdown="{ setSelectedKeys, selectedKeys, confirm, clearFilters, column }"
    >
      <div style="padding: 8px">
        <a-input
          ref="searchInput"
          :placeholder="`Search ${column.dataIndex}`"
          :value="selectedKeys[0]"
          style="width: 188px; margin-bottom: 8px; display: block"
          
          @change="(e: Event) => setSelectedKeys((e.target as HTMLInputElement).value ? [(e.target as HTMLInputElement).value] : [])"
          @pressEnter="handleSearch(selectedKeys, confirm, column.dataIndex)"
        />
        <a-button
          type="primary"
          size="small"
          style="width: 90px; margin-right: 8px"
          @click="handleSearch(selectedKeys, confirm, column.dataIndex)"
        >
          <template #icon><SearchOutlined /></template>
          搜索
        </a-button>
        <a-button size="small" style="width: 90px" @click="handleReset(clearFilters)">
          重置
        </a-button>
      </div>
    </template>

    <template #expandedRowRender="{ record }">

      <div class="container">
        <div class="left-column">
          <a-text>输入数据：</a-text>
          <p style="margin: 0" v-if="getJSON(record.Input)==''">
        {{ record.Input }}
        </p>
         <json-viewer v-else :value="getJSON(record.Input)"></json-viewer>
        </div>
        <div class="right-column">
          <a-text>输出数据：</a-text>
          <p style="margin: 0" v-if="getJSON(record.Output)==''">
        {{ record.Output }}
      </p>
      <json-viewer v-else :value="getJSON(record.Output)"></json-viewer>
        </div>
      </div>

    </template>

    <template #bodyCell="{ column, record,text }">
      <template v-if="column.dataIndex === 'RPAName'">
        <a-tag :bordered="false" :color="getColor(text)">{{ text }}</a-tag>
      </template>

      <template v-if="column.dataIndex === 'CreatedAt'">
        {{ formatTime(text) }}
      </template>

      <template v-if="column.dataIndex === 'UpdatedAt'">
        {{ formatTime(text) }}
      </template>

      <template v-if="column.dataIndex === 'State'">

        <a-tag v-if="record.State == 'error'" :bordered="false" color="error">执行出错</a-tag>
        <a-tag v-else-if="record.State == 'success'" :bordered="false" color="success">执行成功</a-tag>
        <a-tag v-else-if="record.State == 'running'" :bordered="false" color="processing">正在执行</a-tag>
        <a-tag v-else :bordered="false" color="default">等待执行</a-tag>
      </template>
      <template v-if="column.dataIndex === 'Type'">
        <a :href="record.Type">{{ record.Type }}</a>
      </template>
    </template>

  </a-table>

</template>

<script lang="ts" setup>
import { message, type TableColumnType, type TableProps } from 'ant-design-vue';
import { onMounted, reactive, ref, type Ref } from 'vue';
import { getRPADict, getTaskDB, countTask } from "@/request/api";
import { getColor,formatTime } from "@/func";
import type { taskI, filterI } from "@/interface";
import JsonViewer from "vue-json-viewer";
import "vue-json-viewer/style.css";

// 状态和引用
const searchInput = ref();
const RPANameDict: Ref<filterI[]> = ref([]);
const cardData = ref<{ [key: string]: number }>({});
const data: Ref<taskI[]> = ref([]);

const state = reactive({
  searchText: '',
  searchedColumn: '',
});

// 常量
const cardColor = {
  all: { title: '任务数量', color: '#73767a' },
  waiting: { title: '等待执行', color: '#73767a' },
  running: { title: '正在执行', color: '#337ecc' },
  success: { title: '执行成功', color: '#529b2e' },
  error: { title: '执行错误', color: '#c45656' }
};

// 表格列定义
const columns: TableColumnType<taskI>[] = [
  {
    title: 'ID',
    dataIndex: 'ID',
    width: 50,
    align: "center",
    customFilterDropdown: true,
    onFilter: (value, record) => record.ID.toString().includes(value.toString()),
    onFilterDropdownOpenChange: visible => {
      if (visible) {
        setTimeout(() => {
          searchInput.value.focus();
        }, 100);
      }
    },
  },
  {
    title: '任务类型',
    dataIndex: 'RPAName',
    align: "center",
    filters: RPANameDict.value,
    onFilter: (value: string | number | boolean, record: taskI) => record.RPAName.indexOf(String(value)) === 0,
  },
  {
    title: '任务状态',
    dataIndex: 'State',
    align: "center",
    filters: [
      { text: "执行出错", value: "error" },
      { text: "执行成功", value: "success" },
      { text: "正在执行", value: "running" },
      { text: "等待执行", value: "waiting" }
    ],
    onFilter: (value: string | number | boolean, record: taskI) => record.State.indexOf(String(value)) === 0,
  },
  {
    title: '创建时间',
    dataIndex: 'CreatedAt',
    align: "center",
    defaultSortOrder: 'descend',
    sorter: (a: taskI, b: taskI) => new Date(a.CreatedAt).getTime() - new Date(b.CreatedAt).getTime(),
  },
  {
    title: '更新时间',
    dataIndex: 'UpdatedAt',
    align: "center",
    sorter: (a: taskI, b: taskI) => new Date(a.UpdatedAt).getTime() - new Date(b.UpdatedAt).getTime(),
  },
];

// 函数定义
const handleReset = (clearFilters: (params: { confirm: boolean }) => void) => {
  clearFilters({ confirm: true });
  state.searchText = '';
};

// 搜索
const handleSearch = (selectedKeys: string[], confirm: () => void, dataIndex: string) => {
  confirm();
  state.searchText = selectedKeys[0];
  state.searchedColumn = dataIndex;
};

// 获取统计数据
const fetchData = async () => {
  try {
    const res = await countTask();
    cardData.value = res.data.data;
  } catch (err) {
    message.error("获取任务数据失败");
  }
};

// 获取所有任务记录
const getData = (state: string) => {
  getTaskDB(state).then(res => {
    data.value = res.data.data.map((redata: any) => redata);
  });
};

const getJSON = (jsonValue: any) => {
  if (typeof jsonValue === "string") {
    try {
      const obj = JSON.parse(jsonValue);
      if (typeof obj === "object" && obj) {
        return obj;
      }
    } catch (e) {
      console.log("error：" + jsonValue + "!!!" + e);
      return "";
    }
  }
  return jsonValue;
};

// 初始化
onMounted(() => {
  fetchData()
  getData("all")
  
  getRPADict("").then(res => {
  for (let key in res.data.data) {
    RPANameDict.value.push({
      text: res.data.data[key],
      value: key
    }
    )
  }
});
})

</script>

<style scoped>
.cards {
  display: flex;
  justify-content: space-around;
  flex-wrap: wrap;
  width:calc(100% - 40px);
}


.card {
  flex: 1;
  margin: 0px 10px 0px 10px;
  padding: 10px;
  text-align: center;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.card-title {
  line-height: 18px;
  font-size: 18px;
  margin-bottom: 4px;
}

.card-number {
  font-size: 30px;
  font-weight: bold;
}

.delButton {
  margin-left: 10px;
}

.taskTable {
  height: 70vh;
  /* overflow-y: auto; */
}

.container {
  display: flex; /* 使用 Flexbox 布局 */
}

.left-column {
  flex: 1;
}

.right-column {
  flex: 1; 
  padding-left: 20px; 
}
</style>
