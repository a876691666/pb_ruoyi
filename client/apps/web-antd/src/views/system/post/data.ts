import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { DictEnum } from '@vben/constants';
import { getPopupContainer } from '@vben/utils';

import { getDictOptions } from '#/utils/dict';
import { renderDict } from '#/utils/render';

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'post_code',
    label: '岗位编码',
  },
  {
    component: 'Input',
    fieldName: 'post_name',
    label: '岗位名称',
  },
  {
    component: 'Select',
    componentProps: {
      getPopupContainer,
      options: getDictOptions(DictEnum.SYS_NORMAL_DISABLE),
    },
    fieldName: 'status',
    label: '状态',
  },
];

export const columns: VxeGridProps['columns'] = [
  { type: 'checkbox', width: 60 },
  {
    title: '岗位编码',
    field: 'post_code',
  },
  {
    title: '类别编码',
    field: 'post_category',
  },
  {
    title: '岗位名称',
    field: 'post_name',
  },
  {
    title: '排序',
    field: 'post_sort',
  },
  {
    title: '状态',
    field: 'status',
    slots: {
      default: ({ row }) => {
        return renderDict(row.status, DictEnum.SYS_NORMAL_DISABLE);
      },
    },
  },
  {
    title: '创建时间',
    field: 'create_time',
  },
  {
    field: 'action',
    fixed: 'right',
    slots: { default: 'action' },
    title: '操作',
    resizable: false,
    width: 'auto',
  },
];

export const drawerSchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    dependencies: {
      show: () => false,
      triggerFields: [''],
    },
    fieldName: 'id',
    label: 'postId',
  },
  {
    component: 'TreeSelect',
    componentProps: {
      getPopupContainer,
    },
    fieldName: 'dept_id',
    label: '所属部门',
    rules: 'selectRequired',
  },
  {
    component: 'Input',
    fieldName: 'post_name',
    label: '岗位名称',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'post_code',
    label: '岗位编码',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'post_category',
    label: '类别编码',
  },
  {
    component: 'InputNumber',
    fieldName: 'post_sort',
    label: '岗位排序',
    rules: 'required',
    defaultValue: 0,
  },
  {
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: getDictOptions(DictEnum.SYS_NORMAL_DISABLE),
      optionType: 'button',
    },
    defaultValue: '0',
    fieldName: 'status',
    label: '岗位状态',
    rules: 'required',
  },
  {
    component: 'Textarea',
    fieldName: 'remark',
    formItemClass: 'items-start',
    label: '备注',
  },
];
