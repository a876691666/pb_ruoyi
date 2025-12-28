import type { SysConfig } from './model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery, commonExport } from '#/api/helper';
import { pb } from '#/api/request';

const collection = pb.collection<SysConfig>('config');

/**
 * 新增参数设置
 * @param data 参数
 */
export function configAdd(data: Partial<SysConfig>) {
  return collection.create(data);
}

/**
 * 删除参数设置
 * @param configIds ids
 */
export function configRemove(configIds: IDS) {
  return Promise.all(configIds.map((id) => collection.delete(`${id}`)));
}

/**
 * 更新参数设置
 * @param data 参数
 */
export function configUpdate(data: Partial<SysConfig>) {
  return collection.update(`${data.id}`, data);
}

/**
 * 参数设置分页列表
 * @param params 请求参数
 * @returns 列表
 */
export function configList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return collection.getList(currentPage, pageSize, {
    filter,
    sort,
  });
}

/**
 * 参数设置详情
 * @param configId id
 */
export function configInfo(configId: ID) {
  return collection.getOne(`${configId}`);
}

/**
 * 导出参数设置
 * @param params 请求参数
 */
export function configExport(params?: PageQuery) {
  return commonExport('config', buildingQuery(params), { type: 'collection' });
}

/**
 * 获取配置信息
 * @param configKey configKey
 * @returns value
 */
export function configInfoByKey(configKey: string) {
  return pb
    .collection('config')
    .getFirstListItem(`key="${configKey}"`)
    .then((res) => res.value);
}
