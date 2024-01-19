import axios, { AxiosError } from 'axios';

type ErrorResponse = {
	error_msg: string;
	is_error: boolean;
	status_code: number;
};

export function AxiosResponseErr(err: any) {
	if (axios.isAxiosError(err)) {
		return err.response?.data as ErrorResponse;
	}

	const errResp: ErrorResponse = {
		error_msg: 'Failed to connect the server',
		is_error: true,
		status_code: 500
	};

	return errResp;
}
