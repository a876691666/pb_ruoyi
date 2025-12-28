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
    fieldName: 'notice_title',
    label: '公告标题',
  },
  {
    component: 'Input',
    fieldName: 'create_by',
    label: '创建人',
  },
  {
    component: 'Select',
    componentProps: {
      getPopupContainer,
      options: getDictOptions(DictEnum.SYS_NOTICE_TYPE),
    },
    fieldName: 'notice_type',
    label: '公告类型',
  },
];

export const columns: VxeGridProps['columns'] = [
  { type: 'checkbox', width: 60 },
  {
    title: '公告标题',
    field: 'notice_title',
  },
  {
    title: '公告类型',
    field: 'notice_type',
    width: 120,
    slots: {
      default: ({ row }) => {
        return renderDict(row.notice_type, DictEnum.SYS_NOTICE_TYPE);
      },
    },
  },
  {
    title: '状态',
    field: 'status',
    width: 120,
    slots: {
      default: ({ row }) => {
        return renderDict(row.status, DictEnum.SYS_NOTICE_STATUS);
      },
    },
  },
  {
    title: '创建人',
    field: 'create_by',
    width: 150,
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
    fieldName: 'notice_id',
    label: '主键',
  },
  {
    component: 'Input',
    fieldName: 'notice_title',
    label: '公告标题',
    rules: 'required',
  },
  {
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: getDictOptions(DictEnum.SYS_NOTICE_STATUS),
      optionType: 'button',
    },
    defaultValue: '0',
    fieldName: 'status',
    label: '公告状态',
    rules: 'required',
    formItemClass: 'col-span-1',
  },
  {
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: getDictOptions(DictEnum.SYS_NOTICE_TYPE),
      optionType: 'button',
    },
    defaultValue: '1',
    fieldName: 'notice_type',
    label: '公告类型',
    rules: 'required',
    formItemClass: 'col-span-1',
  },
  {
    component: 'RichTextarea',
    componentProps: {
      width: '100%',
    },
    fieldName: 'notice_content',
    label: '公告内容',
    rules: 'required',
  },
];
