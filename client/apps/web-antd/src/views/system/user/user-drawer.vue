<script setup lang="ts">
import type { Role } from '#/api/system/user/model';

import { computed, h, onMounted, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { addFullName, cloneDeep, getPopupContainer } from '@vben/utils';

import { Tag } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { configInfoByKey } from '#/api/system/config';
import { postOptionSelect } from '#/api/system/post';
import {
  findUserInfo,
  getDeptTree,
  userAdd,
  userUpdate,
} from '#/api/system/user';
import { defaultFormValueGetter, useBeforeCloseDiff } from '#/utils/popup';
import { authScopeOptions } from '#/views/system/role/data';

import { drawerSchema } from './data';

const emit = defineEmits<{ reload: [] }>();

const isUpdate = ref(false);
const title = computed(() => {
  return isUpdate.value ? $t('pages.common.edit') : $t('pages.common.add');
});

const [BasicForm, formApi] = useVbenForm({
  commonConfig: {
    formItemClass: 'col-span-2',
    componentProps: {
      class: 'w-full',
    },
    labelWidth: 80,
  },
  schema: drawerSchema(),
  showDefaultActions: false,
  wrapperClass: 'grid-cols-2',
});

/**
 * 生成角色的自定义label
 * 也可以用option插槽来做
 * renderComponentContent: () => ({
    option: ({value, label, [disabled, key, title]}) => '',
  }),
 */
function genRoleOptionlabel(role: Role) {
  const found = authScopeOptions.find((item) => item.value === role.data_scope);
  if (!found) {
    return role.role_name;
  }
  return h('div', { class: 'flex items-center gap-[6px]' }, [
    h('span', null, role.role_name),
    h(Tag, { color: found.color }, () => found.label),
  ]);
}

/**
 * 岗位的加载
 */
async function setupPostOptions(deptId: number | string) {
  const postListResp = await postOptionSelect(deptId);
  const options = postListResp.map((item) => ({
    label: item.post_name,
    value: item.id,
  }));
  const placeholder = options.length > 0 ? '请选择' : '该部门下暂无岗位';
  formApi.updateSchema([
    {
      componentProps: { options, placeholder },
      fieldName: 'post_ids',
    },
  ]);
}

/**
 * 初始化部门选择
 */
async function setupDeptSelect() {
  // updateSchema
  const deptTree = await getDeptTree();
  // 选中后显示在输入框的值 即父节点 / 子节点
  addFullName(deptTree, 'dept_name', ' / ');
  formApi.updateSchema([
    {
      componentProps: (formModel) => ({
        class: 'w-full',
        fieldNames: {
          key: 'id',
          label: 'dept_name',
          value: 'id',
          children: 'children',
        },
        getPopupContainer,
        async onSelect(dept_id: number | string) {
          /** 根据部门ID加载岗位 */
          await setupPostOptions(dept_id);
          /** 变化后需要重新选择岗位 */
          formModel.postIds = [];
        },
        placeholder: '请选择',
        showSearch: true,
        treeData: deptTree,
        treeDefaultExpandAll: true,
        treeLine: { showLeafIcon: false },
        // 筛选的字段
        treeNodeFilterProp: 'dept_name',
        // 选中后显示在输入框的值
        treeNodeLabelProp: 'fullName',
      }),
      fieldName: 'dept_id',
    },
  ]);
}

const defaultPassword = ref('');
onMounted(async () => {
  const password = await configInfoByKey('sys.user.initPassword');
  if (password) {
    defaultPassword.value = password;
  }
});

/**
 * 新增时候 从参数设置获取默认密码
 */
async function loadDefaultPassword(update: boolean) {
  if (!update && defaultPassword.value) {
    formApi.setFieldValue('password', defaultPassword.value);
  }
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
      // 需要重置岗位选择
      formApi.updateSchema([
        {
          componentProps: { options: [], placeholder: '请先选择部门' },
          fieldName: 'post_ids',
        },
      ]);
      return null;
    }
    drawerApi.drawerLoading(true);

    const { id } = drawerApi.getData() as { id?: number | string };
    isUpdate.value = !!id;
    /** update时 禁用用户名修改 不显示密码框 */
    formApi.updateSchema([
      { componentProps: { disabled: isUpdate.value }, fieldName: 'user_name' },
      {
        dependencies: { if: () => !isUpdate.value, triggerFields: ['id'] },
        fieldName: 'password',
      },
    ]);
    // 更新 && 赋值
    const { post_ids, posts, role_ids, roles, user } = await findUserInfo(id);
    const postOptions = (posts ?? []).map((item) => ({
      label: item.post_name,
      value: item.id,
    }));
    formApi.updateSchema([
      {
        componentProps: {
          // title用于选中后回填到输入框 默认为label
          optionLabelProp: 'title',
          options: roles.map((item) => ({
            label: genRoleOptionlabel(item),
            // title用于选中后回填到输入框 默认为label
            title: item.role_name,
            value: item.id,
          })),
        },
        fieldName: 'role_ids',
      },
      {
        componentProps: {
          options: postOptions,
        },
        fieldName: 'post_ids',
      },
    ]);

    // 部门选择、初始密码及用户相关操作并行处理
    const promises = [setupDeptSelect(), loadDefaultPassword(isUpdate.value)];
    if (user) {
      promises.push(
        // 添加基础信息
        formApi.setValues(user),
        // 添加角色和岗位
        formApi.setFieldValue('post_ids', post_ids),
        formApi.setFieldValue('role_ids', role_ids),
        // 更新时不会触发onSelect 需要手动调用
        setupPostOptions(user.dept_id),
      );
    }
    // 并行处理 重构后会带来10-50ms的优化
    await Promise.all(promises);
    await markInitialized();

    drawerApi.drawerLoading(false);
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
    await (isUpdate.value ? userUpdate(data) : userAdd(data));
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
  formApi.resetForm();
  resetInitialized();
}
</script>

<template>
  <BasicDrawer :title="title" class="w-[600px]">
    <BasicForm />
  </BasicDrawer>
</template>
