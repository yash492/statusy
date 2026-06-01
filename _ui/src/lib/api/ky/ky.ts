import { PUBLIC_API_SERVER_ROUTE } from '$env/static/public';
import ky, { HTTPError, TimeoutError, type Options } from 'ky';

type Response<T> =
	| [data: T, error: null]
	| [data: null, error: AppError];

export interface AppError {
	code: string;
	message: string;
}

// Helper guard functions (assuming standard ky instances)
const isHTTPError = (e: any): e is HTTPError => e instanceof HTTPError;
const isTimeoutError = (e: any): e is TimeoutError => e instanceof TimeoutError;
const isNetworkError = (e: any): boolean => e instanceof Error && e.name === 'TypeError' && e.message === 'Failed to fetch';

export async function safeAsync<T>(promise: Promise<T>): Promise<Response<T>> {
	try {
		const data = await promise;
		return [data, null];
	} catch (e) {
		// Handle HTTP Errors and asynchronously parse the error body
		if (isHTTPError(e)) {
			let code = 'HTTP_ERROR';
			let message = e.message;

			try {
				// ky exposes the response object; you must parse the body explicitly
				const errorData = await e.response.json() as Record<string, unknown>;
				if (typeof errorData?.code === 'string') code = errorData.code;
				if (typeof errorData?.message === 'string') message = errorData.message;
			} catch {
				// Fallback if error response is not valid JSON
				message = `HTTP Error ${e.response.status}: ${e.response.statusText}`;
			}

			return [null, { code, message }];
		}

		if (isTimeoutError(e)) {
			return [null, { code: 'TIMEOUT_ERROR', message: 'The request timed out. Please try again.' }];
		}

		if (isNetworkError(e)) {
			return [null, { code: 'NETWORK_ERROR', message: 'Network connection failed. Please check your internet connection.' }];
		}

		const error = e as Error;
		return [null, { code: 'UNKNOWN_ERROR', message: error?.message || 'An unexpected error occurred.' }];
	}
}

const KyClient = ky.create({
	prefix: PUBLIC_API_SERVER_ROUTE,
	headers: {
		accept: 'application/json'
	}
});

// Note: delete methods often return no content (204), .json() might fail on empty bodies.
// Rendered with unknown/void support depending on endpoint expectations.
export const ApiClient = {
	get: <T>(url: string, options?: Options) => safeAsync(KyClient.get(url, options).json<T>()),
	post: <T>(url: string, options?: Options) => safeAsync(KyClient.post(url, options).json<T>()),
	put: <T>(url: string, options?: Options) => safeAsync(KyClient.put(url, options).json<T>()),
	delete: <T = unknown>(url: string, options?: Options) => safeAsync(KyClient.delete(url, options).json<T>())
};