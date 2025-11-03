import { VueQueryPlugin } from "@tanstack/vue-query";
import {
	type ComponentMountingOptions,
	mount,
	type VueWrapper,
} from "@vue/test-utils";
import { queryClient } from "@/shared/api/query-client";

/**
 * Custom mount function that provides common test setup
 * Includes VueQueryPlugin by default for TanStack Query components
 *
 * @example
 * ```ts
 * const wrapper = mountWithPlugins(MyComponent, {
 *   props: { foo: 'bar' }
 * });
 * ```
 *
 * @example
 * // Extend with custom plugins or global config
 * ```ts
 * const wrapper = mountWithPlugins(MyComponent, {
 *   props: { foo: 'bar' },
 *   global: {
 *     stubs: { CustomComponent: true }
 *   }
 * });
 * ```
 */
export function mountWithPlugins<T extends Record<string, unknown>>(
	component: T,
	options?: ComponentMountingOptions<T>,
): VueWrapper {
	return mount(component, {
		...options,
		global: {
			...options?.global,
			plugins: [
				...(options?.global?.plugins || []),
				[VueQueryPlugin, { queryClient }],
			],
		},
	}) as VueWrapper;
}
