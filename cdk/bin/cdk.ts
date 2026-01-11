#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { FrontendStack } from "../lib/stacks/frontend-stack";
import { ApiStack } from "../lib/stacks/api-stack";

const app = new cdk.App();

// S3 bucket for static website hosting
new FrontendStack(app, "SidekickFrontend", {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
});

// Api stack includes the DynamoDB table
new ApiStack(app, "SidekickApi", {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
});
