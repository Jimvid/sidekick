import { Construct } from "constructs";
import { CfnOutput, RemovalPolicy, StackProps } from "aws-cdk-lib";
import { Distribution, ViewerProtocolPolicy } from "aws-cdk-lib/aws-cloudfront";
import { S3StaticWebsiteOrigin } from "aws-cdk-lib/aws-cloudfront-origins";
import { BlockPublicAccess, Bucket } from "aws-cdk-lib/aws-s3";
import { BucketDeployment, Source } from "aws-cdk-lib/aws-s3-deployment";
import * as acm from "aws-cdk-lib/aws-certificatemanager";
import { ARecord, HostedZone, RecordTarget } from "aws-cdk-lib/aws-route53";
import { CloudFrontTarget } from "aws-cdk-lib/aws-route53-targets";

const path = "../apps/web/dist";

interface StaticWebsiteHostingProps extends StackProps {
  domainName: string;
  certificate: acm.ICertificate;
}

export class StaticWebsiteHosting extends Construct {
  constructor(scope: Construct, id: string, props: StaticWebsiteHostingProps) {
    super(scope, id);

    // Domain name
    const { domainName, certificate } = props;
    const rootDomain = domainName.split(".").slice(-2).join(".");
    const subdomainPart = domainName.replace(`.${rootDomain}`, "");

    const hostedZone = HostedZone.fromLookup(this, "HostedZone", {
      domainName: rootDomain,
    });

    // Hosting bucket
    const bucket = new Bucket(this, "FrontendBucket", {
      autoDeleteObjects: true,
      removalPolicy: RemovalPolicy.DESTROY,
      websiteIndexDocument: "index.html",
      websiteErrorDocument: "index.html",
      publicReadAccess: true,
      blockPublicAccess: new BlockPublicAccess({
        blockPublicAcls: false,
        blockPublicPolicy: false,
        ignorePublicAcls: false,
        restrictPublicBuckets: false,
      }),
    });

    // CloudFront Distribution
    const distribution = new Distribution(this, "CloudfrontDistribution", {
      defaultBehavior: {
        origin: new S3StaticWebsiteOrigin(bucket),
        viewerProtocolPolicy: ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
      },
      defaultRootObject: "index.html",
      errorResponses: [
        {
          httpStatus: 404,
          responseHttpStatus: 200,
          responsePagePath: "/index.html",
        },
      ],
      domainNames: [domainName],
      certificate: certificate,
    });

    // Outputs
    new ARecord(this, "ARecord", {
      zone: hostedZone,
      recordName: subdomainPart,
      target: RecordTarget.fromAlias(new CloudFrontTarget(distribution)),
    });

    new BucketDeployment(this, "BucketDeployment", {
      sources: [Source.asset(path)],
      destinationBucket: bucket,
      distribution,
      distributionPaths: ["/*"],
    });

    new CfnOutput(this, "CloudFrontURL", {
      value: distribution.domainName,
      description: "CloudFront Distribution URL",
    });

    new CfnOutput(this, "CustomDomainURL", {
      value: `https://${domainName}`,
      description: "Custom Domain URL",
    });

    new CfnOutput(this, "BucketName", {
      value: bucket.bucketName,
      description: "Hosting bucket for frontend",
    });
  }
}
