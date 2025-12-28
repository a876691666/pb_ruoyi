import type { Dept } from './model';

import type { ID, PageQuery } from '#/api/common';

import { buildingQuery } from '#/api/helper';
import { pb } from '#/api/request';

const collection = pb.collection<Dept>('dept');

/**
 * 部门列表（树形结构，获取全部数据）
 * @returns list
 */
export function deptList(params?: PageQuery) {
  const { filter, sort } = buildingQuery(params);
  return collection.getFullList({
    filter,
    sort,
  });
}

/**
 * 查询部门列表（排除节点）
 * 该接口属于特殊端点，保留原始请求方式
 * @param id 部门ID
 */
export function deptNodeList(id: ID) {
  return collection.getFullList({
    filter: `id != "${id}" && ancestors !~ ",${id}"`,
  });
}

/**
 * 部门详情
 * @param id 部门id
 * @returns 部门信息
 */
export function deptInfo(id: ID) {
  return collection.getOne(`${id}`) as Promise<Dept>;
}

/**
 * 部门新增
 * @param data 参数
 */
export function deptAdd(data: Partial<Dept>) {
  return collection.create(data);
}

/**
 * 部门更新
 * @param data 参数
 */
export function deptUpdate(data: Partial<Dept>) {
  return collection.update(`${data.id}`, data);
}

/**
 * 注意这里只有单删除
 * @param id ID
 */
export function deptRemove(id: ID) {
  return collection.delete(`${id}`);
}
