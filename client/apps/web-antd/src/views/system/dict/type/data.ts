import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import dayjs from 'dayjs';

import { z } from '#/adapter/form';

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'dict_name',
    label: '字典名称',
  },
  {
    component: 'Input',
    fieldName: 'dict_type',
    label: '字典类型',
  },
];

export const columns: VxeGridProps['columns'] = [
  { type: 'checkbox', width: 60 },
  {
    title: '字典名称',
    field: 'dict_name',
  },
  {
    title: '字典类型',
    field: 'dict_type',
  },
  {
    title: '备注',
    field: 'remark',
  },
  {
    title: '创建时间',
    field: 'create_time',
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
    width: 120,
  },
];

export const modalSchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    dependencies: {
      show: () => false,
      triggerFields: [''],
    },
    fieldName: 'id',
    label: '字典ID',
  },
  {
    component: 'Input',
    fieldName: 'dict_name',
    label: '字典名称',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'dict_type',
    help: '使用英文/下划线命名, 如:sys_normal_disable',
    label: '字典类型',
    rules: z
      .string()
      .regex(/^[a-z_]+$/i, { message: '字典类型只能使用英文/下划线命名' }),
  },
  {
    component: 'Textarea',
    fieldName: 'remark',
    label: '备注',
  },
];
