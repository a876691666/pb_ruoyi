import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import dayjs from 'dayjs';

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'file_name',
    label: '文件名',
  },
  {
    component: 'Input',
    fieldName: 'original_name',
    label: '原名',
  },
  {
    component: 'Input',
    fieldName: 'file_suffix',
    label: '拓展名',
  },
  {
    component: 'RangePicker',
    fieldName: 'create_time',
    label: '创建时间',
  },
];

export const columns: VxeGridProps['columns'] = [
  { type: 'checkbox', width: 60 },
  {
    title: '文件名',
    field: 'file_name',
    showOverflow: true,
  },
  {
    title: '文件原名',
    field: 'original_name',
    showOverflow: true,
  },
  {
    title: '文件拓展名',
    field: 'file_suffix',
  },
  {
    title: '文件预览',
    field: 'url',
    showOverflow: true,
    slots: { default: 'url' },
  },
  {
    title: '创建时间',
    field: 'create_time',
    formatter: ({ cellValue }) => {
      return cellValue ? dayjs(cellValue).format('YYYY-MM-DD HH:mm:ss') : '';
    },
    sortable: true,
  },
  {
    title: '上传人',
    field: 'create_by_name',
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
