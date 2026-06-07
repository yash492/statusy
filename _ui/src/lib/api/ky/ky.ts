import ky, { isHTTPError, isNetworkError, isTimeoutError, type Options } from 'ky';

interface AppError {
	message: string;
	code: string;
}

type Response<T> = [data: T | null, error: AppError | null];
export async function safeAsync<T>(promise: Promise<T>): Promise<Response<T>> {
	try {
		const data = (await promise) as T;
		return [data, null];
	} catch (e) {
		if (isHTTPError(e)) {
			const errorData = e.data as Record<string, unknown> | undefined;
			return [
				null,
				{
					code: typeof errorData?.code === 'string' ? errorData.code : 'HTTP_ERROR',
					message: typeof errorData?.message === 'string' ? errorData.message : e.message
				}
			];
		}

		if (isNetworkError(e)) {
			return [
				null,
				{
					code: 'NETWORK_ERROR',
					message: 'Network connection failed. Please check your internet connection.'
				}
			];
		}

		if (isTimeoutError(e)) {
			return [
				null,
				{
					code: 'TIMEOUT_ERROR',
					message: 'The request timed out. Please try again.'
				}
			];
		}

		const error = e as Error;
		return [
			null,
			{
				code: 'UNKNOWN_ERROR',
				message: error?.message || 'An unexpected error occurred.'
			}
		];
	}
}

const KyClient = ky.create({
	prefix: "/api",
	headers: {
		accept: 'application/json'
	},
	hooks: {}
});

export const ApiClient = {
	get: <T>(url: string, options?: Options) => safeAsync(KyClient.get(url, options).json<T>()),
	post: <T>(url: string, options?: Options) => safeAsync(KyClient.post(url, options).json<T>()),
	put: <T>(url: string, options?: Options) => safeAsync(KyClient.put(url, options).json<T>()),
	delete: (url: string, options?: Options) => safeAsync(KyClient.delete(url, options))
};

export default KyClient;
