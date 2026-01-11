import * as cdk from "aws-cdk-lib";
import * as acm from "aws-cdk-lib/aws-certificatemanager";
import { Construct } from "constructs";
import { StaticWebsiteHosting } from "../constructs/static-website-hosting";
import { EnvironmentConfig } from "../config";

interface FrontendStackProps extends cdk.StackProps {
  config: EnvironmentConfig;
  certificate: acm.ICertificate;
}

export class FrontendStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: FrontendStackProps) {
    super(scope, id, props);

    new StaticWebsiteHosting(this, "StaticWebsiteHosting", {
      domainName: props.config.frontendDomainName,
      certificate: props.certificate,
    });
  }
}
