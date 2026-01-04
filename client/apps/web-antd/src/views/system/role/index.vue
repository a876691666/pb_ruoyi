<script setup lang="ts">
import type { VbenFormProps } from '@vben/common-ui';

import type { VxeGridProps } from '#/adapter/vxe-table';
import type { QueryType } from '#/api/common';
import type { Role } from '#/api/system/role/model';

import { useRouter } from 'vue-router';

import { useAccess } from '@vben/access';
import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import { getVxePopupContainer } from '@vben/utils';

import { Modal, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid, vxeCheckboxChecked } from '#/adapter/vxe-table';
import {
  roleChangeStatus,
  roleExport,
  roleList,
  roleRemove,
} from '#/api/system/role';
import { TableSwitch } from '#/components/table';
import { commonDownloadExcel } from '#/utils/file/download';

import { columns, querySchema } from './data';
import roleAuthModal from './role-auth-modal.vue';
import roleDrawer from './role-drawer.vue';

const queryType: QueryType = {
  role_name: 'LIKE',
  role_key: 'LIKE',
  status: 'EQ',
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
    checkMethod: ({ row }) => row.id !== '1',
  },
  columns,
  height: 'auto',
  keepSource: true,
  proxyConfig: {
    ajax: {
      query: async ({ page, sorts }, params = {}) => {
        return await roleList({
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
  id: 'system-role-index',
};

const [BasicTable, tableApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
  gridEvents: {
    sortChange: () => tableApi.query(),
  },
});
const [RoleDrawer, drawerApi] = useVbenDrawer({
  connectedComponent: roleDrawer,
});

function handleAdd() {
  drawerApi.setData({});
  drawerApi.open();
}

async function handleEdit(record: Role) {
  drawerApi.setData({ id: record.id });
  drawerApi.open();
}

async function handleDelete(row: Role) {
  await roleRemove([row.id]);
  await tableApi.query();
}

function handleMultiDelete() {
  const rows = tableApi.grid.getCheckboxRecords();
  const ids = rows.map((row: Role) => row.id);
  Modal.confirm({
    title: '提示',
    okType: 'danger',
    content: `确认删除选中的${ids.length}条记录吗？`,
    onOk: async () => {
      await roleRemove(ids);
      await tableApi.query();
    },
  });
}

function handleDownloadExcel() {
  commonDownloadExcel(roleExport, '角色数据', tableApi.formApi.form.values, {
    fieldMappingTime: formOptions.fieldMappingTime,
  });
}

const { hasAccessByCodes } = useAccess();

const [RoleAuthModal, authModalApi] = useVbenModal({
  connectedComponent: roleAuthModal,
});

function handleAuthEdit(record: Role) {
  authModalApi.setData({ id: record.id });
  authModalApi.open();
}

const router = useRouter();
function handleAssignRole(record: Role) {
  router.push(`/system/role-auth/user/${record.id}`);
}
</script>

<template>
  <Page :auto-content-height="true">
    <BasicTable table-title="角色列表">
      <template #toolbar-tools>
        <Space>
          <a-button
            v-access:code="['role:export']"
            @click="handleDownloadExcel"
          >
            {{ $t('pages.common.export') }}
          </a-button>
          <a-button
            :disabled="!vxeCheckboxChecked(tableApi)"
            danger
            type="primary"
            v-access:code="['role:remove']"
            @click="handleMultiDelete"
          >
            {{ $t('pages.common.delete') }}
          </a-button>
          <a-button
            type="primary"
            v-access:code="['role:add']"
            @click="handleAdd"
          >
            {{ $t('pages.common.add') }}
          </a-button>
        </Space>
      </template>
      <template #status="{ row }">
        <TableSwitch
          v-model:value="row.status"
          :api="() => roleChangeStatus(row)"
          :disabled="
            row.id === '1' ||
            row.role_key === 'admin' ||
            !hasAccessByCodes(['role:edit'])
          "
          @reload="tableApi.query()"
        />
      </template>
      <template #action="{ row }">
        <!-- 租户管理员不可修改admin角色 防止误操作 -->
        <!-- 超级管理员可通过租户切换来操作租户管理员角色 -->
        <template v-if="!(row.id === '1')">
          <Space>
            <ghost-button
              v-access:code="['role:edit']"
              @click.stop="handleEdit(row)"
            >
              {{ $t('pages.common.edit') }}
            </ghost-button>
            <ghost-button
              v-access:code="['role:edit']"
              @click.stop="handleAuthEdit(row)"
            >
              权限
            </ghost-button>
            <ghost-button
              v-access:code="['role:edit']"
              @click.stop="handleAssignRole(row)"
            >
              分配
            </ghost-button>
            <Popconfirm
              :get-popup-container="getVxePopupContainer"
              placement="left"
              title="确认删除？"
              @confirm="handleDelete(row)"
            >
              <ghost-button
                danger
                v-access:code="['role:remove']"
                @click.stop=""
              >
                {{ $t('pages.common.delete') }}
              </ghost-button>
            </Popconfirm>
          </Space>
        </template>
      </template>
    </BasicTable>
    <RoleDrawer @reload="tableApi.query()" />
    <RoleAuthModal @reload="tableApi.query()" />
  </Page>
</template>
