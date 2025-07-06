<template>
  <div class="logs-container">
    <h1>Nginx 日志查看器</h1>

    <!-- 过滤组件 -->
    <FilterBar @filter-changed="applyFilter" />

    <!-- 分页 -->
    <div class="pagination-container">
      <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[50, 100, 300, 500, 1000, 3000]"
          layout="total, sizes, prev, pager, next"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
      />

      <!--      <el-button type="primary" @click="exportCSV">导出CSV</el-button>-->
    </div>

    <!-- 日志表格 -->
    <el-table
        v-loading="loading"
        :data="logs"
        style="width: 100%; margin-top: 20px;"
        border
        stripe
        :default-sort="{ prop: 'Time', order: 'descending' }"
    >
      <el-table-column
          label="序号"
          width="70"
          align="center">
        <template #default="scope">
          {{ (currentPage - 1) * pageSize + scope.$index + 1 }}
        </template>
      </el-table-column>
      <el-table-column prop="Time" label="时间" width="180" :formatter="formatTime" sortable />
      <el-table-column
          prop="ClientIP"
          label="客户端IP"
          width="140"
          sortable
          column-key="ClientIP"
          :filters="getUniqueValues('ClientIP')"
          :filter-method="filterHandler"
      />
      <el-table-column
          prop="Host"
          label="主机名"
          width="130"
          show-overflow-tooltip
          sortable
          column-key="Host"
          :filters="getUniqueValues('Host')"
          :filter-method="filterHandler"
      />
      <el-table-column
          prop="ServerPort"
          label="服务端口"
          width="120"
          sortable
          column-key="ServerPort"
          :filters="getUniqueValues('ServerPort')"
          :filter-method="filterHandler"
      />
      <el-table-column
          prop="Method"
          label="方法"
          width="100"
          show-overflow-tooltip
          sortable
          column-key="Method"
          :filters="getUniqueValues('Method')"
          :filter-method="filterHandler"
      />
      <el-table-column
          prop="Path"
          label="路径"
          min-width="200"
          show-overflow-tooltip
          sortable
          column-key="Path"
          :filters="getUniqueValues('Path')"
          :filter-method="filterHandler"
      />
      <el-table-column
          prop="Status"
          label="状态码"
          width="105"
          sortable
          column-key="Status"
          :filters="getStatusFilters()"
          :filter-method="filterStatusHandler"
      >
        <template #default="scope">
          <span :class="getStatusClass(scope.row.Status)">{{ scope.row.Status }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="BodyBytes" label="响应大小" width="105" :formatter="formatBytes" sortable />
      <el-table-column prop="ReqTime" label="请求时间(s)" width="125" sortable />
      <el-table-column
          prop="Referer"
          label="来源"
          min-width="150"
          show-overflow-tooltip
          sortable
          column-key="Referer"
          :filters="getUniqueValues('Referer')"
          :filter-method="filterHandler"
      />
    </el-table>


  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import FilterBar from '../components/FilterBar.vue'

interface LogRow {
  Time: string;
  ClientIP: string;
  ClientPort: string;
  XFF: string;
  Host: string;
  ServerPort: string;
  Method: string;
  Path: string;
  HTTPVer: string;
  Status: number;
  BodyBytes: number;
  ReqBytes: number;
  ReqTime: number;
  UpConnTime: number | null;
  UpRespTime: number | null;
  UpStatus: number | null;
  UpAddr: string | null;
  Referer: string;
  UA: string;
  TLSProto: string;
  TLSCipher: string;
  ReqID: string;
}

interface FilterParams {
  ip?: string;
  status?: number;
  path?: string;
  from?: string;
  to?: string;
  limit: number; // 新增查询长度参数
}

const logs = ref<LogRow[]>([]);
const loading = ref(false);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(50);
const filterParams = ref<FilterParams>({
  limit: 1000,
});

const allLogs = ref<LogRow[]>([]);

// 获取日志数据
const fetchLogs = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/api/logs', { params: filterParams.value });
    allLogs.value = response.data.logs;
    total.value = response.data.total;
    // 应用前端分页
    applyPagination();
  } catch (error) {
    ElMessage.error('获取日志数据失败');
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// 前端分页处理函数
const applyPagination = () => {
  const startIndex = (currentPage.value - 1) * pageSize.value;
  const endIndex = startIndex + pageSize.value;
  logs.value = allLogs.value.slice(startIndex, endIndex);
};

// 应用过滤条件
const applyFilter = (params: Partial<FilterParams>) => {
  // 如果是空对象，表示重置，只保留分页参数
  if (Object.keys(params).length === 0) {
    filterParams.value = {
      limit: 1000, // 默认查询长度
    };
  } else {
    filterParams.value = {
      ...filterParams.value,
      ...params,
    };
  }
  currentPage.value = 1;
  fetchLogs();
};

// 处理页码变化
const handleCurrentChange = (val: number) => {
  currentPage.value = val;
  applyPagination(); // 重新应用分页
};

// 处理每页条数变化
const handleSizeChange = (val: number) => {
  pageSize.value = val;
  // 如果当前页超出范围，重置为第一页
  if (currentPage.value > Math.ceil(total.value / pageSize.value)) {
    currentPage.value = 1;
  }
  applyPagination(); // 重新应用分页
};

// 获取表格列的唯一值作为筛选选项
const getUniqueValues = (prop: keyof LogRow) => {
  const uniqueSet = new Set<string>();
  allLogs.value.forEach(log => {
    const value = log[prop];
    if (value && typeof value === 'string' && value !== '-') {
      uniqueSet.add(value);
    }
  });

  return Array.from(uniqueSet)
      .slice(0, 20) // 限制数量以避免过多
      .map(value => ({ text: value, value }));
};

// 获取状态码筛选选项
const getStatusFilters = () => {
  const statusSet = new Set<number>();
  allLogs.value.forEach(log => {
    if (log.Status >= 100 && log.Status < 600) {
      statusSet.add(log.Status);
    }
  });

  return Array.from(statusSet)
      .sort((a, b) => a - b) // 按照状态码排序
      .map(status => ({ text: String(status), value: status }));
};

// 通用筛选处理方法
const filterHandler = (value: any, row: LogRow, column: { property: keyof LogRow }) => {
  const cellValue = row[column.property];
  if (cellValue === null || cellValue === undefined) return false;
  return String(cellValue) === String(value);
};

// 状态码筛选处理方法
const filterStatusHandler = (value: number, row: LogRow) => {
  if (value >= 100 && value < 600) {
    // 如果是像200、300、400、500这样的值，按区间匹配
    if (value % 100 === 0) {
      const statusPrefix = Math.floor(value / 100);
      const rowPrefix = Math.floor(row.Status / 100);
      return statusPrefix === rowPrefix;
    }
    // 如果是具体状态码如404，直接匹配
    return row.Status === value;
  }
  return false;
};

// 格式化时间
const formatTime = (row: LogRow) => {
  const date = new Date(row.Time);
  return date.toLocaleString();
};

// 格式化字节数
const formatBytes = (row: any, column: any) => {
  const bytes = row[column.property];
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
};

// 获取状态码样式
const getStatusClass = (status: number) => {
  if (status >= 500) return 'status-error';
  if (status >= 400) return 'status-warning';
  return 'status-success';
};

onMounted(() => {
  fetchLogs();
});
</script>

<style scoped>
.logs-container {
  padding: 20px 0;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-error {
  color: #F56C6C;
  font-weight: bold;
}

.status-warning {
  color: #E6A23C;
  font-weight: bold;
}

.status-success {
  color: #67C23A;
}
</style>