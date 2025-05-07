import { defineConfig, mergeConfig } from "vitest/config";
import viteConfig from "./vite.config";

export default mergeConfig(
  viteConfig,
  defineConfig({
    test: {
      environment: "jsdom",
      dir: "./src",
      setupFiles: "./src/vitest.setup.ts",
    },
  }),
);
