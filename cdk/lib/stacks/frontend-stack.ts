import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { StaticWebsiteHosting } from "../constructs/static-website-hosting";

export class FrontendStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new StaticWebsiteHosting(this, "StaticWebsiteHosting", {
      domainName: "sidekick.jimvid.xyz",
      certificateArn:
        "arn:aws:acm:us-east-1:211125707553:certificate/729577ea-6ecb-4453-a171-93c5a68034aa",
    });
  }
}
