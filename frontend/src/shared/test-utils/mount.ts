import {
	mount,
	type ComponentMountingOptions,
	type VueWrapper,
} from "@vue/test-utils";
import { VueQueryPlugin } from "@tanstack/vue-query";
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
export function mountWithPlugins<T>(
	component: T,
	options?: ComponentMountingOptions<T>,
): VueWrapper<any> {
	return mount(component, {
		...options,
		global: {
			...options?.global,
			plugins: [
				...(options?.global?.plugins || []),
				[VueQueryPlugin, { queryClient }],
			],
		},
	});
}
