import * as dotenv from "dotenv";
import path = require("path");

dotenv.config({ path: path.resolve(__dirname, "../.env") });

export type ConfigProps = {
  CLERK_SECRET: string;
};

export const getConfig = (): ConfigProps => ({
  CLERK_SECRET: process.env.CLERK_SECRET || "",
});
