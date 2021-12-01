

export interface FormSource {
  label?: string;
  key: string;
  type: string;
  icon?: string;
  required?: boolean;
  multiple?: boolean;
  initValue: string | string[];
  rules?: object;
  minRows?: number;
  maxRows?: number;
  // options?: {
  //   label: any,
  //   value: string
  // }[]
  // options?: Option[]
  options?: any;
  filterable?: boolean;
  allowCreate?: boolean;
  createDefaultFirst?: boolean;
  class?: string;
  width?: string;
  placeholder?: string;
}

/**
 * 登陆
 */
export interface LoginModule {
  corpCode: string;
  loginName: string;
  password: string;
  verifyCode: string;
  captchaKey: string;
  type: number;
}

/**
 * 用户信息
 */
export interface UserInfoModule {
  id: string;
  loginName: string;
  telephone: string;
  email: string;
  businessId: string;
  picturePath: string;
  status: number;
  createTime: number;
  createUserId: string;
  lastUpdateUserId: string;
  lastUpdateTime: number;
}
