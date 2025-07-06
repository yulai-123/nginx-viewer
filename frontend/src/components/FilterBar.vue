<template>
  <div class="filter-bar">
    <el-form :inline="true" :model="formData" class="filter-form">
      <el-form-item label="IP地址">
        <el-input v-model="formData.ip" placeholder="客户端IP" clearable />
      </el-form-item>

      <el-form-item label="状态码">
        <el-input v-model="formData.status" placeholder="状态码" clearable style="width: 150px;" />
      </el-form-item>

      <el-form-item label="路径">
        <el-input v-model="formData.path" placeholder="请求路径" clearable />
      </el-form-item>

      <el-form-item label="查询长度">
        <el-input-number v-model="formData.queryLimit" :min="100" :max="1000000" :step="100" placeholder="最大返回条数" style="width: 150px;" />
      </el-form-item>

      <el-form-item label="时间范围">
        <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DDTHH:mm:ss.SSSZ"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleSearch">查询</el-button>
        <el-button @click="handleReset">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

// 定义事件
const emit = defineEmits(['filter-changed']);

// 表单数据
const formData = ref({
  ip: '',
  status: '',
  path: '',
  from: '',
  to: '',
  queryLimit: 1000 // 默认查询长度
});

// 日期范围选择器
const dateRange = ref<[string, string] | null>(null);

// 监听日期范围变化
watch(dateRange, (newVal) => {
  if (newVal) {
    formData.value.from = newVal[0];
    formData.value.to = newVal[1];
  } else {
    formData.value.from = '';
    formData.value.to = '';
  }
});

// 处理搜索按钮点击
const handleSearch = () => {
  const params: Record<string, any> = {};

  if (formData.value.ip) params.ip = formData.value.ip;
  if (formData.value.status) params.status = Number(formData.value.status);
  if (formData.value.path) params.path = formData.value.path;
  if (formData.value.from) params.from = formData.value.from;
  if (formData.value.to) params.to = formData.value.to;
  if (formData.value.queryLimit) params.limit = formData.value.queryLimit;

  emit('filter-changed', params);
};

// 处理重置按钮点击
const handleReset = () => {
  formData.value = {
    ip: '',
    status: '',
    path: '',
    from: '',
    to: '',
    queryLimit: 1000
  };
  dateRange.value = null;
  emit('filter-changed', {});
};
</script>

<style scoped>
.filter-bar {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
}
</style>