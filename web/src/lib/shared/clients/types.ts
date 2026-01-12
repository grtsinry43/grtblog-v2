// 定义 API 响应的通用结构
export interface ApiResponse<T> {
	code: number; // 0 成功，非 0 失败
	bizError: string; // 业务错误信息，code 非 0 时有值
	msg: string; // 服务器消息，用户友好信息，可直接展示在前端，code 0 时通常是成功提示，非 0 时可能是错误提示
	data: T;
	meta: {
		requestId: string; // 服务器生成的请求 ID，便于排查问题，可以让用户提供反馈或者上报
		timestamp: string; // 服务器响应时间，ISO 格式字符串
	};
}

export class BusinessError extends Error {
	code: number;
	bizError: string;

	constructor(code: number, msg: string, bizError: string) {
		super(msg);
		this.name = 'BusinessError';
		this.code = code;
		this.bizError = bizError;
	}
}