<script setup lang="ts">
import type { VbenFormProps } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import type { VxeGridProps } from '#/adapter/vxe-table';
import type { QueryType } from '#/api/common';
import type { GenInfo } from '#/api/tool/gen/model';

import { useRouter } from 'vue-router';

import { Page, useVbenModal } from '@vben/common-ui';
import { getVxePopupContainer } from '@vben/utils';

import { message, Modal, Popconfirm, Space } from 'ant-design-vue';
import dayjs from 'dayjs';

import { useVbenVxeGrid, vxeCheckboxChecked } from '#/adapter/vxe-table';
import { generatedList, genRemove, syncDb } from '#/api/tool/gen';
import { downloadByData } from '#/utils/file/download';

import codePreviewModal from './code-preview-modal.vue';
import { columns, querySchema } from './data';
import { createZipBlob, generateCodeFromGenInfo } from './index';
import tableImportModal from './table-import-modal.vue';

const queryType: QueryType = {
  table_name: 'LIKE',
  table_comment: 'LIKE',
  create_time: 'BETWEEN',
};

const formOptions: VbenFormProps = {
  schema: querySchema(),
  commonConfig: {
    labelWidth: 80,
    componentProps: {
      allowClear: true,
    },
  },
  wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
};

const gridOptions: VxeGridProps = {
  checkboxConfig: {
    // 高亮
    highlight: true,
    // 翻页时保留选中状态
    reserve: true,
    // 点击行选中
    trigger: 'row',
  },
  columns,
  height: 'auto',
  keepSource: true,
  proxyConfig: {
    ajax: {
      query: async ({ page, sorts }, params = {}) => {
        return await generatedList({
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
  id: 'tool-gen-index',
};

const [BasicTable, tableApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
  gridEvents: {
    sortChange: () => tableApi.query(),
  },
});

const [CodePreviewModal, previewModalApi] = useVbenModal({
  connectedComponent: codePreviewModal,
});

function handlePreview(record: Recordable<any>) {
  previewModalApi.setData(record);
  previewModalApi.open();
}

const router = useRouter();
function handleEdit(record: Recordable<any>) {
  router.push(`/tool/gen-edit/index/${record.id}`);
}

async function handleSync(record: GenInfo) {
  await syncDb(record);
  await tableApi.query();
}

/**
 * 批量生成代码
 */
async function handleBatchGen() {
  const rows = tableApi.grid.getCheckboxRecords();
  const ids = rows.map((row: any) => row.id);
  if (ids.length === 0) {
    message.info('请选择需要生成代码的表');
    return;
  }
  const hideLoading = message.loading('下载中...');
  try {
    // 获取所有表的信息并生成代码
    const allFiles: { [key: string]: string } = {};

    for (const row of rows) {
      const files = generateCodeFromGenInfo(row);

      // 将文件添加到总文件列表中,使用表名作为子目录
      for (const [filePath, content] of Object.entries(files)) {
        allFiles[filePath] = content;
      }
    }

    // 打包成 zip
    const blob = await createZipBlob(allFiles);
    const timestamp = Date.now();
    downloadByData(blob, `批量代码生成_${timestamp}.zip`);
  } finally {
    hideLoading();
  }
}

async function handleDownload(record: GenInfo) {
  const hideLoading = message.loading('加载中...');
  try {
    // 生成代码文件
    const files = await generateCodeFromGenInfo(record);

    // 打包成 zip
    const blob = await createZipBlob(files);
    const filename = `代码生成_${record.name}_${dayjs().valueOf()}.zip`;
    downloadByData(blob, filename);
  } catch (error) {
    console.error(error);
  } finally {
    hideLoading();
  }
}

/**
 * 删除
 * @param record
 */
async function handleDelete(record: Recordable<any>) {
  await genRemove([record.id]);
  await tableApi.query();
}

function handleMultiDelete() {
  const rows = tableApi.grid.getCheckboxRecords();
  const ids = rows.map((row: any) => row.tableId);
  Modal.confirm({
    title: '提示',
    okType: 'danger',
    content: `确认删除选中的${ids.length}条记录吗？`,
    onOk: async () => {
      await genRemove(ids);
      await tableApi.query();
    },
  });
}

const [TableImportModal, tableImportModalApi] = useVbenModal({
  connectedComponent: tableImportModal,
});

function handleImport() {
  tableImportModalApi.open();
}
</script>

<template>
  <Page :auto-content-height="true">
    <BasicTable table-title="代码生成列表">
      <template #toolbar-tools>
        <a
          class="text-primary mr-2"
          href="https://dapdap.top/other/template.html"
          target="_blank"
          >👉关于代码生成模板
        </a>
        <Space>
          <a-button
            :disabled="!vxeCheckboxChecked(tableApi)"
            danger
            type="primary"
            v-access:code="['tool:gen:remove']"
            @click="handleMultiDelete"
          >
            {{ $t('pages.common.delete') }}
          </a-button>
          <a-button
            :disabled="!vxeCheckboxChecked(tableApi)"
            v-access:code="['tool:gen:code']"
            @click="handleBatchGen"
          >
            {{ $t('pages.common.generate') }}
          </a-button>
          <a-button
            type="primary"
            v-access:code="['tool:gen:import']"
            @click="handleImport"
          >
            {{ $t('pages.common.import') }}
          </a-button>
        </Space>
      </template>
      <template #action="{ row }">
        <a-button
          size="small"
          type="link"
          v-access:code="['tool:gen:preview']"
          @click.stop="handlePreview(row)"
        >
          {{ $t('pages.common.preview') }}
        </a-button>
        <a-button
          size="small"
          type="link"
          v-access:code="['tool:gen:edit']"
          @click.stop="handleEdit(row)"
        >
          {{ $t('pages.common.edit') }}
        </a-button>
        <Popconfirm
          :get-popup-container="getVxePopupContainer"
          :title="`确认同步[${row.name}]?`"
          placement="left"
          @confirm="handleSync(row)"
        >
          <a-button
            size="small"
            type="link"
            v-access:code="['tool:gen:edit']"
            @click.stop=""
          >
            {{ $t('pages.common.sync') }}
          </a-button>
        </Popconfirm>
        <a-button
          size="small"
          type="link"
          v-access:code="['tool:gen:code']"
          @click.stop="handleDownload(row)"
        >
          生成代码
        </a-button>
        <Popconfirm
          :get-popup-container="getVxePopupContainer"
          :title="`确认删除[${row.name}]?`"
          placement="left"
          @confirm="handleDelete(row)"
        >
          <a-button
            danger
            size="small"
            type="link"
            v-access:code="['tool:gen:remove']"
            @click.stop=""
          >
            {{ $t('pages.common.delete') }}
          </a-button>
        </Popconfirm>
      </template>
    </BasicTable>
    <CodePreviewModal />
    <TableImportModal @reload="tableApi.query()" />
  </Page>
</template>
