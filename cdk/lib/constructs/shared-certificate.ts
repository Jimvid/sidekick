import * as acm from "aws-cdk-lib/aws-certificatemanager";
import * as route53 from "aws-cdk-lib/aws-route53";
import { Construct } from "constructs";

export class SharedCertificate extends Construct {
  public readonly certificate: acm.Certificate;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    // Look up the hosted zone
    const hostedZone = route53.HostedZone.fromLookup(this, "HostedZone", {
      domainName: "jimvid.xyz",
    });

    // Create wildcard certificate that covers all subdomains
    this.certificate = new acm.Certificate(this, "WildcardCertificate", {
      domainName: "*.sidekick.jimvid.xyz",
      subjectAlternativeNames: [
        "sidekick.jimvid.xyz",
        "*.api.sidekick.jimvid.xyz", // Covers dev.api, prod api, etc.
      ],
      validation: acm.CertificateValidation.fromDns(hostedZone),
    });
  }
}
