<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';
import {
  addFullName,
  cloneDeep,
  getPopupContainer,
  listToTree,
} from '@vben/utils';

import { Input, Skeleton } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { menuAdd, menuInfo, menuList, menuUpdate } from '#/api/system/menu';
import { defaultFormValueGetter, useBeforeCloseDiff } from '#/utils/popup';

import { drawerSchema } from './data';

interface ModalProps {
  id?: number | string;
  update: boolean;
}

const emit = defineEmits<{ reload: [] }>();

const isUpdate = ref(false);
const title = computed(() => {
  return isUpdate.value ? $t('pages.common.edit') : $t('pages.common.add');
});
const loading = ref(false);

const [BasicForm, formApi] = useVbenForm({
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
    formItemClass: 'col-span-2',
    labelWidth: 90,
  },
  schema: drawerSchema(),
  showDefaultActions: false,
  wrapperClass: 'grid-cols-2',
});

async function setupMenuSelect() {
  // menu
  const menuArray = await menuList();
  // support i18n
  menuArray.forEach((item) => {
    item.menu_name = $t(item.menu_name);
  });
  // const folderArray = menuArray.filter((item) => item.menu_type === 'M');
  /**
   * 这里需要过滤掉按钮类型
   * 不允许在按钮下添加数据
   */
  const filteredList = menuArray.filter((item) => item.menu_type !== 'F');
  const menuTree = listToTree(filteredList, { id: 'id', pid: 'parent_id' });
  const fullMenuTree = [
    {
      id: 0,
      menu_name: $t('menu.root'),
      children: menuTree,
    },
  ];
  addFullName(fullMenuTree, 'menu_name', ' / ');

  formApi.updateSchema([
    {
      componentProps: {
        fieldNames: {
          label: 'menu_name',
          value: 'id',
        },
        getPopupContainer,
        // 设置弹窗滚动高度 默认256
        listHeight: 300,
        showSearch: true,
        treeData: fullMenuTree,
        treeDefaultExpandAll: false,
        // 默认展开的树节点
        treeDefaultExpandedKeys: [0],
        treeLine: { showLeafIcon: false },
        // 筛选的字段
        treeNodeFilterProp: 'menu_name',
        treeNodeLabelProp: 'fullName',
      },
      fieldName: 'parent_id',
    },
  ]);
}

const { onBeforeClose, markInitialized, resetInitialized } = useBeforeCloseDiff(
  {
    initializedGetter: defaultFormValueGetter(formApi),
    currentGetter: defaultFormValueGetter(formApi),
  },
);

const [BasicDrawer, drawerApi] = useVbenDrawer({
  onBeforeClose,
  onClosed: handleClosed,
  onConfirm: handleConfirm,
  async onOpenChange(isOpen) {
    if (!isOpen) {
      return null;
    }
    drawerApi.drawerLoading(true);
    loading.value = true;

    const { id, update } = drawerApi.getData() as ModalProps;
    isUpdate.value = update;

    if (id) {
      await formApi.setFieldValue('parent_id', id);
      setupMenuSelect();
      if (update) {
        // 没有依赖关系 同时加载
        const record = await menuInfo(id);
        await formApi.setValues(record);
      }
    } else {
      // 加载菜单树选择
      await setupMenuSelect();
    }
    await markInitialized();

    drawerApi.drawerLoading(false);
    loading.value = false;
  },
});

async function handleConfirm() {
  try {
    drawerApi.lock(true);
    const { valid } = await formApi.validate();
    if (!valid) {
      return;
    }
    const data = cloneDeep(await formApi.getValues());
    await (isUpdate.value ? menuUpdate(data) : menuAdd(data));
    resetInitialized();
    emit('reload');
    drawerApi.close();
  } catch (error) {
    console.error(error);
  } finally {
    drawerApi.lock(false);
  }
}

async function handleClosed() {
  await formApi.resetForm();
  resetInitialized();
}
</script>

<template>
  <BasicDrawer :title="title" class="w-[600px]">
    <Skeleton active v-if="loading" />
    <BasicForm v-show="!loading">
      <template #remark="slotProps">
        <div class="flex flex-col gap-2">
          <Input v-bind="slotProps" />
          <span class="text-[14px] leading-[1.5] text-black/45">
            在ele作为activePath使用 但是非json格式 v5无法使用
            建议自行在apps/web-antd/src/router/access.ts更改
          </span>
        </div>
      </template>
    </BasicForm>
  </BasicDrawer>
</template>
