import type { BaseCollectionModel } from 'pocketbase';

export interface Dept extends BaseCollectionModel {
  deptId: number;
  parentId: number;
  parentName?: any;
  ancestors: string;
  deptName: string;
  orderNum: number;
  leader: string;
  phone?: any;
  email: string;
  status: string;
  createTime?: any;
}

export interface Role extends BaseCollectionModel {
  roleId: number;
  roleName: string;
  roleKey: string;
  roleSort: number;
  dataScope: string;
  menu_check_strictly?: any;
  dept_check_strictly?: any;
  status: string;
  remark: string;
  createTime?: any;
  flag: boolean;
  superAdmin: boolean;
}

export interface User extends BaseCollectionModel {
  id: number;
  tenant_id: string;
  dept_id: number;
  user_name: string;
  nick_name: string;
  user_type: string;
  email: string;
  phonenumber: string;
  sex: string;
  avatar: string;
  status: string;
  login_ip: string;
  login_date: string;
  remark: string;
  create_time: string;
  dept: Dept;
  roles: Role[];
  roleIds?: string[];
  postIds?: string[];
  roleId: number;
  dept_name: string;
}

/**
 * @description 用户个人主页信息
 * @param user 用户信息
 * @param roleGroup 角色名称
 * @param postGroup 岗位名称
 */
export interface UserProfile {
  user: User;
  roleGroup: string;
  postGroup: string;
}

export interface UpdatePasswordParam {
  oldPassword: string;
  newPassword: string;
}

interface FileCallBack {
  name: string;
  file: Blob;
  filename: string;
}
