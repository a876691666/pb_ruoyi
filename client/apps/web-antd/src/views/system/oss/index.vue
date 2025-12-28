<script setup lang="ts">
import type { VbenFormProps } from '@vben/common-ui';

import type { VxeGridProps } from '#/adapter/vxe-table';
import type { QueryType } from '#/api/common';
import type { OssFile } from '#/api/system/oss/model';

import { onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import { useAppConfig } from '@vben/hooks';
import { $t } from '@vben/locales';
import { stringify } from '@vben/request';
import { useAccessStore } from '@vben/stores';
import { getVxePopupContainer } from '@vben/utils';

import {
  Image,
  Modal,
  Popconfirm,
  Space,
  Spin,
  Switch,
  Tooltip,
} from 'ant-design-vue';

import { useVbenVxeGrid, vxeCheckboxChecked } from '#/adapter/vxe-table';
import { configInfoByKey } from '#/api/system/config';
import { checkLoginBeforeDownload, ossList, ossRemove } from '#/api/system/oss';
import { downloadByUrl } from '#/utils/file/download';

import { supportImageList } from './constant';
import { columns, querySchema } from './data';
import fallbackImageBase64 from './fallback-image.txt?raw';
import fileUploadModal from './file-upload-modal.vue';
import imageUploadModal from './image-upload-modal.vue';

const queryType: QueryType = {
  file_name: 'LIKE',
  original_name: 'LIKE',
  file_suffix: 'LIKE',
  url: 'LIKE',
  service: 'EQ',
  create_time: 'BETWEEN',
};

const formOptions: VbenFormProps = {
  commonConfig: {
    labelWidth: 80,
    componentProps: {
      allowClear: true,
    },
  },
  schema: querySchema(),
  wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
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
        return await ossList({
          page,
          sorts,
          params,
          queryType,
        });
      },
    },
  },
  headerCellConfig: {
    height: 44,
  },
  cellConfig: {
    height: 65,
  },
  rowConfig: {
    keyField: 'id',
  },
  sortConfig: {
    // 远程排序
    remote: true,
    // 支持多字段排序 默认关闭
    multiple: false,
  },
  id: 'system-oss-index',
};

const [BasicTable, tableApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
  gridEvents: {
    // 排序 重新请求接口
    sortChange: () => tableApi.query(),
  },
});

const { apiURL, clientId } = useAppConfig(
  import.meta.env,
  import.meta.env.PROD,
);
const accessStore = useAccessStore();

/**
 * 浏览器直接接管下载 相较于axios请求 不会阻塞
 * @param row oss信息
 */
async function handleDownload(row: OssFile) {
  await checkLoginBeforeDownload();

  const params = {
    clientid: clientId,
    Authorization: `Bearer ${accessStore.accessToken}`,
  };

  const downloadLink = `${apiURL}/resource/oss/download/${row.id}?${stringify(params)}`;
  // file_name不设置也行 默认会取header里的名称
  downloadByUrl({ fileName: row.file_name, url: downloadLink });
}

async function handleDelete(row: OssFile) {
  await ossRemove([row.id]);
  await tableApi.query();
}

function handleMultiDelete() {
  const rows = tableApi.grid.getCheckboxRecords();
  const ids = rows.map((row: OssFile) => row.id);
  Modal.confirm({
    title: '提示',
    okType: 'danger',
    content: `确认删除选中的${ids.length}条记录吗？`,
    onOk: async () => {
      await ossRemove(ids);
      await tableApi.query();
    },
  });
}

const preview = ref(false);
onMounted(async () => {
  const previewStr = await configInfoByKey('sys.oss.previewListResource');
  preview.value = previewStr === 'true';
});

/**
 * 根据拓展名判断是否是图片
 * @param ext 拓展名
 */
function isImageFile(ext: string) {
  return supportImageList.some((item) =>
    ext.toLocaleLowerCase().includes(item),
  );
}

/**
 * 判断是否是pdf文件
 * @param ext 扩展名
 */
function isPdfFile(ext: string) {
  return ext.toLocaleLowerCase().includes('pdf');
}

/**
 * pdf预览 使用浏览器接管
 * @param url 文件地址
 */
function pdfPreview(url: string) {
  window.open(url);
}

const [ImageUploadModal, imageUploadApi] = useVbenModal({
  connectedComponent: imageUploadModal,
});

const [FileUploadModal, fileUploadApi] = useVbenModal({
  connectedComponent: fileUploadModal,
});
</script>

<template>
  <Page :auto-content-height="true">
    <BasicTable table-title="文件列表">
      <template #toolbar-tools>
        <Space>
          <Tooltip title="预览图片">
            <Switch v-model:checked="preview" />
          </Tooltip>
          <a-button
            :disabled="!vxeCheckboxChecked(tableApi)"
            danger
            type="primary"
            v-access:code="['system:oss:remove']"
            @click="handleMultiDelete"
          >
            {{ $t('pages.common.delete') }}
          </a-button>
          <a-button
            v-access:code="['system:oss:upload']"
            @click="fileUploadApi.open"
          >
            文件上传
          </a-button>
          <a-button
            v-access:code="['system:oss:upload']"
            @click="imageUploadApi.open"
          >
            图片上传
          </a-button>
        </Space>
      </template>
      <template #url="{ row }">
        <!-- placeholder为图片未加载时显示的占位图 -->
        <!-- fallback为图片加载失败时显示 -->
        <!-- 需要设置key属性 否则切换翻页会有延迟 -->
        <Image
          :key="row.id"
          v-if="preview && isImageFile(row.url)"
          :src="row.url"
          height="50px"
          :fallback="fallbackImageBase64"
        >
          <template #placeholder>
            <div class="flex size-full items-center justify-center">
              <Spin />
            </div>
          </template>
        </Image>
        <!-- pdf预览 使用浏览器开新窗口 -->
        <span
          v-else-if="preview && isPdfFile(row.url)"
          class="icon-[vscode-icons--file-type-pdf2] size-10 cursor-pointer"
          @click.stop="pdfPreview(row.url)"
        ></span>
        <span v-else>{{ row.url }}</span>
      </template>
      <template #action="{ row }">
        <Space>
          <ghost-button
            v-access:code="['system:oss:download']"
            @click="handleDownload(row)"
          >
            {{ $t('pages.common.download') }}
          </ghost-button>
          <Popconfirm
            :get-popup-container="getVxePopupContainer"
            placement="left"
            title="确认删除？"
            @confirm="handleDelete(row)"
          >
            <ghost-button
              danger
              v-access:code="['system:oss:remove']"
              @click.stop=""
            >
              {{ $t('pages.common.delete') }}
            </ghost-button>
          </Popconfirm>
        </Space>
      </template>
    </BasicTable>
    <ImageUploadModal @reload="tableApi.query" />
    <FileUploadModal @reload="tableApi.query" />
  </Page>
</template>
