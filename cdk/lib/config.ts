import * as dotenv from "dotenv";
import path = require("path");

dotenv.config({ path: path.resolve(__dirname, "../.env") });

export type Environment = "dev" | "prod";

export interface EnvironmentConfig {
  env: Environment;
  apiDomainName: string;
  frontendDomainName: string;
  clerkSecret: string;
}

const configs: Record<Environment, EnvironmentConfig> = {
  dev: {
    env: "dev",
    apiDomainName: "dev.api.sidekick.jimvid.xyz",
    frontendDomainName: "dev.sidekick.jimvid.xyz",
    clerkSecret: process.env.CLERK_SECRET || "",
  },
  prod: {
    env: "prod",
    apiDomainName: "api.sidekick.jimvid.xyz",
    frontendDomainName: "sidekick.jimvid.xyz",
    clerkSecret: process.env.CLERK_SECRET || "",
  },
};

export const getConfig = (env: Environment): EnvironmentConfig => {
  return configs[env];
};
