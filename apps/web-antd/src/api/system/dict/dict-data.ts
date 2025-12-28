import type { DictData } from './dict-data-model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery, commonExport } from '#/api/helper';
import { pb } from '#/api/request';

const collection = pb.collection<DictData>('dict_data');

// enum Api {
//   dictDataList = '/system/dict/data/list',
//   root = '/system/dict/data',
// }

/**
 * 主要是DictTag组件使用
 * @param dictType 字典类型
 * @returns 字典数据
 */
export function dictDataInfo(dictType: string) {
  return collection.getFullList({
    filter: `dict_type = "${dictType}"`,
    fields:
      'id,dict_sort,dict_label,dict_value,css_class,list_class,is_default,remark,create_time',
    requestKey: `dict_data_info_${dictType}`,
  });
}

/**
 * 字典数据
 * @param params 查询参数
 * @returns 字典数据列表
 */
export function dictDataList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return collection.getList(currentPage, pageSize, {
    filter,
    fields:
      'id,dict_sort,dict_label,dict_value,dict_type,css_class,list_class,is_default,remark,create_time',
    sort,
  });
}

/**
 * 导出字典数据
 * @param params 查询参数
 * @returns blob
 */
export function dictDataExport(params?: PageQuery) {
  return commonExport('dict_data', buildingQuery(params), {
    type: 'collection',
  });
}

/**
 * 删除
 * @param dictIds 字典ID Array
 * @returns void
 */
export function dictDataRemove(dictIds: IDS) {
  return Promise.all(dictIds.map((id) => collection.delete(`${id}`)));
}

/**
 * 新增
 * @param data 表单参数
 * @returns void
 */
export function dictDataAdd(data: Partial<DictData>) {
  return collection.create(data);
}

/**
 * 修改
 * @param data 表单参数
 * @returns void
 */
export function dictDataUpdate(data: Partial<DictData>) {
  const id = `${data.id ?? ''}`;
  return collection.update(id, data);
}

/**
 * 查询字典数据详细
 * @param dictCode 字典编码
 * @returns 字典数据
 */
export function dictDetailInfo(dictCode: ID) {
  return collection.getOne(`${dictCode}`);
}
