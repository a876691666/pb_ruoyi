import type { DeptTree, ResetPwdParam, User, UserImportParam } from './model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery, commonExport, ContentTypeEnum } from '#/api/helper';
import { pb, requestClient } from '#/api/request';

import { postList } from '../post';
import { roleList } from '../role';

const usersCollection = pb.collection<User>('users');
const userRoleCollection = pb.collection<{ role: string }>('user_role');
const userPostCollection = pb.collection<{ post: string }>('user_post');

enum Api {
  deptTree = '/system/user/deptTree',
  listDeptUsers = '/system/user/list/dept',
  root = '/system/user',
  userAuthRole = '/system/user/authRole',
  userExport = '/system/user/export',
  userImport = '/system/user/importData',
  userImportTemplate = '/system/user/importTemplate',
  userList = '/system/user/list',
  userResetPassword = '/system/user/resetPwd',
  userStatusChange = '/system/user/changeStatus',
}

/**
 *  获取用户列表
 * @param params
 * @returns User
 */
export function userList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return usersCollection.getList(currentPage, pageSize, {
    filter,
    sort,
  });
}

/**
 * 导出excel
 * @param params 查询参数
 * @returns blob
 */
export function userExport(params?: PageQuery) {
  return commonExport('users', buildingQuery(params), { type: 'collection' });
}

/**
 * 从excel导入用户
 * @param data
 * @returns void
 */
export function userImportData(data: UserImportParam) {
  return requestClient.post<{ code: number; msg: string }>(
    Api.userImport,
    data,
    {
      headers: {
        'Content-Type': ContentTypeEnum.FORM_DATA,
      },
      isTransformResponse: false,
    },
  );
}

/**
 * 下载用户导入模板
 * @returns blob
 */
export function downloadImportTemplate() {
  return requestClient.post<Blob>(
    Api.userImportTemplate,
    {},
    {
      isTransformResponse: false,
      responseType: 'blob',
    },
  );
}

/**
 * 可以不传ID  返回部门和角色options 需要获得原始数据
 * 不传ID时一定要带最后的/
 * @param userId 用户ID
 * @returns 用户信息
 */
export async function findUserInfo(userId?: ID) {
  let user;
  let role_ids: string[] = [];
  let post_ids: string[] = [];
  const [rolesResult, postResult] = await Promise.all([roleList(), postList()]);

  if (userId) {
    const [userResult, roleIdsResult, postIdsResult] = await Promise.all([
      usersCollection.getOne(`${userId}`),
      userRoleCollection.getFullList({
        filter: `user = "${userId}"`,
        fields: 'role',
      }),
      userPostCollection.getFullList({
        filter: `user = "${userId}"`,
        fields: 'post',
      }),
    ]);
    user = userResult;
    role_ids = roleIdsResult.map((item) => item.role);
    post_ids = postIdsResult.map((item) => item.post);
  }

  const roles = rolesResult.items;
  const posts = postResult.items;

  return {
    user,
    roles,
    posts,
    role_ids,
    post_ids,
  };
}

/**
 * 新增用户
 * @param data data
 * @returns void
 */
export function userAdd(data: Partial<User>) {
  return usersCollection.create(
    { ...data, passwordConfirm: (data as any).password },
    {
      headers: { 'X-Dept': 'true', 'X-Post': 'true', 'X-Role': 'true' },
    },
  );
}

/**
 * 更新用户
 * @param data data
 * @returns void
 */
export function userUpdate(data: Partial<User>) {
  return usersCollection.update(`${data.id}`, data, {
    headers: { 'X-Dept': 'true', 'X-Post': 'true', 'X-Role': 'true' },
  });
}

/**
 * 更新用户状态
 * @param data data
 * @returns void
 */
export function userStatusChange(data: Partial<User>) {
  const requestData = {
    id: data.id,
    status: data.status,
  };
  return usersCollection.update(`${data.id}`, requestData);
}

/**
 * 删除用户
 * @param userIds 用户ID数组
 * @returns void
 */
export function userRemove(userIds: IDS) {
  return Promise.all(userIds.map((id) => usersCollection.delete(`${id}`)));
}

/**
 * 重置用户密码 需要加密
 * @param data
 * @returns void
 */
export function userResetPassword(data: ResetPwdParam) {
  return requestClient.putWithMsg<void>(Api.userResetPassword, data);
}

/**
 * 获取部门树
 * @returns 部门树数组
 */
export function getDeptTree() {
  return requestClient.get<DeptTree[]>(Api.deptTree);
}

/**
 * 获取部门下的所有用户信息
 */
export function listUserByDeptId(deptId: ID) {
  return requestClient.get<User[]>(`${Api.listDeptUsers}/${deptId}`);
}
