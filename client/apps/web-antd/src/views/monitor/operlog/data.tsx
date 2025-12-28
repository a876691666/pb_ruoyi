import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { DictEnum } from '@vben/constants';

import { getDictOptions } from '#/utils/dict';
import { renderDict } from '#/utils/render';

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'title',
    label: '系统模块',
  },
  {
    component: 'Input',
    fieldName: 'oper_name',
    label: '操作人员',
  },
  {
    component: 'Select',
    componentProps: {
      options: getDictOptions(DictEnum.SYS_OPER_TYPE),
    },
    fieldName: 'business_type',
    label: '操作类型',
  },
  {
    component: 'Input',
    fieldName: 'oper_ip',
    label: '操作IP',
  },
  {
    component: 'Select',
    componentProps: {
      options: getDictOptions(DictEnum.SYS_COMMON_STATUS),
    },
    fieldName: 'status',
    label: '状态',
  },
  {
    component: 'RangePicker',
    fieldName: 'create_time',
    label: '操作时间',
    componentProps: {
      valueFormat: 'YYYY-MM-DD HH:mm:ss',
    },
  },
];

export const columns: VxeGridProps['columns'] = [
  { type: 'checkbox', width: 60 },
  { field: 'title', title: '系统模块' },
  {
    title: '操作类型',
    field: 'business_type',
    slots: {
      default: ({ row }) => {
        return renderDict(row.business_type, DictEnum.SYS_OPER_TYPE);
      },
    },
  },
  { field: 'oper_name', title: '操作人员' },
  { field: 'oper_ip', title: 'IP地址' },
  { field: 'oper_location', title: 'IP信息' },
  {
    field: 'status',
    title: '操作状态',
    slots: {
      default: ({ row }) => {
        return renderDict(row.status, DictEnum.SYS_COMMON_STATUS);
      },
    },
  },
  { field: 'oper_time', title: '操作日期', sortable: true },
  {
    field: 'cost_time',
    title: '操作耗时',
    sortable: true,
    formatter({ cellValue }) {
      return `${cellValue} ms`;
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
