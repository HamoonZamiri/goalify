import { z } from "zod";

/**
 * Generic server response wrapper schema
 */
export function createServerResponseSchema<TData extends z.ZodTypeAny>(
	schema: TData,
) {
	return z.object({
		data: schema,
	});
}

/**
 * Generic array schema helper
 */
export function createArraySchema<TData extends z.ZodTypeAny>(schema: TData) {
	return z.array(schema);
}

/**
 * Error response type from the server
 */
export type ErrorResponse = {
	statusCode?: number;
	message: string;
	errors?: Record<string, string>;
};

/**
 * Type guard to check if a response is an error
 */
export function isErrorResponse(
	res: unknown | ErrorResponse,
): res is ErrorResponse {
	const casted = res as ErrorResponse;
	return (
		casted.errors !== undefined ||
		(casted.statusCode !== undefined && casted.statusCode >= 400)
	);
}
