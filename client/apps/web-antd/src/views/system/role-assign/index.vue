<script setup lang="ts">
import type { VbenFormProps } from '@vben/common-ui';

import type { VxeGridProps } from '#/adapter/vxe-table';
import type { QueryType } from '#/api/common';
import type { User } from '#/api/system/user/model';

import { useRoute } from 'vue-router';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { getVxePopupContainer } from '@vben/utils';

import { Modal, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid, vxeCheckboxChecked } from '#/adapter/vxe-table';
import {
  roleAllocatedList,
  roleAuthCancel,
  roleAuthCancelAll,
} from '#/api/system/role';

import { columns, querySchema } from './data';
import roleAssignDrawer from './role-assign-drawer.vue';

const route = useRoute();
const role_id = route.params.roleId as string;

const queryType: QueryType = {
  user_name: 'LIKE',
  phonenumber: 'LIKE',
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
        params.role_id = role_id;
        return await roleAllocatedList({
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
  id: 'system-role-assign-index',
};

const [BasicTable, tableApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
  gridEvents: {
    sortChange: () => tableApi.query(),
  },
});

const [RoleAssignDrawer, drawerApi] = useVbenDrawer({
  connectedComponent: roleAssignDrawer,
});

function handleAdd() {
  drawerApi.setData({});
  drawerApi.open();
}

/**
 * 取消授权 一条记录
 */
async function handleAuthCancel(record: User) {
  await roleAuthCancel(record.id);
  await tableApi.query();
}

/**
 * 批量取消授权
 */
function handleMultipleAuthCancel() {
  const rows = tableApi.grid.getCheckboxRecords();
  const ids = rows.map((row: User) => row.id);
  Modal.confirm({
    title: '提示',
    okType: 'danger',
    content: `确认取消选中的${ids.length}条授权记录吗？`,
    onOk: async () => {
      await roleAuthCancelAll(ids);
      await tableApi.query();
      tableApi.grid.clearCheckboxRow();
    },
  });
}
</script>

<template>
  <Page :auto-content-height="true">
    <BasicTable table-title="已分配的用户列表">
      <template #toolbar-tools>
        <Space>
          <a-button
            :disabled="!vxeCheckboxChecked(tableApi)"
            danger
            type="primary"
            v-access:code="['system:role:remove']"
            @click="handleMultipleAuthCancel"
          >
            取消授权
          </a-button>
          <a-button
            type="primary"
            v-access:code="['system:role:add']"
            @click="handleAdd"
          >
            {{ $t('pages.common.add') }}
          </a-button>
        </Space>
      </template>
      <template #action="{ row }">
        <Popconfirm
          :get-popup-container="getVxePopupContainer"
          :title="`是否取消授权用户[${row?.expand?.user?.user_name} - ${row?.expand?.user?.nick_name}]?`"
          placement="left"
          @confirm="handleAuthCancel(row)"
        >
          <ghost-button
            danger
            v-access:code="['system:role:remove']"
            @click.stop=""
          >
            取消授权
          </ghost-button>
        </Popconfirm>
      </template>
    </BasicTable>
    <RoleAssignDrawer @reload="tableApi.query()" />
  </Page>
</template>
