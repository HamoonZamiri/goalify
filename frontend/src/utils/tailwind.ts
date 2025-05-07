export type FlexDirection = "row" | "col";

export const paddings = ["p-2", "p-4", "p-6", "p-8"] as const;
export type Padding = (typeof paddings)[number];

export const gaps = ["gap-2", "gap-4", "gap-6", "gap-8"] as const;
export type Gap = (typeof gaps)[number];

export const roundings = [
  "rounded",
  "rounded-md",
  "rounded-lg",
  "rounded-xl",
  "rounded-2xl",
] as const;
export type Rounding = (typeof roundings)[number];

export const shadows = [
  "shadow-none",
  "shadow-sm",
  "shadow",
  "shadow-md",
  "shadow-lg",
] as const;
export type Shadow = (typeof shadows)[number];

export const buttonVariants = ["primary", "secondary", "ghost"] as const;
export type ButtonVariant = (typeof buttonVariants)[number];

export const baseButtonClasses =
  "inline-flex items-center justify-center font-medium transition-colors duration-200 disabled:opacity-50 disabled:pointer-events-none";

export const buttonVariantClasses: Record<ButtonVariant, string> = {
  primary: "bg-blue-600 hover:bg-blue-500 text-white",
  secondary: "bg-gray-700 hover:bg-gray-600 text-white",
  ghost: "bg-transparent hover:bg-gray-800 text-white",
};

export const heights = [
  "h-0",
  "h-2",
  "h-4",
  "h-6",
  "h-8",
  "h-10",
  "h-12",
  "h-14",
  "h-16",
  "h-full",
  "h-1/2",
] as const;
export type Height = (typeof heights)[number];

export const widths = [
  "w-0",
  "w-2",
  "w-4",
  "w-6",
  "w-8",
  "w-8",
  "w-10",
  "w-12",
  "w-14",
  "w-16",
  "w-full",
  "w-1/2",
] as const;
export type Width = (typeof widths)[number];
