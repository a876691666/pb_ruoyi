import type { DeptTree } from '../user/model';
import type { Post } from './model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery, commonExport } from '#/api/helper';
import { pb, requestClient } from '#/api/request';

const collection = pb.collection<Post>('post');

// enum Api {
//   postList = '/system/post/list',
//   postSelect = '/system/post/optionselect',
//   root = '/system/post',
// }

/**
 * 获取岗位列表
 * @param params 参数
 * @returns Post[]
 */
export function postList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return collection.getList(currentPage, pageSize, {
    filter,
    sort,
  });
}

/**
 * 导出岗位信息
 * @param params 查询参数
 * @returns blob
 */
export function postExport(params?: PageQuery) {
  return commonExport('post', buildingQuery(params), { type: 'collection' });
}

/**
 * 查询岗位信息
 * @param postId id
 * @returns 岗位信息
 */
export function postInfo(postId: ID) {
  return collection.getOne(`${postId}`);
}

/**
 * 岗位新增
 * @param data 参数
 * @returns void
 */
export function postAdd(data: Partial<Post>) {
  return collection.create(data);
}

/**
 * 岗位更新
 * @param data 参数
 * @returns void
 */
export function postUpdate(data: Partial<Post>) {
  return collection.update(`${data.id}`, data);
}

/**
 * 岗位删除
 * @param postIds ids
 * @returns void
 */
export function postRemove(postIds: IDS) {
  return collection.delete(`${postIds}`);
}

/**
 * 根据部门id获取岗位下拉列表
 * @param deptId 部门id
 * @returns 岗位
 */
export function postOptionSelect(deptId: ID) {
  return collection.getFullList({
    filter: `dept_id="${deptId}"`,
  });
}

/**
 * 岗位专用 - 获取部门树
 * @returns 部门树
 */
export function postDeptTreeSelect() {
  return requestClient.get<DeptTree[]>('/system/post/deptTree');
}
