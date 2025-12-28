export type ID = number | string;
export type IDS = (number | string)[];

export interface BaseEntity {
  createBy?: string;
  createDept?: string;
  createTime?: string;
  updateBy?: string;
  updateTime?: string;
}

/**
 * 分页信息
 * @param rows 结果集
 * @param total 总数
 */
export interface PageResult<T = any> {
  rows: T[];
  total: number;
}

/**
 * 查询类型
 */
export interface QueryType {
  [key: string]:
    | 'AEQ'
    | 'AGE'
    | 'AGT'
    | 'AIN'
    | 'ALE'
    | 'ALIKE'
    | 'ALT'
    | 'ANE'
    | 'BETWEEN'
    | 'EQ'
    | 'GE'
    | 'GT'
    | 'IN'
    | 'LE'
    | 'LIKE'
    | 'LT'
    | 'NE';
}

/**
 * PB查询参数
 */
export interface PBQuery {
  params?: { [key: string]: any };
  sorts?: {
    field: string;
    order: 'asc' | 'desc' | string;
    sortTime?: numbner;
  }[];
  queryType?: QueryType;
  page?: {
    currentPage: number;
    pageSize: number;
  };
}

/**
 * 分页查询参数
 *
 * 排序支持的用法如下:
 * {isAsc:"asc",orderByColumn:"id"} order by id asc
 * {isAsc:"asc",orderByColumn:"id,createTime"} order by id asc,create_time asc
 * {isAsc:"desc",orderByColumn:"id,createTime"} order by id desc,create_time desc
 * {isAsc:"asc,desc",orderByColumn:"id,createTime"} order by id asc,create_time desc
 *
 * @param pageNum 当前页
 * @param pageSize 每页大小
 * @param orderByColumn 排序字段
 * @param isAsc 是否升序
 */
export interface PageQuery extends PBQuery {
  isAsc?: string;
  orderByColumn?: string;
  pageNum?: number;
  pageSize?: number;
  [key: string]: any;
}
