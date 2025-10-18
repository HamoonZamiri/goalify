import { useQuery } from "@tanstack/vue-query";
import { zodFetch } from "@/shared/api";
import { isErrorResponse } from "@/shared/schemas";
import { LevelSchema, type Level } from "@/features/levels/schemas";
import { levelKeys } from "./queryKeys";
import { API_BASE } from "@/utils/constants";

async function levelInfoQueryDataFn(levelId: number): Promise<Level> {
	const result = await zodFetch(`${API_BASE}/levels/${levelId}`, LevelSchema);

	if (isErrorResponse(result)) {
		throw new Error(result.message);
	}

	return result;
}

export function useLevelInfo(levelId: number) {
	return useQuery({
		queryKey: levelKeys.detail(levelId),
		queryFn: () => levelInfoQueryDataFn(levelId),
		enabled: !!levelId,
	});
}
