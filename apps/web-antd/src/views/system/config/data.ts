import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { DictEnum } from '@vben/constants';
import { getPopupContainer } from '@vben/utils';

import dayjs from 'dayjs';

import { getDictOptions } from '#/utils/dict';
import { renderDict } from '#/utils/render';

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'name',
    label: '参数名称',
  },
  {
    component: 'Input',
    fieldName: 'key',
    label: '参数键名',
  },
  {
    component: 'Select',
    componentProps: {
      getPopupContainer,
      options: getDictOptions(DictEnum.SYS_YES_NO),
    },
    fieldName: 'type',
    label: '系统内置',
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
    title: '参数名称',
    field: 'name',
  },
  {
    title: '参数KEY',
    field: 'key',
  },
  {
    title: '参数Value',
    field: 'value',
    sortable: true,
  },
  {
    title: '系统内置',
    field: 'type',
    width: 120,
    slots: {
      default: ({ row }) => {
        return renderDict(row.type, DictEnum.SYS_YES_NO);
      },
    },
  },
  {
    title: '备注',
    field: 'remark',
    sortable: true,
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
    width: 'auto',
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
    label: '参数主键',
  },
  {
    component: 'Input',
    fieldName: 'name',
    label: '参数名称',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'key',
    label: '参数键名',
    rules: 'required',
  },
  {
    component: 'Textarea',
    formItemClass: 'items-start',
    fieldName: 'value',
    label: '参数键值',
    componentProps: {
      autoSize: true,
    },
    rules: 'required',
  },
  {
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: getDictOptions(DictEnum.SYS_YES_NO),
      optionType: 'button',
    },
    defaultValue: 'N',
    fieldName: 'type',
    label: '是否内置',
    rules: 'required',
  },
  {
    component: 'Textarea',
    fieldName: 'remark',
    formItemClass: 'items-start',
    label: '备注',
  },
];
