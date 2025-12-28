import type { DictType } from './dict-type-model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery, commonExport } from '#/api/helper';
import { pb } from '#/api/request';

const collection = pb.collection<DictType>('dict_type');

/**
 * 获取字典类型列表
 * @param params 请求参数
 * @returns list
 */
export function dictTypeList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return collection.getList(currentPage, pageSize, {
    filter,
    fields: 'id,dict_name,dict_type,create_time,remark',
    sort,
  });
}

/**
 * 导出字典类型列表
 * @param params 查询参数
 * @returns blob
 */
export function dictTypeExport(params?: PageQuery) {
  return commonExport('dict_type', buildingQuery(params), {
    type: 'collection',
  });
}

/**
 * 删除字典类型
 * @param dictIds 字典类型id数组
 * @returns void
 */
export function dictTypeRemove(dictIds: IDS) {
  return Promise.all(dictIds.map((id) => collection.delete(`${id}`)));
}
/**
 * 新增
 * @param data 表单参数
 * @returns void
 */
export function dictTypeAdd(data: Partial<DictType>) {
  return collection.create(data);
}

/**
 * 修改
 * @param data 表单参数
 * @returns void
 */
export function dictTypeUpdate(data: Partial<DictType>) {
  const id = `${data.id ?? ''}`;
  return collection.update(id, data);
}

/**
 * 查询详情
 * @param dictId 字典类型id
 * @returns 信息
 */
export function dictTypeInfo(dictId: ID) {
  return collection.getOne(`${dictId}`);
}

/**
 * 这个在ele用到 v5用不上
 * 下拉框  返回值和list一样
 * @returns options
 */
export function dictOptionSelectList() {
  return collection.getFullList();
}
