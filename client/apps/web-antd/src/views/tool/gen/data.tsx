import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import dayjs from 'dayjs';

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'tableName',
    label: '表名称',
  },
  {
    component: 'Input',
    fieldName: 'tableComment',
    label: '表描述',
  },
  {
    component: 'RangePicker',
    fieldName: 'createTime',
    label: '创建时间',
  },
];

export const columns: VxeGridProps['columns'] = [
  { type: 'checkbox', width: 60 },
  {
    field: 'name',
    title: '表名称',
  },
  {
    field: 'tableComment',
    title: '表描述',
  },
  {
    field: 'fields_length',
    title: '字段数量',
    formatter: ({ row }) => {
      return row.fields ? row.fields.length : 0;
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
    field: 'update_time',
    title: '更新时间',
    formatter: ({ cellValue }) => {
      return cellValue ? dayjs(cellValue).format('YYYY-MM-DD HH:mm:ss') : '';
    },
  },
  {
    field: 'action',
    fixed: 'right',
    slots: { default: 'action' },
    title: '操作',
    width: 300,
  },
];
