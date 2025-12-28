import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { DictEnum } from '@vben/constants';
import { getPopupContainer } from '@vben/utils';

import dayjs from 'dayjs';

import { z } from '#/adapter/form';
import { getDictOptions } from '#/utils/dict';
import { renderDict } from '#/utils/render';

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'dept_name',
    label: '部门名称',
  },
  {
    component: 'Select',
    componentProps: {
      getPopupContainer,
      options: getDictOptions(DictEnum.SYS_NORMAL_DISABLE),
    },
    fieldName: 'status',
    label: '部门状态',
  },
];

export const columns: VxeGridProps['columns'] = [
  {
    field: 'dept_name',
    title: '部门名称',
    treeNode: true,
  },
  {
    field: 'dept_category',
    title: '类别编码',
  },
  {
    field: 'order_num',
    title: '排序',
  },
  {
    field: 'status',
    title: '状态',
    slots: {
      default: ({ row }) => {
        return renderDict(row.status, DictEnum.SYS_NORMAL_DISABLE);
      },
    },
  },
  {
    field: 'create_time',
    title: '创建时间',
    formatter: ({ cellValue }) => {
      return cellValue ? dayjs(cellValue).format('YYYY-MM-DD HH:mm:ss') : '';
    },
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
  },
  {
    component: 'TreeSelect',
    componentProps: {
      getPopupContainer,
    },
    dependencies: {
      show: (model) => model.parent_id !== 0,
      triggerFields: ['parent_id'],
    },
    fieldName: 'parent_id',
    label: '上级部门',
    rules: 'selectRequired',
  },
  {
    component: 'Input',
    fieldName: 'dept_name',
    label: '部门名称',
    rules: 'required',
  },
  {
    component: 'InputNumber',
    fieldName: 'order_num',
    label: '显示排序',
    rules: 'required',
    defaultValue: 0,
  },
  {
    component: 'Input',
    fieldName: 'dept_category',
    label: '类别编码',
  },
  {
    component: 'Select',
    componentProps: {
      // 选中了就只能修改 不能重置为无负责人
      allowClear: false,
      getPopupContainer,
    },
    fieldName: 'leader',
    label: '负责人',
  },
  {
    component: 'Input',
    fieldName: 'phone',
    label: '联系电话',
    rules: z
      .string()
      .regex(/^1[3,4578]\d{9}$/, { message: '请输入正确的手机号' })
      .optional()
      .or(z.literal('')),
  },
  {
    component: 'Input',
    fieldName: 'email',
    label: '邮箱',
    rules: z
      .string()
      .email({ message: '请输入正确的邮箱' })
      .optional()
      .or(z.literal('')),
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
    label: '状态',
  },
];
