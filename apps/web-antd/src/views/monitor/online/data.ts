import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import dayjs from 'dayjs';

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'dept_name',
    label: '部门名称',
  },
  {
    component: 'Input',
    fieldName: 'user_name',
    label: '用户账号',
  },
];

export const columns: VxeGridProps['columns'] = [
  {
    title: '登录账号',
    field: 'user_name',
  },
  {
    title: '部门名称',
    field: 'dept_name',
  },
  {
    title: 'IP地址',
    field: 'phonenumber',
  },
  {
    title: '登录时间',
    field: 'login_time',
    formatter: ({ cellValue }) => {
      return dayjs(cellValue).format('YYYY-MM-DD HH:mm:ss');
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
