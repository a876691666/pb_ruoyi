<script setup lang="ts">
import type { Ref } from 'vue';

import type { GenInfo } from '#/api/tool/gen/model';

import { inject, onMounted } from 'vue';

import { useVbenForm } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { addFullName, listToTree } from '@vben/utils';

import { Col, Row } from 'ant-design-vue';

import { menuPageList } from '#/api/system/menu';

import { formSchema } from './basic';

/**
 * 从父组件注入
 */
const genInfoData = inject('genInfoData') as Ref<GenInfo>;

const [BasicForm, formApi] = useVbenForm({
  commonConfig: {
    componentProps: {
      class: 'w-full',
      formItemClass: 'col-span-1',
    },
    labelWidth: 150,
  },
  schema: formSchema(),
  showDefaultActions: false,
  wrapperClass: 'grid-cols-2',
});

/**
 * 树表需要用到的数据
 */
async function initTreeSelect(columns: any[]) {
  const options = columns.map((item) => {
    const label = `${item.name}${item.comment ? ` (${item.comment})` : ''}`;
    return { label, value: item.name };
  });
  formApi.updateSchema([
    {
      componentProps: {
        options,
      },
      fieldName: 'options.tree_code',
    },
    {
      componentProps: {
        options,
      },
      fieldName: 'options.tree_parent_code',
    },
    {
      componentProps: {
        options,
      },
      fieldName: 'options.tree_name',
    },
  ]);
}

/**
 * 加载菜单选择
 */
async function initMenuSelect() {
  const list = await menuPageList();
  const tree = listToTree(list, { id: 'id', pid: 'parent_id' });
  const treeData = [
    {
      fullName: $t('menu.root'),
      id: 0,
      menu_name: $t('menu.root'),
      children: tree,
    },
  ];
  addFullName(treeData, 'menu_name', ' / ');

  formApi.updateSchema([
    {
      componentProps: {
        fieldNames: {
          label: 'menu_name',
          value: 'id',
        },
        // 设置弹窗滚动高度 默认256
        listHeight: 300,
        treeData,
        treeDefaultExpandAll: false,
        // 默认展开的树节点
        treeDefaultExpandedKeys: [0],
        treeLine: { showLeafIcon: false },
        treeNodeLabelProp: 'fullName',
      },
      fieldName: 'options.parent_menu_id',
    },
  ]);
}

onMounted(async () => {
  const info = genInfoData.value;
  await formApi.setValues(info);
  await Promise.all([initTreeSelect(info.fields), initMenuSelect()]);
});

/**
 * 校验表单
 */
async function validateForm() {
  const { valid } = await formApi.validate();
  if (!valid) {
    return false;
  }
  return true;
}

/**
 * 获取表单值
 */
async function getFormValues() {
  return await formApi.getValues();
}

defineExpose({
  validateForm,
  getFormValues,
});
</script>

<template>
  <Row justify="center">
    <Col v-bind="{ xs: 24, sm: 24, md: 20, lg: 16, xl: 16 }">
      <BasicForm />
    </Col>
  </Row>
</template>
