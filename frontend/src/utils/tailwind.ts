export type FlexDirection = "row" | "col";

export type Padding = "p-2" | "p-4" | "p-6" | "p-8";

export type Gap = "gap-1" | "gap-2" | "gap-4" | "gap-6" | "gap-8";

export type Rounding =
	| "rounded"
	| "rounded-md"
	| "rounded-lg"
	| "rounded-xl"
	| "rounded-2xl";

export type Shadow =
	| "shadow-none"
	| "shadow-sm"
	| "shadow"
	| "shadow-md"
	| "shadow-lg";

export type ButtonVariant = "primary" | "secondary" | "ghost";

export const baseButtonClasses =
	"inline-flex items-center justify-center font-medium transition-colors duration-200 disabled:opacity-50 disabled:pointer-events-none";

export const buttonVariantClasses: Record<ButtonVariant, string> = {
	primary: "bg-blue-600 hover:bg-blue-500 text-white",
	secondary: "bg-gray-700 hover:bg-gray-600 text-white",
	ghost: "bg-transparent hover:bg-gray-800 text-white",
};

export type Height =
	| "h-0"
	| "h-2"
	| "h-4"
	| "h-6"
	| "h-8"
	| "h-10"
	| "h-12"
	| "h-14"
	| "h-16"
	| "h-full"
	| "h-1/2"
	| "h-1/3"
	| "h-1/4"
	| "h-4/5";

export type Width =
	| "w-0"
	| "w-2"
	| "w-4"
	| "w-6"
	| "w-8"
	| "w-10"
	| "w-12"
	| "w-14"
	| "w-16"
	| "w-full"
	| "w-1/2"
	| "w-1/3"
	| "w-1/4"
	| "w-4/5"
	| "w-auto";

export type TextColor = "light" | "dark" | "error";

export const textColorMap: Record<TextColor, string> = {
	light: "text-gray-300",
	dark: "text-gray-900",
	error: "text-red-500",
};
