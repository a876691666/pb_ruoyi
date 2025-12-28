<script setup lang="ts">
import type { VbenFormProps } from '@vben/common-ui';

import type { VxeGridProps } from '#/adapter/vxe-table';
import type { QueryType } from '#/api/common';
import type { Menu } from '#/api/system/menu/model';

import { computed } from 'vue';

import { useAccess } from '@vben/access';
import { Fallback, Page, useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { eachTree, getVxePopupContainer, listToTree } from '@vben/utils';

import { Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { menuList, menuRemove } from '#/api/system/menu';

import { columns, querySchema } from './data';
import menuDrawer from './menu-drawer.vue';

/**
 * 不要问为什么有两个根节点 v-if会控制只会渲染一个
 */

const queryType: QueryType = {
  menu_name: 'LIKE',
  status: 'EQ',
  visible: 'EQ',
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
  columns,
  height: 'auto',
  keepSource: true,
  pagerConfig: {
    enabled: false,
  },
  proxyConfig: {
    ajax: {
      query: async ({ sorts }, params = {}) => {
        const resp = await menuList({
          sorts,
          params,
          queryType,
        });
        // 手动转为树结构
        const treeData = listToTree(resp, { id: 'id', pid: 'parent_id' });
        // 添加hasChildren字段
        eachTree(treeData, (item) => {
          item.hasChildren = !!(item.children && item.children.length > 0);
        });
        console.log(treeData);

        return { rows: treeData };
      },
    },
  },
  rowConfig: {
    keyField: 'id',
    // 高亮点击行
    isCurrent: true,
  },
  sortConfig: {
    // 远程排序
    remote: true,
    // 支持多字段排序 默认关闭
    multiple: true,
  },
  /**
   * 开启虚拟滚动
   * 数据量小可以选择关闭
   * 如果遇到样式问题(空白、错位 滚动等)可以选择关闭虚拟滚动
   *
   * 由于已经重构为懒加载 不需要虚拟滚动(如果你觉得卡顿 依旧可以选择开启)
   */
  // scrollY: {
  //   enabled: true,
  //   gt: 0,
  // },
  treeConfig: {
    parentField: 'parent_id',
    rowField: 'id',
    // 使用懒加载需要自行构造hasChild字段 不需要自动转换为树结构
    transform: false,
    // 刷新接口后 记录展开行的情况
    reserve: true,
    // 是否存在子节点的字段
    hasChildField: 'hasChildren',
    // 开启展开 懒加载
    lazy: true,
    // 懒加载方法 直接返回children
    loadMethod: ({ row }) => row.children ?? [],
  },
  id: 'system-menu-index',
};

const [BasicTable, tableApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
  gridEvents: {
    sortChange: () => tableApi.query(),
    cellDblclick: (e) => {
      const { row = {} } = e;
      if (!row?.children) {
        return;
      }
      const isExpanded = row?.expand;
      tableApi.grid.setTreeExpand(row, !isExpanded);
      row.expand = !isExpanded;
    },
    // 需要监听使用箭头展开的情况 否则展开/折叠的数据不一致
    toggleTreeExpand: (e) => {
      const { row = {}, expanded } = e;
      row.expand = expanded;
    },
  },
});
const [MenuDrawer, drawerApi] = useVbenDrawer({
  connectedComponent: menuDrawer,
});

function handleAdd() {
  drawerApi.setData({});
  drawerApi.open();
}

function handleSubAdd(row: Menu) {
  const { id } = row;
  drawerApi.setData({ id, update: false });
  drawerApi.open();
}

async function handleEdit(record: Menu) {
  drawerApi.setData({ id: record.id, update: true });
  drawerApi.open();
}

async function handleDelete(row: Menu) {
  await menuRemove([row.id]);
  await tableApi.query();
}

function removeConfirmTitle(row: Menu) {
  const menu_name = $t(row.menu_name);
  return `是否确认删除 [${menu_name}] ?`;
}

/**
 * 编辑/添加成功后刷新表格
 */
async function afterEditOrAdd() {
  tableApi.query();
}

/**
 * 全部展开/折叠
 * @param expand 是否展开
 */
function setExpandOrCollapse(expand: boolean) {
  eachTree(tableApi.grid.getData(), (item) => (item.expand = expand));
  tableApi.grid?.setAllTreeExpand(expand);
}

/**
 * 与后台逻辑相同
 * 只有租户管理和超级管理能访问菜单管理
 * 注意: 只有超管才能对菜单进行`增删改`操作
 * 注意: 只有超管才能对菜单进行`增删改`操作
 * 注意: 只有超管才能对菜单进行`增删改`操作
 */
const { hasAccessByRoles } = useAccess();
const isAdmin = computed(() => {
  return hasAccessByRoles(['admin', 'superadmin']);
});
</script>

<template>
  <Page v-if="isAdmin" :auto-content-height="true">
    <BasicTable
      id="system-menu-table"
      table-title="菜单列表"
      table-title-help="双击展开/收起子菜单"
    >
      <template #toolbar-tools>
        <Space>
          <a-button @click="setExpandOrCollapse(false)">
            {{ $t('pages.common.collapse') }}
          </a-button>

          <a-button
            type="primary"
            v-access:code="['system:menu:add']"
            v-access:role="['superadmin']"
            @click="handleAdd"
          >
            {{ $t('pages.common.add') }}
          </a-button>
        </Space>
      </template>
      <template #action="{ row }">
        <Space>
          <ghost-button
            v-access:code="['system:menu:edit']"
            v-access:role="['superadmin']"
            @click="handleEdit(row)"
          >
            {{ $t('pages.common.edit') }}
          </ghost-button>
          <!-- '按钮类型'无法再添加子菜单 -->
          <ghost-button
            v-if="row.menuType !== 'F'"
            class="btn-success"
            v-access:code="['system:menu:add']"
            v-access:role="['superadmin']"
            @click="handleSubAdd(row)"
          >
            {{ $t('pages.common.add') }}
          </ghost-button>
          <Popconfirm
            :get-popup-container="getVxePopupContainer"
            placement="left"
            :title="removeConfirmTitle(row)"
            @confirm="handleDelete(row)"
          >
            <ghost-button
              danger
              v-access:code="['system:menu:remove']"
              v-access:role="['superadmin']"
              @click.stop=""
            >
              {{ $t('pages.common.delete') }}
            </ghost-button>
          </Popconfirm>
        </Space>
      </template>
    </BasicTable>
    <MenuDrawer @reload="afterEditOrAdd" />
  </Page>
  <Fallback v-else description="您没有菜单管理的访问权限" status="403" />
</template>

<style lang="scss">
#system-menu-table > .vxe-grid {
  --vxe-ui-table-row-current-background-color: hsl(var(--primary-100));

  html.dark & {
    --vxe-ui-table-row-current-background-color: hsl(var(--primary-800));
  }
}
</style>
