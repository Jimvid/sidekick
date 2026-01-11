import * as cdk from "aws-cdk-lib";
import * as acm from "aws-cdk-lib/aws-certificatemanager";
import { Construct } from "constructs";
import { ApiWithDynamo } from "../constructs/api-with-dynamo";
import { EnvironmentConfig } from "../config";

interface ApiStackProps extends cdk.StackProps {
  config: EnvironmentConfig;
  certificate: acm.ICertificate;
}

export class ApiStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: ApiStackProps) {
    super(scope, id, props);

    new ApiWithDynamo(this, "ApiWithDynamo", {
      domainName: props.config.apiDomainName,
      certificate: props.certificate,
      config: props.config,
    });
  }
}
