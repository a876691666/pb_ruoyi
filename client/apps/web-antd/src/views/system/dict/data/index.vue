<script setup lang="ts">
import type { VbenFormProps } from '@vben/common-ui';

import type { VxeGridProps } from '#/adapter/vxe-table';
import type { QueryType } from '#/api/common';
import type { DictData } from '#/api/system/dict/dict-data-model';

import { ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { getVxePopupContainer } from '@vben/utils';

import { Modal, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid, vxeCheckboxChecked } from '#/adapter/vxe-table';
import {
  dictDataExport,
  dictDataList,
  dictDataRemove,
} from '#/api/system/dict/dict-data';
import { commonDownloadExcel } from '#/utils/file/download';

import { emitter } from '../mitt';
import { columns, querySchema } from './data';
import dictDataDrawer from './dict-data-drawer.vue';

const dict_type = ref('');

const queryType: QueryType = {
  dict_label: 'LIKE',
  dict_type: 'EQ',
  status: 'EQ',
};

const formOptions: VbenFormProps = {
  commonConfig: {
    labelWidth: 80,
    componentProps: {
      allowClear: true,
    },
  },
  schema: querySchema(),
  wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
};

const gridOptions: VxeGridProps = {
  checkboxConfig: {
    // 高亮
    highlight: true,
    // 翻页时保留选中状态
    reserve: true,
    // 点击行选中
    // trigger: 'row',
  },
  columns,
  height: 'auto',
  keepSource: true,
  proxyConfig: {
    ajax: {
      query: async ({ page, sorts }, params = {}) => {
        if (dict_type.value) {
          params.dict_type = dict_type.value;
        }

        return await dictDataList({
          page,
          sorts,
          params,
          queryType,
        });
      },
    },
  },
  rowConfig: {
    keyField: 'id',
  },
  sortConfig: {
    // 远程排序
    remote: true,
    // 支持多字段排序 默认关闭
    multiple: true,
  },
  id: 'system-dict-data-index',
};

const [BasicTable, tableApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
  gridEvents: {
    sortChange: () => tableApi.query(),
  },
});

const [DictDataDrawer, drawerApi] = useVbenDrawer({
  connectedComponent: dictDataDrawer,
});

function handleAdd() {
  drawerApi.setData({ dict_type: dict_type.value });
  drawerApi.open();
}

async function handleEdit(record: DictData) {
  drawerApi.setData({
    dict_type: dict_type.value,
    id: record.id,
  });
  drawerApi.open();
}

async function handleDelete(row: DictData) {
  await dictDataRemove([row.id]);
  await tableApi.query();
}

function handleMultiDelete() {
  const rows = tableApi.grid.getCheckboxRecords();
  const ids = rows.map((row: DictData) => row.id);
  Modal.confirm({
    title: '提示',
    okType: 'danger',
    content: `确认删除选中的${ids.length}条记录吗？`,
    onOk: async () => {
      await dictDataRemove(ids);
      await tableApi.query();
    },
  });
}

function handleDownloadExcel() {
  commonDownloadExcel(dictDataExport, '字典数据', tableApi.formApi.form.values);
}

emitter.on('rowClick', async (value) => {
  dict_type.value = value;
  await tableApi.query();
});
</script>

<template>
  <div>
    <BasicTable id="dict-data" table-title="字典数据列表">
      <template #toolbar-tools>
        <Space>
          <a-button
            v-access:code="['dict:export']"
            @click="handleDownloadExcel"
          >
            {{ $t('pages.common.export') }}
          </a-button>
          <a-button
            :disabled="!vxeCheckboxChecked(tableApi)"
            danger
            type="primary"
            v-access:code="['dict:remove']"
            @click="handleMultiDelete"
          >
            {{ $t('pages.common.delete') }}
          </a-button>
          <a-button
            :disabled="dict_type === ''"
            type="primary"
            v-access:code="['dict:add']"
            @click="handleAdd"
          >
            {{ $t('pages.common.add') }}
          </a-button>
        </Space>
      </template>
      <template #action="{ row }">
        <Space>
          <ghost-button
            v-access:code="['dict:edit']"
            @click="handleEdit(row)"
          >
            {{ $t('pages.common.edit') }}
          </ghost-button>
          <Popconfirm
            :get-popup-container="
              (node) => getVxePopupContainer(node, 'dict-data')
            "
            placement="left"
            title="确认删除？"
            @confirm="handleDelete(row)"
          >
            <ghost-button
              danger
              v-access:code="['dict:remove']"
              @click.stop=""
            >
              {{ $t('pages.common.delete') }}
            </ghost-button>
          </Popconfirm>
        </Space>
      </template>
    </BasicTable>
    <DictDataDrawer @reload="tableApi.query()" />
  </div>
</template>
