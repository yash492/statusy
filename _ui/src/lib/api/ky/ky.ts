import { PUBLIC_API_SERVER_ROUTE } from '$env/static/public';
import ky, { HTTPError } from 'ky';

interface AppError {
	message: string
	code: string
}

type Response<T> = [data: T | null, error: AppError | null];
export async function safeAsync<T>(promise: Promise<T>): Promise<Response<T>> {
	try {
		const data = (await promise) as T
		return [data, null]
	} catch (e) {
		//TODO: make it more robust and handle network errors
		const error = e as HTTPError
		return [
			null, {
				code: (error.data as any).code,
				message: (error.data as any).message
			}
		]
	}
}

const KyClient = ky.create({
	prefix: PUBLIC_API_SERVER_ROUTE,
	headers: {
		accept: 'application/json'
	},
	hooks: {

	}

});

export default KyClient;


