import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { ApiWithDynamo } from "../constructs/api-with-dynamo";

export class ApiStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new ApiWithDynamo(this, "ApiWithDynamo", {
      domainName: "api.sidekick.jimvid.xyz",
      certificateArn:
        "arn:aws:acm:us-east-1:211125707553:certificate/96aa108c-f031-4efc-b522-1b849005aa83",
    });
  }
}
