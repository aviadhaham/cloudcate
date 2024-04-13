import react from '@vitejs/plugin-react-swc'
import path from 'path';
import { defineConfig } from 'vite'

const defaultConfig = {
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
};

export default defineConfig(({ command, mode }) => {
  if (command === "serve") {
    const isDev = mode === "development";

    return {
      ...defaultConfig,
      server: {
        proxy: {
          "/api": isDev ? "http://localhost:80" : "/api",
        },
      },
    };
  } else {
    return defaultConfig;
  }
});
