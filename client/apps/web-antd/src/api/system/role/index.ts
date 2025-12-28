import type { DeptResp, Role } from './model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery, commonExport } from '#/api/helper';
import { Ands, pb } from '#/api/request';
import { buildTree } from '#/utils/tree';

const roleCollection = pb.collection<Role>('role');
const userRoleCollection = pb.collection('user_role');
const usersCollection = pb.collection('users');
const roleDeptCollection = pb.collection('role_dept');
const deptCollection = pb.collection('dept');

/**
 * 查询角色分页列表
 * @param params 搜索条件
 * @returns 分页列表
 */
export function roleList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return roleCollection.getList(currentPage, pageSize, {
    filter,
    sort,
  });
}

/**
 * 导出角色信息
 * @param params 查询参数
 * @returns blob
 */
export function roleExport(params?: PageQuery) {
  return commonExport('role', buildingQuery(params), { type: 'collection' });
}

/**
 * 查询角色信息
 * @param id 角色id
 * @returns 角色信息
 */
export function roleInfo(id: ID) {
  return roleCollection.getOne(`${id}`);
}

/**
 * 角色新增
 * @param data 参数
 * @returns void
 */
export function roleAdd(data: Partial<Role>) {
  return roleCollection.create(data, {
    headers: {
      'X-Menu': 'true',
    },
  });
}

/**
 * 角色更新
 * @param data 参数
 * @returns void
 */
export function roleUpdate(data: Partial<Role>) {
  return roleCollection.update(`${data.id}`, data, {
    headers: {
      'X-Menu': 'true',
    },
  });
}

/**
 * 修改角色状态
 * @param data 参数
 * @returns void
 */
export function roleChangeStatus(data: Partial<Role>) {
  roleCollection.update(`${data.id}`, { status: data.status });
  return Promise.resolve();
}

/**
 * 角色删除
 * @param roleIds ids
 * @returns void
 */
export function roleRemove(roleIds: IDS) {
  return Promise.all(roleIds.map((id) => roleCollection.delete(`${id}`)));
}

/**
 * 更新数据权限
 * @param data
 * @returns void
 */
export function roleDataScope(data: any) {
  return roleCollection.update(`${data.id}`, data, {
    headers: {
      'X-Dept': 'true',
    },
  });
}

/**
 * 已分配角色的用户分页
 * @param params 请求参数
 * @returns 分页
 */
export function roleAllocatedList(params?: PageQuery) {
  const { sort, pageSize, currentPage } = buildingQuery(params);
  return userRoleCollection.getList(currentPage, pageSize, {
    filter: Ands([
      `role="${params?.params?.role_id}"`,
      params?.params?.user_name
        ? `user.user_name ~ "${params?.params?.user_name}"`
        : undefined,
      params?.params?.phonenumber
        ? `user.phonenumber ~ "${params?.params?.phonenumber}"`
        : undefined,
    ]),
    fields: 'id,expand.user',
    expand: 'user',
    sort,
  });
}

/**
 * 未授权的用户
 * @param params
 * @returns void
 */
export function roleUnallocatedList(params: any) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  // 添加额外的角色过滤条件
  const roleFilter = `user_role_via_user.role !="${params?.role_id}"`;
  const combinedFilter = filter ? `(${filter}) && ${roleFilter}` : roleFilter;

  return usersCollection.getList(currentPage, pageSize, {
    filter: combinedFilter,
    expand: 'user_role_via_user',
    sort,
  });
}

/**
 * 取消用户角色授权
 * @returns void
 */
export async function roleAuthCancel(id: string) {
  return userRoleCollection.delete(id);
}

/**
 * 批量取消授权
 * @param ids 角色ID
 * @returns void
 */
export function roleAuthCancelAll(ids: IDS) {
  return Promise.all(ids.map((id) => userRoleCollection.delete(`${id}`)));
}

/**
 * 批量授权用户
 * @param id 角色ID
 * @param userIds 用户ID集合
 * @returns void
 */
export function roleSelectAll(id: ID, userIds: IDS) {
  return Promise.all(
    userIds.map((userId) =>
      userRoleCollection.create(
        {
          role: id,
          user: userId,
        },
        { requestKey: `roleSelectAll-${id}-${userId}` },
      ),
    ),
  );
}

/**
 * 根据角色id获取部门树
 * @param id 角色id
 * @returns DeptResp
 */
export async function roleDeptTree(id: ID) {
  const roleDepts = await roleDeptCollection.getFullList<{
    dept: string;
  }>({
    filter: `role="${id}"`,
    fields: 'dept',
  });
  const checkedKeys = roleDepts.map((m) => m.dept);

  const deptsRaw = await deptCollection.getFullList<{
    dept_name: string;
    id: string;
    order_num?: number;
    parent_id?: null | string;
    status: string;
  }>({
    filter: `status="0"`,
    // 只取构建树所需字段，减少网络负载
    fields: 'id,parent_id,dept_name,status,order_num',
  });

  const depts = buildTree<any, DeptResp['depts'][number]>(deptsRaw, {
    getId: (d) => d.id,
    getParentId: (d) => d.parent_id ?? '0',
    rootPid: '0',
    getLabel: (d) => d.dept_name,
    getWeight: (d) => d.order_num ?? 0,
    assign: (d, node: any) => {
      node.key = d.id;
    },
  });

  return { checkedKeys, depts };
}
