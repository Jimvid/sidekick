#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { FrontendStack } from "../lib/stacks/frontend-stack";
import { ApiStack } from "../lib/stacks/api-stack";
import { getConfig, Environment } from "../lib/config";
import { SharedCertificate } from "../lib/constructs/shared-certificate";

const app = new cdk.App();

// Get environment from context (defaults to dev)
const environment = (app.node.tryGetContext("env") || "dev") as Environment;
const config = getConfig(environment);

// Stack name suffix based on environment
const suffix = environment === "prod" ? "" : `-${environment}`;

// Create certificate stack in us-east-1 (required for CloudFront and API Gateway)
const certStack = new cdk.Stack(app, `SidekickCertificate${suffix}`, {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: "us-east-1", // Must be us-east-1 for CloudFront
  },
  crossRegionReferences: true,
});

const sharedCert = new SharedCertificate(certStack, "SharedCertificate");

// S3 bucket for static website hosting
new FrontendStack(app, `SidekickFrontend${suffix}`, {
  config,
  certificate: sharedCert.certificate,
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
  crossRegionReferences: true,
});

// Api stack includes the DynamoDB table
new ApiStack(app, `SidekickApi${suffix}`, {
  config,
  certificate: sharedCert.certificate,
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
  crossRegionReferences: true,
});
