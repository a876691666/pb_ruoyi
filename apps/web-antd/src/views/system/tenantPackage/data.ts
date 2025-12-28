import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'package_name',
    label: '套餐名称',
  },
];

export const columns: VxeGridProps['columns'] = [
  { type: 'checkbox', width: 60 },
  {
    title: '套餐名称',
    field: 'package_name',
  },
  {
    title: '备注',
    field: 'remark',
  },
  {
    title: '状态',
    field: 'status',
    slots: { default: 'status' },
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
    component: 'Radio',
    dependencies: {
      show: () => false,
      triggerFields: [''],
    },
    fieldName: 'menu_check_strictly',
  },
  {
    component: 'Input',
    fieldName: 'package_name',
    label: '套餐名称',
    rules: 'required',
  },
  {
    component: 'menuIds',
    defaultValue: [],
    fieldName: 'menu_ids',
    label: '关联菜单',
  },
  {
    component: 'Textarea',
    fieldName: 'remark',
    label: '备注',
  },
];
