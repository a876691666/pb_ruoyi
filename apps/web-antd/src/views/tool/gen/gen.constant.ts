/** Entity基类字段 */
export const BASE_ENTITY: string[] = [
  'create_by',
  'create_time',
  'del_flag',
  'update_by',
  'update_time',
  'remark',
];

/** 页面不需要编辑字段 */
export const COLUMNNAME_NOT_EDIT: string[] = [
  'id',
  'create_by',
  'create_time',
  'del_flag',
];

/** 页面不需要插入字段 */
/** 页面不需要插入字段 */
export const COLUMNNAME_NOT_INSERT: string[] = [
  'id',
  'create_by',
  'create_time',
  'del_flag',
];

/** 页面不需要显示的列表字段 */
export const COLUMNNAME_NOT_LIST: string[] = [
  'create_by',
  'create_time',
  'del_flag',
  'update_by',
  'update_time',
];

/** 页面不需要查询字段 */
/** 页面不需要查询字段 */
export const COLUMNNAME_NOT_QUERY: string[] = [
  'id',
  'create_by',
  'create_time',
  'del_flag',
  'update_by',
  'update_time',
  'remark',
];

/** 数据库数字类型 */
export const COLUMNTYPE_NUMBER: string[] = [
  'smallint',
  'mediumint',
  'int',
  'number',
  'integer',
  'bit',
  'bigint',
  'float',
  'double',
  'decimal',
];

/** 数据库字符串类型 */
/** 数据库字符串类型 */
export const COLUMNTYPE_STR: string[] = [
  'char',
  'varchar',
  'nvarchar',
  'varchar2',
];

/** 数据库文本类型 */
export const COLUMNTYPE_TEXT: string[] = ['text', 'mediumtext', 'longtext'];

/** 数据库时间类型 */
/** 数据库时间类型 */
export const COLUMNTYPE_TIME: string[] = [
  'datetime',
  'time',
  'date',
  'timestamp',
];

/** 复选框 */
export const HTML_CHECKBOX: string = 'checkbox';
/** 日期控件 */
export const HTML_DATETIME: string = 'datetime';

/** 富文本控件 */
export const HTML_EDITOR: string = 'editor';

/** 文件上传控件 */
export const HTML_FILE_UPLOAD: string = 'fileUpload';

/** 图片上传控件 */
export const HTML_IMAGE_UPLOAD: string = 'imageUpload';

/** 文本框 */
export const HTML_INPUT: string = 'input';

/** 单选框 */
export const HTML_RADIO: string = 'radio';

/** 下拉框 */
export const HTML_SELECT: string = 'select';

/** 文本域 */
export const HTML_TEXTAREA: string = 'textarea';

/** 不需要 */
export const NOT_REQUIRE: string = '0';

/** 上级菜单ID字段 */
export const PARENT_MENU_ID: string = 'parentMenuId';

/** 上级菜单名称字段 */
export const PARENT_MENU_NAME: string = 'parentMenuName';

/** 日期区间查询 */
export const QUERY_BETWEEN: string = 'BETWEEN';

/** 相等查询 */
export const QUERY_EQ: string = 'EQ';

/** 大于查询 */
export const QUERY_GT: string = 'GT';

/** 大于等于查询 */
export const QUERY_GTE: string = 'GTE';

/** 模糊查询 */
export const QUERY_LIKE: string = 'LIKE';

/** 小于查询 */
export const QUERY_LT: string = 'LT';

/** 小于等于查询 */
export const QUERY_LTE: string = 'LTE';

/** 不等于查询 */
export const QUERY_NE: string = 'NE';

/** 需要 */
export const REQUIRE: string = '1';

/** 单表（增删改查） */
export const TPL_CRUD: string = 'crud';

/** 主子表（增删改查） */
export const TPL_SUB: string = 'sub';

/** 树表（增删改查） */
export const TPL_TREE: string = 'tree';

/** 树编码字段 */
export const TREE_CODE: string = 'treeCode';

/** Tree基类字段 */
/** Tree基类字段 */
export const TREE_ENTITY: string[] = [
  'parentName',
  'parentId',
  'orderNum',
  'ancestors',
  'children',
];

/** 树名称字段 */
export const TREE_NAME: string = 'treeName';

/** 树父编码字段 */
export const TREE_PARENT_CODE: string = 'treeParentCode';

/** 高精度计算类型 */
export const TYPE_BIGDECIMAL: string = 'BigDecimal';

/** 时间类型 */
export const TYPE_DATE: string = 'Date';

/** 浮点型 */
export const TYPE_DOUBLE: string = 'Double';

/** 整型 */
export const TYPE_INTEGER: string = 'Integer';

/** 长整型 */
export const TYPE_LONG: string = 'Long';

export const TYPE_NUMBER: string = 'Number';

/** 字符串类型 */
export const TYPE_STRING: string = 'String';
