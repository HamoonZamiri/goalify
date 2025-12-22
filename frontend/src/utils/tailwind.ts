export type FlexDirection = "row" | "col";

export type Padding = "p-2" | "p-4" | "p-6" | "p-8";

export type Gap = "gap-1" | "gap-2" | "gap-4" | "gap-6" | "gap-8";

export type Rounding =
	| "rounded-none"
	| "rounded-sm"
	| "rounded-md"
	| "rounded-lg"
	| "rounded-xl"
	| "rounded-2xl";

export type Shadow =
	| "shadow-none"
	| "shadow-xs"
	| "shadow-sm"
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
	| "h-4/5"
	| "h-auto";

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

export type IconName =
	| "arrow-path"
	| "check-outline"
	| "chevron-up"
	| "gear"
	| "gift"
	| "home"
	| "signout"
	| "trash"
	| "x-mark";

export type IconSize = "sm" | "md" | "lg";

export const iconSizeClasses: Record<IconSize, string> = {
	sm: "size-5",
	md: "size-6",
	lg: "size-8",
};

export type IconFill =
	| "none"
	| "currentColor"
	| "green"
	| "red"
	| "blue"
	| "yellow";

export const dialogPatterns = {
	base: "fixed inset-0 z-50 flex",
	backdrop: "fixed inset-0 bg-black/50",
	centered: "items-center justify-center",
	sideDrawer: "items-stretch justify-end",
	sizes: {
		sm: "w-[90vw] max-w-md",
		md: "w-[90vw] max-w-lg",
		lg: "w-[90vw] max-w-2xl",
		drawer: "w-full sm:w-1/2 max-w-2xl",
	},
} as const;

export type DialogSize = keyof typeof dialogPatterns.sizes;
export type DialogPattern = Exclude<
	keyof typeof dialogPatterns,
	"base" | "backdrop" | "sizes"
>;
