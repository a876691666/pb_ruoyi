import type { FormSchemaGetter } from '#/adapter/form';

import { getPopupContainer } from '@vben/utils';

export const formSchema: FormSchemaGetter = () => [
  {
    component: 'Divider',
    componentProps: {
      orientation: 'left',
    },
    fieldName: 'divider1',
    formItemClass: 'col-span-2',
    label: '基本信息',
  },
  {
    component: 'Input',
    fieldName: 'name',
    label: '表名称',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'comment',
    label: '表描述',
    rules: 'required',
  },
  {
    component: 'Divider',
    componentProps: {
      orientation: 'left',
    },
    fieldName: 'divider2',
    formItemClass: 'col-span-2',
    label: '生成信息',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      getPopupContainer,
      options: [
        { label: '单表(增删改查)', value: 'crud' },
        { label: '树表(增删改查)', value: 'tree' },
      ],
    },
    defaultValue: 'crud',
    fieldName: 'tpl_category',
    label: '模板类型',
    rules: 'selectRequired',
  },
  {
    component: 'Select',
    componentProps: {
      getPopupContainer,
    },
    dependencies: {
      show: (values) => values.tpl_category === 'tree',
      triggerFields: ['tpl_category'],
    },
    fieldName: 'options.tree_code',
    helpMessage: '树节点显示的编码字段名， 如: dept_id (相当于id)',
    label: '树编码字段',
    rules: 'selectRequired',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
    },
    dependencies: {
      show: (values) => values.tpl_category === 'tree',
      triggerFields: ['tpl_category'],
    },
    fieldName: 'options.tree_parent_code',
    help: '树节点显示的父编码字段名， 如: parent_Id (相当于parentId)',
    label: '树父编码字段',
    rules: 'selectRequired',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
    },
    dependencies: {
      show: (values) => values.tpl_category === 'tree',
      triggerFields: ['tpl_category'],
    },
    fieldName: 'options.tree_name',
    help: '树节点的显示名称字段名， 如: dept_name (相当于label)',
    label: '树名称字段',
    rules: 'selectRequired',
  },
  {
    component: 'Input',
    fieldName: 'module_name',
    help: '可理解为子系统名，例如 system',
    label: '生成模块名',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'business_name',
    help: '可理解为功能英文名，例如 user',
    label: '生成业务名',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'function_name',
    help: '用作类描述，例如 用户',
    label: '生成功能名',
    rules: 'required',
  },
  {
    component: 'TreeSelect',
    componentProps: {
      allowClear: false,
      getPopupContainer,
    },
    defaultValue: 0,
    fieldName: 'options.parent_menu_id',
    label: '上级菜单',
  },
  {
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: [
        { label: 'modal弹窗', value: 'modal' },
        { label: 'drawer抽屉', value: 'drawer' },
      ],
      optionType: 'button',
    },
    help: '自定义功能, 需要后端支持',
    defaultValue: 'modal',
    fieldName: 'options.popup_type',
    label: '弹窗组件类型',
  },
  {
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: [
        { label: 'useVbenForm', value: 'useForm' },
        { label: 'antd原生表单', value: 'native' },
      ],
      optionType: 'button',
    },
    help: '自定义功能, 需要后端支持\n复杂(布局, 联动等)表单建议用antd原生表单',
    defaultValue: 'useForm',
    fieldName: 'options.form_type',
    label: '生成表单类型',
  },
  {
    component: 'Textarea',
    fieldName: 'remark',
    formItemClass: 'col-span-2 items-baseline',
    label: '备注',
  },
];
