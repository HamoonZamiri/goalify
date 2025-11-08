import {
	type QueryFunction,
	type UseQueryOptions,
	useQuery,
} from "@tanstack/vue-query";
import { type Level, LevelSchema } from "@/features/levels/schemas";
import { zodFetch } from "@/shared/api";
import { isErrorResponse } from "@/shared/schemas";
import { API_BASE } from "@/utils/constants";
import { getLevelByIdParams, levelKeys } from "./queryKeys";

type LevelInfoQueryKey = ReturnType<typeof levelKeys.detail>;

const levelInfoQueryDataFn: QueryFunction<Level, LevelInfoQueryKey> = async ({
	queryKey,
}) => {
	const { levelId } = getLevelByIdParams(queryKey);

	const result = await zodFetch(`${API_BASE}/levels/${levelId}`, LevelSchema);

	if (isErrorResponse(result)) {
		throw new Error(result.message);
	}

	return result;
};

export const useLevelInfoQuery = (
	params: ReturnType<typeof getLevelByIdParams>,
	options?: Partial<
		UseQueryOptions<Level, Error, Level, Level, LevelInfoQueryKey>
	>,
) => {
	return useQuery({
		...options,
		queryKey: levelKeys.detail(params),
		queryFn: levelInfoQueryDataFn,
	});
};
