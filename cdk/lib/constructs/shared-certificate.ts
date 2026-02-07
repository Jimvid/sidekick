import * as acm from "aws-cdk-lib/aws-certificatemanager";
import * as route53 from "aws-cdk-lib/aws-route53";
import { Construct } from "constructs";

export interface SharedCertificateProps {
  rootDomainName: string;
  appDomainName: string;
}

export class SharedCertificate extends Construct {
  public readonly certificate: acm.Certificate;

  constructor(scope: Construct, id: string, props: SharedCertificateProps) {
    super(scope, id);

    const { rootDomainName, appDomainName } = props;

    // Look up the hosted zone
    const hostedZone = route53.HostedZone.fromLookup(this, "HostedZone", {
      domainName: rootDomainName,
    });

    // Create wildcard certificate that covers all subdomains
    this.certificate = new acm.Certificate(this, "WildcardCertificate", {
      domainName: `*.${appDomainName}`,
      subjectAlternativeNames: [
        appDomainName,
        `*.api.${appDomainName}`, // Covers dev.api, prod api, etc.
      ],
      validation: acm.CertificateValidation.fromDns(hostedZone),
    });
  }
}
