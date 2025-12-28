import type { GenInfo } from '#/api/tool/gen/model';

import { cloneDeep } from '@vben/utils';

// 将下划线改为小驼峰
function toCamelCase(str: string) {
  return str.replaceAll(/_([a-z])/g, (g) => g[1]?.toUpperCase() ?? g);
}

/**
 * 序列化代码生成配置数据
 * 将前端表单数据转换为后端接口需要的格式
 * @param genInfoData 原始的生成配置数据
 * @param formValues 表单值
 * @param tableRecords 表格记录
 * @returns 序列化后的请求数据
 */
export function serializeGenData(
  genInfoData: GenInfo,
  formValues: any,
  tableRecords: any[],
) {
  const requestData = cloneDeep(genInfoData);

  // 合并表单数据
  Object.assign(requestData, formValues);

  // 从表格获取最新的列配置
  requestData.fields = tableRecords;

  return requestData;
}

/**
 * 反序列化代码生成配置数据
 * 将后端接口返回的数据转换为前端表单需要的格式
 * @param data 后端返回的数据
 * @returns 反序列化后的数据
 */
export function deserializeGenData(data: any) {
  // 如果需要反序列化逻辑，可以在这里添加
  return data;
}

// 系统字段列表 - 这些字段通常由系统自动维护
const systemFields = new Set([
  'collectionId',
  'collectionName',
  'create_by',
  'create_dept',
  'created',
  'id',
  'tenant_id',
  'update_by',
  'updated',
]);

// 隐藏字段列表 - 不在表格中显示
const hiddenFields = new Set([
  'collectionId',
  'collectionName',
  'create_by',
  'create_dept',
  'created',
  'remark',
  'tenant_id',
  'update_by',
  'updated',
]);

// 主键字段列表
const pkFields = new Set(['id']);

/**
 * 判断字段是否参与编辑
 * @param name 字段名
 * @returns 是否可编辑
 */
function getEdit(name: string) {
  if (systemFields.has(name)) {
    return false;
  }
  return true;
}

/**
 * 判断字段是否参与新增
 * @param name 字段名
 * @returns 是否参与新增
 */
function getInsert(name: string) {
  if (systemFields.has(name)) {
    return false;
  }
  return true;
}

/**
 * 判断字段是否在列表显示
 * @param name 字段名
 * @returns 是否显示
 */
function getList(name: string) {
  if (hiddenFields.has(name)) {
    return false;
  }
  return true;
}

/**
 * 判断字段是否作为查询条件
 * @param name 字段名
 * @returns 是否作为查询条件
 */
function getQuery(name: string) {
  if (systemFields.has(name)) {
    return false;
  }
  if (name === 'remark' || name === 'content' || name === 'description') {
    return false;
  }
  return true;
}

/**
 * 判断字段是否必填
 * @param name 字段名
 * @param _type 字段类型
 * @returns 是否必填
 */
function getRequired(name: string, _type?: string) {
  // 系统字段通常不需要用户填写
  if (systemFields.has(name)) {
    return false;
  }
  // 备注等字段通常不必填
  if (name === 'remark' || name === 'description' || name === 'note') {
    return false;
  }
  return false; // 默认不必填，由用户根据业务需要调整
}

/**
 * 判断字段是否为主键
 * @param name 字段名
 * @returns 是否为主键
 */
function getPk(name: string) {
  return pkFields.has(name);
}

/**
 * 获取查询方式
 * @param name 字段名
 * @param type 字段类型
 * @returns 查询方式
 */
function getQueryType(name: string, type?: string) {
  // 名称、标题等字段使用模糊查询
  if (
    name.includes('name') ||
    name.includes('title') ||
    name.includes('nickname')
  ) {
    return 'LIKE';
  }
  // 时间字段使用范围查询
  if (
    type?.includes('_date') ||
    type?.includes('_time') ||
    name.includes('_time') ||
    name.includes('_date')
  ) {
    return 'BETWEEN';
  }
  // 默认使用等于查询
  return 'EQ';
}

/**
 * 获取HTML显示类型
 * @param name 字段名
 * @param type 字段类型
 * @returns HTML显示类型
 */
function getHtmlType(name: string, type?: string) {
  // 长文本使用textarea
  if (name === 'remark' || name === 'content' || name === 'description') {
    return 'textarea';
  }
  // 时间字段
  if (
    type?.includes('_date') ||
    type?.includes('_time') ||
    name.includes('_time') ||
    name.includes('_date')
  ) {
    return 'datetime';
  }
  // 数字类型
  if (type?.includes('number') || type?.includes('int')) {
    return 'input';
  }
  // 布尔类型
  if (type?.includes('bool')) {
    return 'radio';
  }
  // 默认使用input
  return 'input';
}

/**
 * 序列化导入表数据
 * 将前端表数据转换为后端接口需要的格式
 * @param table 表数据
 * @returns 序列化后的表数据
 */
export function serializeImportTableData(table: any) {
  return {
    business_name: toCamelCase(table.name),
    comment: table.comment || '',
    fields: table.fields.map((field: any, index: number) => ({
      comment: field.comment || '',
      edit: field.primaryKey ? true : getEdit(field.name),
      htmlType: getHtmlType(field.name, field.type),
      insert: getInsert(field.name),
      list: field.primaryKey ? true : getList(field.name),
      name: field.name,
      pk: getPk(field.name),
      query: getQuery(field.name),
      queryType: getQueryType(field.name, field.type),
      required: getRequired(field.name, field.type),
      sort: index,
      type: field.type,
    })),
    function_name: table.name,
    module_name: 'system',
    name: table.name,
    tpl_category: 'crud',
  };
}
