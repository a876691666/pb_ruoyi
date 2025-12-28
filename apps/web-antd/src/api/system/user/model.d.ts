import type { BaseCollectionModel } from 'pocketbase';

/**
 * @description: 用户导入
 * @param updateSupport 是否覆盖数据
 * @param file excel文件
 */
export interface UserImportParam {
  update_support: boolean;
  file: Blob | File;
}

/**
 * @description: 重置密码
 */
export interface ResetPwdParam {
  id: string;
  password: string;
}

export interface Dept extends BaseCollectionModel {
  id: number;
  parent_id: number;
  parent_name?: string;
  ancestors: string;
  dept_name: string;
  order_num: number;
  leader: string;
  phone?: string;
  email?: string;
  status: string;
  create_time?: string;
}

export interface Role extends BaseCollectionModel {
  id: string;
  role_name: string;
  role_key: string;
  role_sort: number;
  data_scope: string;
  menu_check_strictly?: boolean;
  dept_check_strictly?: boolean;
  status: string;
  remark: string;
  create_time?: string;
  flag: boolean;
}

export interface User extends BaseCollectionModel {
  id: string;
  tenant_id: string;
  dept_id: number;
  user_name: string;
  nick_name: string;
  user_type: string;
  email: string;
  phonenumber: string;
  sex: string;
  avatar?: string;
  status: string;
  login_ip: string;
  login_date: string;
  remark: string;
  create_time: string;
  dept: Dept;
  roles: Role[];
  role_ids?: string[];
  post_ids?: number[];
  role_id: string;
  dept_name: string;
}

export interface Post extends BaseCollectionModel {
  post_id: number;
  post_code: string;
  post_name: string;
  post_sort: number;
  status: string;
  remark: string;
  create_time: string;
}

/**
 * @description 用户信息
 * @param user 用户个人信息
 * @param roleIds 角色IDS 不传id为空
 * @param roles 所有的角色
 * @param postIds 岗位IDS 不传id为空
 * @param posts 所有的岗位
 */
export interface UserInfoResponse {
  user?: User;
  role_ids?: string[];
  roles: Role[];
  post_ids?: number[];
  posts?: Post[];
}

/**
 * @description: 部门树
 */
export interface DeptTree {
  id: number;
  /**
   * antd组件必须要这个属性 实际是没有这个属性的
   */
  key: string;
  parent_id: number;
  label: string;
  weight: number;
  children?: DeptTree[];
}

export interface DeptTreeData {
  id: number;
  label: string;
  children?: DeptTreeData[];
}
