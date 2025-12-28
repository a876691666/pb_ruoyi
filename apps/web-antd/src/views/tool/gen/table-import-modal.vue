<script setup lang="ts">
import type { VbenFormProps } from '@vben/common-ui';

import type { VxeGridProps } from '#/adapter/vxe-table';

import { useVbenModal } from '@vben/common-ui';

import dayjs from 'dayjs';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { importTable, readyToGenList } from '#/api/tool/gen';

const emit = defineEmits<{ reload: [] }>();

// 用于存储完整数据的缓存
let cachedGeneratedList: any[] = [];

/**
 * 前端分页函数 - 从缓存中获取分页数据
 * @param currentPage 当前页码
 * @param pageSize 每页数量
 * @param forceRefresh 是否强制刷新缓存
 * @returns 分页数据和总数
 */
async function getPagedDataFromCache(
  currentPage: number = 1,
  pageSize: number = 10,
  name: string = '',
  forceRefresh: boolean = false,
) {
  // 如果缓存为空或强制刷新，则从API加载数据
  if (cachedGeneratedList.length === 0 || forceRefresh) {
    const result = await readyToGenList();
    cachedGeneratedList = result.filter((item) => !item.system);
  }

  // 计算分页
  const startIndex = (currentPage - 1) * pageSize;
  const endIndex = startIndex + pageSize;
  const pageData = cachedGeneratedList
    .filter((item) => item.name.includes(name))
    .slice(startIndex, endIndex);

  return {
    items: pageData,
    totalItems: cachedGeneratedList.length,
  };
}

const formOptions: VbenFormProps = {
  schema: [
    {
      label: '表名称',
      fieldName: 'name',
      component: 'Input',
    },
  ],
  commonConfig: {
    labelWidth: 60,
  },
  showCollapseButton: false,
  wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
};

const gridOptions: VxeGridProps = {
  checkboxConfig: {
    highlight: true,
    reserve: true,
    trigger: 'row',
  },
  columns: [
    {
      type: 'checkbox',
      width: 60,
    },
    {
      title: '表名称',
      field: 'name',
      align: 'left',
    },
    {
      field: 'fields_length',
      title: '字段数量',
      formatter: ({ row }) => {
        return row.fields ? row.fields.length : 0;
      },
    },
    {
      title: '创建时间',
      field: 'created',
      formatter: ({ cellValue }) => {
        return cellValue ? dayjs(cellValue).format('YYYY-MM-DD HH:mm:ss') : '';
      },
    },
    {
      title: '更新时间',
      field: 'updated',
      formatter: ({ cellValue }) => {
        return cellValue ? dayjs(cellValue).format('YYYY-MM-DD HH:mm:ss') : '';
      },
    },
  ],
  keepSource: true,
  size: 'small',
  minHeight: 400,
  pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues = {}) => {
        return await getPagedDataFromCache(
          page.currentPage,
          page.pageSize,
          formValues.name,
        );
      },
    },
  },
  rowConfig: {
    keyField: 'tableName',
  },
  toolbarConfig: {
    enabled: false,
  },
  id: 'import-table-modal',
  cellClassName: 'cursor-pointer',
};

const [BasicTable, tableApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [BasicModal, modalApi] = useVbenModal({
  onOpenChange: async (isOpen) => {
    if (!isOpen) {
      tableApi.grid.clearCheckboxRow();
      return null;
    }
  },
  onConfirm: handleSubmit,
});

async function handleSubmit() {
  try {
    const records = tableApi.grid.getCheckboxRecords();
    if (records.length === 0) {
      modalApi.close();
      return;
    }
    modalApi.modalLoading(true);
    await importTable(records);
    emit('reload');
    modalApi.close();
  } catch (error) {
    console.warn(error);
  } finally {
    modalApi.modalLoading(false);
  }
}
</script>

<template>
  <BasicModal class="w-[800px]" title="导入表">
    <BasicTable />
  </BasicModal>
</template>
